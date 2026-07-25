package managers

import (
	"fmt"
	"github.com/Sinanaas/gotth-financial-tracker/internal/constants"
	"github.com/Sinanaas/gotth-financial-tracker/internal/models"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"google.golang.org/appengine/log"
	"gorm.io/gorm"
	"time"
)

func (m *BasicManager) GetRecurrings(userId string) ([]models.Recurring, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	var recurrings []models.Recurring
	if err := m.DB.Preload("Category").Preload("Account").Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&recurrings).Error; err != nil {
		return nil, err
	}

	return recurrings, nil
}

func (m *BasicManager) CreateRecurring(payload models.RecurringRequest) error {
	userUUID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return err
	}

	categoryUUID, err := uuid.Parse(payload.CategoryID)
	if err != nil {
		return err
	}

	accountUUID, err := uuid.Parse(payload.AccountID)
	if err != nil {
		return err
	}

	transactionDate, err := time.Parse("2006-01-02", payload.StartDate)
	if err != nil {
		return err
	}

	recurring := models.Recurring{
		UserID:          userUUID,
		Name:            payload.Name,
		StartDate:       transactionDate,
		Amount:          payload.Amount,
		TransactionType: constants.TransactionType(payload.TransactionType),
		Periodicity:     constants.Periodicity(payload.Periodicity),
		CategoryID:      categoryUUID,
		AccountID:       accountUUID,
	}

	fn := constants.Fn(func() {
		transactionRequest := models.TransactionRequest{
			UserID:      payload.UserID,
			Amount:      payload.Amount,
			Type:        payload.TransactionType,
			Description: payload.Name,
			CategoryID:  payload.CategoryID,
			Account:     payload.AccountID,
			Date:        time.Now().Format("2006-01-02"),
		}

		if err := m.CreateTransaction(transactionRequest); err != nil {
			log.Errorf(nil, "❌ failed to create transaction: %v", err)
		}
	})

	if err := m.DB.Create(&recurring).Error; err != nil {
		return err
	}

	if err := m.SetUserCRONJob(recurring, *m.GoCRON, fn); err != nil {
		return err
	}

	return nil
}

func (m *BasicManager) FindUserRecurringById(userId, id string) (models.Recurring, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return models.Recurring{}, err
	}
	recUUID, err := uuid.Parse(id)
	if err != nil {
		return models.Recurring{}, err
	}
	var rec models.Recurring
	if err := m.DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", recUUID, userUUID).First(&rec).Error; err != nil {
		return models.Recurring{}, err
	}
	return rec, nil
}

// UpdateRecurring edits the row (scoped to user) and replaces its scheduled job.
func (m *BasicManager) UpdateRecurring(id string, payload models.RecurringRequest) error {
	userUUID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return err
	}
	recUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	categoryUUID, err := uuid.Parse(payload.CategoryID)
	if err != nil {
		return err
	}
	accountUUID, err := uuid.Parse(payload.AccountID)
	if err != nil {
		return err
	}
	startDate, err := time.Parse("2006-01-02", payload.StartDate)
	if err != nil {
		return err
	}

	var recurring models.Recurring
	if err := m.DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", recUUID, userUUID).First(&recurring).Error; err != nil {
		return err
	}

	// Drop the old schedule before applying the new one.
	if err := m.RemoveCRONJob(recurring, *m.GoCRON); err != nil {
		return err
	}

	recurring.Name = payload.Name
	recurring.StartDate = startDate
	recurring.Amount = payload.Amount
	recurring.TransactionType = constants.TransactionType(payload.TransactionType)
	recurring.Periodicity = constants.Periodicity(payload.Periodicity)
	recurring.CategoryID = categoryUUID
	recurring.AccountID = accountUUID

	if err := m.DB.Save(&recurring).Error; err != nil {
		return err
	}

	fn := constants.Fn(func() {
		if err := m.CreateTransaction(models.TransactionRequest{
			UserID:      payload.UserID,
			Amount:      payload.Amount,
			Type:        payload.TransactionType,
			Description: payload.Name,
			CategoryID:  payload.CategoryID,
			Account:     payload.AccountID,
			Date:        time.Now().Format("2006-01-02"),
		}); err != nil {
			log.Errorf(nil, "❌ failed to create recurring transaction: %v", err)
		}
	})

	return m.SetUserCRONJob(recurring, *m.GoCRON, fn)
}

