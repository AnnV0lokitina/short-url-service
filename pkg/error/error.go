package error

import (
	"fmt"
	"strings"
)

// LabelError Store error and its description.
type LabelError struct {
	Label string
	Err   error
}

// Possible error descriptions.
const (
	TypeConflict = "CONFLICT"  // error type conflict
	TypeNotFound = "NOT FOUND" // error type not found
	TypeGone     = "GONE"      // error type gone
)

// Error Return error as string.
func (le *LabelError) Error() string {
	return fmt.Sprintf("[%s] %v", le.Label, le.Err)
}

// NewLabelError Create new LabelError.
func NewLabelError(label string, err error) error {
	return &LabelError{
		Label: strings.ToUpper(label),
		Err:   err,
	}
}
