// Package to hold all cross cutting common errors across and utils to deal with errors
package errors

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	temporalsdk "go.temporal.io/sdk/temporal"

	"gorm.io/gorm"
)

type ErrorType string

const (
	/*
		List of common error types that are valid irrespective of the domain
	*/

	ErrRecordNotFoundStr             ErrorType = "record not found"
	ErrInvalidArgumentStr            ErrorType = "Invalid Argument"
	ErrBadRequestStr                 ErrorType = "Bad Request"
	ErrAlreadyExistsStr              ErrorType = "record already exists"
	ErrInternalServerStr             ErrorType = "internal server error"
	ErrInProgressStr                 ErrorType = "in progress error"
	ErrPermissionDeniedStr           ErrorType = "permission denied"
	ErrFailedPreconditionStr         ErrorType = "failed precondition"
	ErrActivityRateLimitExhaustedStr ErrorType = "Activity Rate Limit Exhausted"
	ErrResourceExhaustedStr          ErrorType = "resource exhausted"
)

var (
	ErrInvalidEnvironment = errors.New("invalid environment")
	ErrEnvironmentNotSet  = errors.New("environment not set")
	// Signifies retryable failure
	ErrTransient = errors.New("transient failure")
	// Signifies non-retryable failure
	ErrPermanent = errors.New("permanent failure")
	// ErrRequestCanceled signifies canceled request
	// eg: workflow cancellation signal to workflow
	ErrRequestCanceled    = errors.New("request canceled")
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrAlreadyExists      = errors.New("already exists")
	ErrRecordNotFound     = errors.New("record not found")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrFailedPrecondition = errors.New("failed precondition")
	ErrTimedOut           = errors.New("timed out")
	ErrResourceExhausted  = errors.New("resource exhausted")
)

var InvalidEnvironmentErrFn = func(env string) error {
	return fmt.Errorf("unexpected env: %s: %w", env, ErrInvalidEnvironment)
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, format string, a ...any) error {
	return errors.Wrap(err, fmt.Sprintf(format, a...))
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// UnwrapRootCause unwraps the error until the root cause is identified
func UnwrapRootCause(err error) error {
	for {
		// Check if the error has an Unwrap method
		unwrapErr, ok := err.(interface {
			Unwrap() error
		})
		if !ok {
			// If the error does not have an Unwrap method, return it
			return err
		}

		// Unwrap the error
		nextErr := unwrapErr.Unwrap()
		if nextErr == nil {
			// If the unwrapped error is nil, return the current error
			return err
		}

		// Set the next error as the current error for the next iteration
		err = nextErr
	}
}

// UnwrapLastN unwraps the error until the last Nth error is identified
func UnwrapLastN(err error, n int) error {
	if err == nil {
		return nil
	}

	if n == 0 {
		return UnwrapRootCause(err)
	}

	var errorChain []error
	errorChain = append(errorChain, err)

	for err != nil {
		err = Unwrap(err)
		if err == nil {
			continue
		}

		errorChain = append(errorChain, err)
	}

	if n >= len(errorChain) {
		return errorChain[0]
	}

	return errorChain[len(errorChain)-n]
}

// New returns an error with the supplied message.
func New(format string, a ...any) error {
	return errors.New(fmt.Sprintf(format, a...))
}

// Errorf returns an error with the supplied message.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Is returns true if err is of type target
func Is(err, target error) bool {
	return errors.Is(err, target)
}

func IsRecordNotFound(err error) bool {
	return Is(err, gorm.ErrRecordNotFound) || Is(err, ErrRecordNotFound)
}

func IsDuplicateKeyConstraintErr(err error) bool {
	return Is(err, ErrAlreadyExists) || strings.EqualFold(err.Error(), "duplicate key value violates unique constraint")
}

// IsErrorOfType checks whether an error is of particular type or not
// this can be typically used in the workflow code where workflow business logic on activity failure
// needs more granular information than just retryable and non-retryable failure information
func IsErrorOfType(err error, errType ErrorType) bool {
	if err == nil {
		return false
	}

	appErr := &temporalsdk.ApplicationError{}
	if ok := errors.As(err, &appErr); ok {
		return appErr.Type() == string(errType)
	}

	errResp := &ErrorResponse{}
	if ok := errors.As(err, &errResp); ok {
		return errResp.ErrorType == errType
	}

	return false
}