func (m *BasicManager) GetUserUpcomingRecurring(userId string) (models.Recurring, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return models.Recurring{}, err
	}

	var recurrings []models.Recurring
	if err := m.DB.Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&recurrings).Error; err != nil {
		return models.Recurring{}, err
	}

	if len(recurrings) == 0 {
		return models.Recurring{}, nil
	}

	today := time.Now()
	var closestRecurring models.Recurring
	minDiff := time.Duration(1<<63 - 1) // Max duration

	for _, recurring := range recurrings {
		nextOccurrence := GetNextOccurrence(recurring.StartDate, recurring.Periodicity, today)
		diff := nextOccurrence.Sub(today)
		if diff >= 0 && diff < minDiff {
			minDiff = diff
			closestRecurring = recurring
		}
	}

	if minDiff == time.Duration(1<<63-1) {
		return models.Recurring{}, gorm.ErrRecordNotFound
	}

	return closestRecurring, nil
}

func (m *BasicManager) DeleteRecurringById(recurringId, userId string) error {
	recurringUUID, err := uuid.Parse(recurringId)
	if err != nil {
		return err
	}
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}

	var recurring models.Recurring
	if err := m.DB.Where("id = ? AND user_id = ?", recurringUUID, userUUID).First(&recurring).Error; err != nil {
		return err
	}

	if err := m.RemoveCRONJob(recurring, *m.GoCRON); err != nil {
		return err
	}

	recurring.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	if err := m.DB.Save(&recurring).Error; err != nil {
		return err
	}

	return nil
}

func (m *BasicManager) GetAllUpcomingRecurring(userId string) ([]models.Recurring, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	var recurrings []models.Recurring
	if err := m.DB.Preload("Category").Preload("Account").Where("user_id = ? AND deleted_at IS NULL", userUUID).Find(&recurrings).Error; err != nil {
		return nil, err
	}

	return recurrings, nil
}

func GetNextOccurrence(startDate time.Time, frequency constants.Periodicity, today time.Time) time.Time {
	nextOccurrence := startDate
	for nextOccurrence.Before(today) {
		switch frequency {
		case constants.Daily:
			nextOccurrence = nextOccurrence.AddDate(0, 0, 1)
		case constants.Weekly:
			nextOccurrence = nextOccurrence.AddDate(0, 0, 7)
		case constants.Monthly:
			nextOccurrence = nextOccurrence.AddDate(0, 1, 0)
		}
	}
	return nextOccurrence
}

func (m *BasicManager) SetUserCRONJob(recurring models.Recurring, scheduler gocron.Scheduler, taskFn constants.Fn) error {
	var err error
	var j gocron.Job

	switch recurring.Periodicity {
	case constants.Daily:
		j, err = scheduler.NewJob(
			gocron.DailyJob(
				1,
				gocron.NewAtTimes(
					gocron.NewAtTime(0, 0, 0)),
			),
			gocron.NewTask(taskFn),
		)
	case constants.Weekly:
		j, err = scheduler.NewJob(
			gocron.WeeklyJob(
				1,
				gocron.NewWeekdays(recurring.StartDate.Weekday()),
				gocron.NewAtTimes(
					gocron.NewAtTime(0, 0, 0)),
			),
			gocron.NewTask(taskFn),
		)
	case constants.Monthly:
		j, err = scheduler.NewJob(
			gocron.MonthlyJob(
				1,
				gocron.NewDaysOfTheMonth(recurring.StartDate.Day()),
				gocron.NewAtTimes(
					gocron.NewAtTime(0, 0, 0)),
			),
			gocron.NewTask(taskFn),
		)
	default:
		return fmt.Errorf("❌ invalid periodicity")
	}

	if err != nil {
		return fmt.Errorf("❌ failed to create new job: %w", err)
	}

	tx := m.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	recurring.JobID = j.ID()
	recurring.JobName = recurring.Name

	if err := m.DB.Save(&recurring).Error; err != nil {
		tx.Rollback()
	}

	tx.Commit()

	return nil
}

func (m *BasicManager) RemoveCRONJob(recurring models.Recurring, scheduler gocron.Scheduler) error {
	if err := scheduler.RemoveJob(recurring.JobID); err != nil {
		return err
	}

	return nil
}

func (m *BasicManager) LoadAndScheduleJobs() error {
	var recurrings []models.Recurring
	if err := m.DB.Find(&recurrings).Error; err != nil {
		return err
	}

	for _, recurring := range recurrings {
		r := recurring // capture loop variable
		fn := constants.Fn(func() {
			transactionRequest := models.TransactionRequest{
				UserID:      r.UserID.String(),
				Amount:      r.Amount,
				Type:        int(r.TransactionType),
				Description: r.Name,
				CategoryID:  r.CategoryID.String(),
				Account:     r.AccountID.String(),
				Date:        time.Now().Format("2006-01-02"),
			}

			if err := m.CreateTransaction(transactionRequest); err != nil {
				log.Errorf(nil, "❌ failed to create transaction: %v", err)
			}
		})

		if err := m.SetUserCRONJob(r, *m.GoCRON, fn); err != nil {
			return err
		}
	}
	return nil
}
