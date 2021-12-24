package errors

import "strconv"

type errnoGenerator interface {
	Value() string
	Next() errnoGenerator
}

type errnoSequenceGenerator int

func (e *errnoSequenceGenerator) Value() string {
	return strconv.Itoa(int(*e))
}

func (e *errnoSequenceGenerator) Next() errnoGenerator {
	*e++
	return e
}

func newErrnoSequenceGenerator() errnoGenerator {
	var start errnoSequenceGenerator = 1
	return &start
}
