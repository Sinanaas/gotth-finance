package constants

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

func (t TransactionType) ToIndex() int {
	return int(t)
}

func (t TransactionType) ToArrayString() []string {
	return []string{"Expenses", "Income"}
}

func GetTransactionTypes() []TransactionType {
	return []TransactionType{Expenses, Income}
}
