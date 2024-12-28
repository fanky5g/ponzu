package constants

type ComparisonOperator int

const (
	Equal ComparisonOperator = iota + 1
	LessThan
	GreaterThan
	GreaterThanOrEqualTo
	LessThanOrEqualTo
)
