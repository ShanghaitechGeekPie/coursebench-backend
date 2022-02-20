package events

import "strconv"

type eventNoGenerator interface {
	Value() string
	Next() eventNoGenerator
}

type eventNoSequenceGenerator int

func (e *eventNoSequenceGenerator) Value() string {
	return strconv.Itoa(int(*e))
}

func (e *eventNoSequenceGenerator) Next() eventNoGenerator {
	*e++
	return e
}

func newErrnoSequenceGenerator() eventNoGenerator {
	var start eventNoSequenceGenerator = 1
	return &start
}
