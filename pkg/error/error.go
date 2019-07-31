package error

import "fmt"

type errType = int

const (
	errInvalidProperty errType = iota + 1
	errInvalidDocument
	errUnsupportedAsyncapiVersion
)

// Error represents conversion error
type Error struct {
	errType
	msg string
}

func (err Error) Error() string {
	return err.msg
}

func isErrorType(errType errType, err error) bool {
	if err, ok := err.(Error); ok {
		return err.errType == errType
	}
	return false
}

// IsInvalidProperty returns true if err is InvalidProperty error,
// otherwise returns false.
//
// See NewInvalidProperty
func IsInvalidProperty(err error) bool {
	return isErrorType(errInvalidProperty, err)
}

// IsInvalidDocument returns true if err is InvalidDocument error,
// otherwise returns false.
//
// See NewInvalidDocument
func IsInvalidDocument(err error) bool {
	return isErrorType(errInvalidDocument, err)
}

// IsUnsupportedAsyncapiVersion returns true if err is UnsupportedAsyncapiVersion error,
// otherwise returns false.
//
// See UnsupportedAsyncapiVersion
func IsUnsupportedAsyncapiVersion(err error) bool {
	return isErrorType(errUnsupportedAsyncapiVersion, err)
}

func newError(errType errType, msg string) Error {
	return Error{
		errType: errType,
		msg:     msg,
	}
}

// NewInvalidProperty creates new invalid property error.
// This error is returned by converter when one of the document
// properties is invalid or it is missing.
func NewInvalidProperty(context interface{}) Error {
	msg := fmt.Sprintf("asyncapi: error invalid property %v", context)
	return newError(errInvalidProperty, msg)
}

// NewInvalidDocument creates new invalid document error.
// This error is returned by converter when the document has invalid structure.
func NewInvalidDocument() Error {
	return newError(errInvalidDocument, "asyncapi: unable to decode document")
}

// NewUnsupportedAsyncapiVersion creates new unsupported asyncapi version error.
// This error is returned when converter does not recognize the version of the
// converted AsyncAPI document.
func NewUnsupportedAsyncapiVersion(context interface{}) Error {
	msg := fmt.Sprintf("asyncapi: unsupported asyncapi version '%v'", context)
	return newError(errUnsupportedAsyncapiVersion, msg)
}
