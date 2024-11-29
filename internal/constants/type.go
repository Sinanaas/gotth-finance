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
