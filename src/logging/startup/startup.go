package startup

import (
	"errors"
	"fmt"
)

func NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{errors: []error{}}
}

type ErrorCollector struct {
	errors []error
}

func (ec *ErrorCollector) Append(err error) {
	if err != nil {
		ec.errors = append(ec.errors, err)
	}
}

func (ec *ErrorCollector) Error(msg string) {
	ec.Append(errors.New(msg))
}

func (ec *ErrorCollector) Print() {
	for i := 0; i < len(ec.errors); i++ {
		fmt.Println((ec.errors[i]).Error())
	}
}

func (ec *ErrorCollector) IsEmpty() bool {
	return len(ec.errors) == 0
}
