// Package errorso provides common error declarations
package errorso

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
)
