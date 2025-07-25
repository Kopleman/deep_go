package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	Errors []error
}

func (e *MultiError) Error() string {
	if e == nil || len(e.Errors) == 0 {
		return ""
	}

	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%d errors occured:\n", len(e.Errors)))

	for _, err := range e.Errors {
		b.WriteString(fmt.Sprintf("\t* %s", err.Error()))
	}

	b.WriteString("\n")
	return b.String()
}

func Append(err error, errs ...error) error {
	var validErrs []error
	for _, e := range errs {
		if e != nil {
			validErrs = append(validErrs, e)
		}
	}

	var multiError *MultiError
	if err == nil {
		multiError = &MultiError{
			Errors: validErrs,
		}

		return multiError
	}

	multiError, ok := err.(*MultiError)
	if !ok {
		multiError = &MultiError{
			Errors: append([]error{err}, validErrs...),
		}
		return multiError
	}
	multiError.Errors = append(multiError.Errors, validErrs...)

	return multiError
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}

func TestMultiErrorNil(t *testing.T) {
	var multiErr *MultiError
	assert.Equal(t, "", multiErr.Error())
}

func TestMultiErrorEmpty(t *testing.T) {
	multiErr := &MultiError{Errors: []error{}}
	assert.Equal(t, "", multiErr.Error())
}

func TestMultiErrorSingle(t *testing.T) {
	err := Append(nil, errors.New("single error"))
	expectedMessage := "1 errors occured:\n\t* single error\n"
	assert.EqualError(t, err, expectedMessage)
}

func TestMultiErrorMultiple(t *testing.T) {
	err := Append(nil,
		errors.New("first error"),
		errors.New("second error"),
		errors.New("third error"),
	)
	expectedMessage := "3 errors occured:\n\t* first error\t* second error\t* third error\n"
	assert.EqualError(t, err, expectedMessage)
}

func TestAppendToRegularError(t *testing.T) {
	baseErr := errors.New("base error")
	err := Append(baseErr, errors.New("additional error"))
	expectedMessage := "2 errors occured:\n\t* base error\t* additional error\n"
	assert.EqualError(t, err, expectedMessage)
}

func TestAppendToMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))
	err = Append(err, errors.New("error 3"))

	expectedMessage := "3 errors occured:\n\t* error 1\t* error 2\t* error 3\n"
	assert.EqualError(t, err, expectedMessage)
}

func TestAppendEmptyErrors(t *testing.T) {
	baseErr := errors.New("base error")
	multiErr := Append(baseErr)
	expectedMessage := "1 errors occured:\n\t* base error\n"
	assert.EqualError(t, multiErr, expectedMessage)
}

func TestAppendToNil(t *testing.T) {
	err := Append(nil, errors.New("new error"))
	expectedMessage := "1 errors occured:\n\t* new error\n"
	assert.EqualError(t, err, expectedMessage)
}

func TestAppendReturnType(t *testing.T) {
	err := Append(nil, errors.New("test error"))
	assert.IsType(t, &MultiError{}, err)
}

func TestNestedMultiError(t *testing.T) {
	innerErr := Append(nil, errors.New("inner 1"), errors.New("inner 2"))
	outerErr := Append(innerErr, errors.New("outer 1"))

	expectedMessage := "3 errors occured:\n\t* inner 1\t* inner 2\t* outer 1\n"
	assert.EqualError(t, outerErr, expectedMessage)
}

func TestMultiErrorCount(t *testing.T) {
	err := Append(nil,
		errors.New("error 1"),
		errors.New("error 2"),
		errors.New("error 3"),
	)

	multiErr, ok := err.(*MultiError)
	assert.True(t, ok)
	assert.Equal(t, 3, len(multiErr.Errors))
}

func TestAppendNilErrors(t *testing.T) {
	err := Append(nil, nil, errors.New("valid error"), nil)
	expectedMessage := "1 errors occured:\n\t* valid error\n"
	assert.EqualError(t, err, expectedMessage)
}

func TestImmutability(t *testing.T) {
	originalErr := errors.New("original")
	appendedErr := Append(originalErr, errors.New("appended"))

	assert.EqualError(t, originalErr, "original")

	expectedMessage := "2 errors occured:\n\t* original\t* appended\n"
	assert.EqualError(t, appendedErr, expectedMessage)
}

func TestAppendAllNilErrors(t *testing.T) {
	err := Append(nil, nil, nil, nil)
	assert.Equal(t, "", err.Error())
}

func TestAppendMixedNilErrors(t *testing.T) {
	err := Append(nil, nil, errors.New("valid 1"), nil, errors.New("valid 2"), nil)
	expectedMessage := "2 errors occured:\n\t* valid 1\t* valid 2\n"
	assert.EqualError(t, err, expectedMessage)
}

func TestAppendToNilWithNilErrors(t *testing.T) {
	err := Append(nil, nil)
	assert.Equal(t, "", err.Error())
}

func TestAppendNilToRegularError(t *testing.T) {
	baseErr := errors.New("base error")
	err := Append(baseErr, nil)
	expectedMessage := "1 errors occured:\n\t* base error\n"
	assert.EqualError(t, err, expectedMessage)
}

func TestMultiErrorUniqueness(t *testing.T) {
	err1 := errors.New("same error")
	err2 := errors.New("same error")

	err := Append(nil, err1, err2)
	expectedMessage := "2 errors occured:\n\t* same error\t* same error\n"
	assert.EqualError(t, err, expectedMessage)
}

type CustomError struct {
	Message string
}

func (e CustomError) Error() string {
	return e.Message
}

func TestMultiErrorWithCustomErrors(t *testing.T) {
	customErr1 := CustomError{Message: "custom error 1"}
	customErr2 := CustomError{Message: "custom error 2"}

	err := Append(nil, customErr1, customErr2)
	expectedMessage := "2 errors occured:\n\t* custom error 1\t* custom error 2\n"
	assert.EqualError(t, err, expectedMessage)
}

func TestMultiErrorWithSpecialCharacters(t *testing.T) {
	err := Append(nil,
		errors.New("error with \"quotes\""),
		errors.New("error with\nnewline"),
		errors.New("error with\ttab"),
	)

	expectedMessage := "3 errors occured:\n\t* error with \"quotes\"\t* error with\nnewline\t* error with\ttab\n"
	assert.EqualError(t, err, expectedMessage)
}

func TestMultiErrorWithLongMessages(t *testing.T) {
	longMessage := strings.Repeat("very long error message ", 50)
	err := Append(nil, errors.New(longMessage))

	assert.Contains(t, err.Error(), longMessage)
	assert.Contains(t, err.Error(), "1 errors occured:")
}

func TestMultiErrorWithEmptyMessages(t *testing.T) {
	err := Append(nil, errors.New(""), errors.New("normal error"), errors.New(""))
	expectedMessage := "3 errors occured:\n\t* \t* normal error\t* \n"
	assert.EqualError(t, err, expectedMessage)
}
