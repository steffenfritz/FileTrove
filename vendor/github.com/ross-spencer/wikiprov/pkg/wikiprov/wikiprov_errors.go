package wikiprov

import (
	"fmt"
)

// ResponseError defines an error type that can be inspected by callers
// of spargo.
type ResponseError struct {
	expectedCode int
	receivedCode int
	Err          error
}

// makeError is an internal function that allows spargo to construct
// meaningful errors and return that error to the caller.
func (err ResponseError) makeError(expectedCode int, receivedCode int) error {
	err.expectedCode = expectedCode
	err.receivedCode = receivedCode
	return err
}

// Error enables ResponseError to implement the Errors interface.
func (err ResponseError) Error() string {
	return fmt.Sprintf("wikiprov: unexpected response from server: %d", err.receivedCode)
}
