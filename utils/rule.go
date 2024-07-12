package utils

import "fmt"

import "errors"

type Rule struct {}

func NewRule() *Rule {
	return &Rule{}
}

// Example core validators
func (r *Rule) ValidateRequired(field string, value string) error {
    if value == "" {
        return errors.New(field + " is required")
    }
    return nil
}

func  (r *Rule) ValidateMinLength(field string, value string, minLength int) error {
    if len(value) < minLength {
        return errors.New(field + " should have at least " + fmt.Sprint(minLength) + " characters")
    }
    return nil
}