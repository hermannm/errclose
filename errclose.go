// Package errclose provides a single function, [errclose.Close], for handling errors when closing
// resources.
package errclose

import (
	"fmt"
)

// Close closes the given resource, and handles close errors.
//
// You'll typically call this in a defer statement (to close a resource when the function exits),
// using named returns to give a pointer to the error returned by your function:
//
//	import (
//		"os"
//
//		"hermannm.dev/errclose"
//	)
//
//	// Use a named error return value, so we can pass a pointer to errclose.Close
//	func example() (returnedErr error) {
//		file, err := os.Open("/some/path")
//		if err != nil {
//			return err
//		}
//		defer errclose.Close(file, &returnedErr, "file")
//
//		// Use file
//	}
//
// It's recommended to give the error returned by your function a unique name (e.g. 'returnedErr'),
// so you don't accidentally give errclose.Close a pointer to a local error.
//
// # Error format
//
// If resource.Close returns an error, then we set the error pointed to by returnedErr to the close
// error. The close error is wrapped with the given resource name for context, on the following
// format:
//
//	failed to close <resourceName>: <close error>
//
// If returnedErr points to an existing non-nil error, then we combine the existing error and the
// close error on the following format:
//
//	<existing error> (and failed to close <resourceName>: <close error>)
//
// We use the %w verb with [fmt.Errorf] when formatting the error string, so that the underlying
// errors can be checked with [errors.Is] and [errors.As].
func Close(
	resource interface{ Close() error },
	returnedErr *error,
	resourceName string,
) {
	closeErr := resource.Close()
	if closeErr == nil {
		return
	}

	currentReturnedErr := *returnedErr
	if currentReturnedErr != nil {
		*returnedErr = fmt.Errorf(
			"%w (and failed to close %s: %w)",
			currentReturnedErr,
			resourceName,
			closeErr,
		)
	} else {
		*returnedErr = fmt.Errorf("failed to close %s: %w", resourceName, closeErr)
	}
}
