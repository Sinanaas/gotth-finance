package constants

type TransactionType int

const (
	Expenses TransactionType = iota
	Income
	Transfer
	Investment
)

type Periodicity int

const (
	Daily Periodicity = iota
	Weekly
	Monthly
)

func (t TransactionType) ToString() string {
	return [...]string{"Expenses", "Income", "Transfer", "Investment"}[t]
}

func GetTransactionTypes() []TransactionType {
	return []TransactionType{Expenses, Income, Transfer, Investment}
}
