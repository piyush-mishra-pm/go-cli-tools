package main

import (
	"errors"
	"fmt"
)

var (
	ErrValidation = errors.New("Validation failed")
)

type errStep struct {
	step  string
	msg   string
	cause error
}

func (s *errStep) Error() string {
	return fmt.Sprintf("Step: %q: %s: Cause: %v", s.step, s.msg, s.cause)
}

func (s *errStep) Is(target error) bool {
	t, ok := target.(*errStep)
	if !ok {
		return false
	}

	return t.step == s.step
}

func (s *errStep) Unwrap() error {
	return s.cause
}
