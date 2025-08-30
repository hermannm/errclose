// Package errclose provides [errclose.Close], a function for handling errors when closing
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
// so you don't accidentally give a pointer to a local error.
//
// # Error format
//
// If resource.Close returns an error, then the error pointed to by returnedErr is set to the close
// error. The close error is wrapped with the given resource name for context, on the following
// format:
//
//	failed to close <resourceName>: <close error>
//
// If returnedErr points to an existing non-nil error, then the existing error and the close error
// are combined on the following format:
//
//	<existing error> (and failed to close <resourceName>: <close error>)
//
// The error string formatting uses [fmt.Errorf] with the %w verb, so that the underlying errors can
// be checked with [errors.Is] and [errors.As].
//
// If you want to use format args to format the resource name, call [errclose.Closef].
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

// Closef closes the given resource, and handles close errors.
//
// It takes a format string and args to construct a name for the resource (with [fmt.Sprintf]),
// which is added to the close error for context (see 'Error format' below). The formatting is only
// performed if there is a close error, so this is more efficient than calling [fmt.Sprintf]
// yourself and passing the result to [errclose.Close], as that will perform formatting even when
// there's no error.
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
//	// Use a named error return value, so we can pass a pointer to errclose.Closef
//	func example(filePath string) (returnedErr error) {
//		file, err := os.Open(filePath)
//		if err != nil {
//			return err
//		}
//		defer errclose.Closef(file, &returnedErr, "file at path %s", filePath)
//
//		// Use file
//	}
//
// It's recommended to give the error returned by your function a unique name (e.g. 'returnedErr'),
// so you don't accidentally give a pointer to a local error.
//
// # Error format
//
// If resource.Close returns an error, then the error pointed to by returnedErr is set to the close
// error. The close error is wrapped with a resource name for context, using the given resource name
// format string and args to format the name. The complete error string looks like this:
//
//	failed to close <formatted resource name>: <close error>
//
// If returnedErr points to an existing non-nil error, then the existing error and the close error
// are combined on the following format:
//
//	<existing error> (and failed to close <formatted resource name>: <close error>)
//
// The error string formatting uses [fmt.Errorf] with the %w verb, so that the underlying errors can
// be checked with [errors.Is] and [errors.As].
func Closef(
	resource interface{ Close() error },
	returnedErr *error,
	resourceNameFormat string,
	formatArgs ...any,
) {
	closeErr := resource.Close()
	if closeErr == nil {
		return
	}

	resourceName := fmt.Sprintf(resourceNameFormat, formatArgs...)

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
