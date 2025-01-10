package constants

type Fn func()

type TransactionType int

const (
	Expenses TransactionType = iota
	Income
)

type Periodicity int

const (
	Daily Periodicity = iota
	Weekly
	Monthly
)

func (t TransactionType) ToString() string {
	return [...]string{"Expenses", "Income"}[t]
}

func (t TransactionType) ToArrayString() []string {
	return []string{"Expenses", "Income"}
}

func (p Periodicity) ToString() string {
	return [...]string{"Daily", "Weekly", "Monthly"}[p]
}

func (p Periodicity) ToArrayString() []string {
	return []string{"Daily", "Weekly", "Monthly"}
}
