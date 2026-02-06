package apperr

import "fmt"

// Error codes
const (
	// 400xx - Client errors
	CodeInvalidInput  = 40001
	CodeMissingField  = 40002
	CodeInvalidFormat = 40003
	CodeInvalidState  = 40004

	// 401xx - Authentication
	CodeUnauthorized  = 40101
	CodeTokenExpired  = 40102
	CodeTokenInvalid  = 40103
	CodeLoginFailed   = 40104
	CodePasswordWrong = 40105

	// 403xx - Authorization
	CodeForbidden        = 40301
	CodePermissionDenied = 40302

	// 404xx - Not found
	CodeNotFound     = 40401
	CodeUserNotFound = 40402

	// 409xx - Conflict
	CodeAlreadyExists = 40901
	CodeHasRelation   = 40902

	// 500xx - Server errors
	CodeInternal = 50001
	CodeDBError  = 50002

	// 503xx - Service unavailable
	CodeServiceUnavailable = 50301
	CodeMediaMTXUnavail    = 50302
)

// Error is the unified application error type.
type Error struct {
	Code    int
	Message string
	Cause   error
	Details map[string]interface{}
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func New(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

func Wrap(cause error, code int, message string) *Error {
	return &Error{Code: code, Message: message, Cause: cause}
}

func WithDetails(err *Error, details map[string]interface{}) *Error {
	e := *err
	e.Details = details
	return &e
}

func NotFound(entity string, id interface{}) *Error {
	return &Error{
		Code:    CodeNotFound,
		Message: fmt.Sprintf("%s not found: %v", entity, id),
	}
}

func InvalidInput(message string) *Error {
	return &Error{Code: CodeInvalidInput, Message: message}
}

func Conflict(message string) *Error {
	return &Error{Code: CodeAlreadyExists, Message: message}
}

func HasRelation(message string) *Error {
	return &Error{Code: CodeHasRelation, Message: message}
}

func Unauthorized(message string) *Error {
	return &Error{Code: CodeUnauthorized, Message: message}
}

func Forbidden(message string) *Error {
	return &Error{Code: CodeForbidden, Message: message}
}

func Internal(message string, cause error) *Error {
	return &Error{Code: CodeInternal, Message: message, Cause: cause}
}

func ServiceUnavailable(message string) *Error {
	return &Error{Code: CodeServiceUnavailable, Message: message}
}

func IsNotFound(err error) bool {
	e, ok := err.(*Error)
	return ok && e.Code == CodeNotFound
}

func IsConflict(err error) bool {
	e, ok := err.(*Error)
	return ok && (e.Code == CodeAlreadyExists || e.Code == CodeHasRelation)
}

func IsUnauthorized(err error) bool {
	e, ok := err.(*Error)
	return ok && (e.Code >= 40101 && e.Code <= 40199)
}

func IsForbidden(err error) bool {
	e, ok := err.(*Error)
	return ok && (e.Code >= 40301 && e.Code <= 40399)
}
