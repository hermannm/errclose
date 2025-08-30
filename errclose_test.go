package errclose_test

import (
	"errors"
	"reflect"
	"testing"

	"hermannm.dev/errclose"
)

func TestCloseError(t *testing.T) {
	var file *mockFile

	useFile := func() (returnedErr error) {
		file = openFileWithCloseError()
		defer errclose.Close(file, &returnedErr, "file")

		return nil
	}

	err := useFile()
	assertEqual(t, file.closeWasCalled, true, "file.closeWasCalled")
	assertEqual(t, err.Error(), "failed to close file: close error", "error string")
	assertEqual(t, errors.Is(err, file.closeError), true, "errors.Is result")
}

func TestCloseErrorWithExistingError(t *testing.T) {
	var file *mockFile

	useFile := func() (returnedErr error) {
		file = openFileWithCloseError()
		defer errclose.Close(file, &returnedErr, "file")

		return fallibleOperation()
	}

	err := useFile()
	assertEqual(t, file.closeWasCalled, true, "file.closeWasCalled")
	assertEqual(
		t,
		err.Error(),
		"operation failed (and failed to close file: close error)",
		"error string",
	)
	assertEqual(t, errors.Is(err, file.closeError), true, "errors.Is(closeError)")
	assertEqual(t, errors.Is(err, errFallibleOperation), true, "errors.Is(errFallibleOperation)")
}

func TestNilCloseError(t *testing.T) {
	var file *mockFile

	useFile := func() (returnedErr error) {
		file = openFileWithoutCloseError()
		defer errclose.Close(file, &returnedErr, "file")

		return nil
	}

	err := useFile()
	assertEqual(t, file.closeWasCalled, true, "file.closeWasCalled")
	assertEqual(t, err, nil, "error")
}

func TestNilCloseErrorWithExistingError(t *testing.T) {
	var file *mockFile

	useFile := func() (returnedErr error) {
		file = openFileWithoutCloseError()
		defer errclose.Close(file, &returnedErr, "file")

		return fallibleOperation()
	}

	err := useFile()
	assertEqual(t, file.closeWasCalled, true, "file.closeWasCalled")
	assertEqual(t, err, errFallibleOperation, "error")
}

func TestClosef(t *testing.T) {
	var file *mockFile

	useFile := func() (returnedErr error) {
		file = openFileWithCloseError()
		defer errclose.Closef(file, &returnedErr, "file at path %s", "/example/path")

		return nil
	}

	err := useFile()
	assertEqual(t, file.closeWasCalled, true, "file.closeWasCalled")
	assertEqual(
		t,
		err.Error(),
		"failed to close file at path /example/path: close error",
		"error string",
	)
	assertEqual(t, errors.Is(err, file.closeError), true, "errors.Is result")
}

func TestClosefWithExistingError(t *testing.T) {
	var file *mockFile

	useFile := func() (returnedErr error) {
		file = openFileWithCloseError()
		defer errclose.Closef(file, &returnedErr, "file at path %s", "/example/path")

		return fallibleOperation()
	}

	err := useFile()
	assertEqual(t, file.closeWasCalled, true, "file.closeWasCalled")
	assertEqual(
		t,
		err.Error(),
		"operation failed (and failed to close file at path /example/path: close error)",
		"error string",
	)
	assertEqual(t, errors.Is(err, file.closeError), true, "errors.Is(closeError)")
	assertEqual(t, errors.Is(err, errFallibleOperation), true, "errors.Is(errFallibleOperation)")
}

func TestClosefWithoutCloseError(t *testing.T) {
	var file *mockFile

	useFile := func() (returnedErr error) {
		file = openFileWithoutCloseError()
		defer errclose.Closef(file, &returnedErr, "file at path %s", "/example/path")

		return nil
	}

	err := useFile()
	assertEqual(t, file.closeWasCalled, true, "file.closeWasCalled")
	assertEqual(t, err, nil, "error")
}

func TestClosefWithoutCloseErrorWithExistingError(t *testing.T) {
	var file *mockFile

	useFile := func() (returnedErr error) {
		file = openFileWithoutCloseError()
		defer errclose.Closef(file, &returnedErr, "file at path %s", "/example/path")

		return fallibleOperation()
	}

	err := useFile()
	assertEqual(t, file.closeWasCalled, true, "file.closeWasCalled")
	assertEqual(t, err, errFallibleOperation, "error")
}

type mockFile struct {
	closeWasCalled bool
	closeError     error
}

func (file *mockFile) Close() error {
	file.closeWasCalled = true
	return file.closeError
}

func openFileWithCloseError() *mockFile {
	return &mockFile{closeWasCalled: false, closeError: errors.New("close error")}
}

func openFileWithoutCloseError() *mockFile {
	return &mockFile{closeWasCalled: false, closeError: nil}
}

func fallibleOperation() error {
	return errFallibleOperation
}

var errFallibleOperation = errors.New("operation failed")

func assertEqual(t *testing.T, actual any, expected any, descriptor string) {
	t.Helper()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf(
			`Unexpected %s
Want: %+v
 Got: %+v`,
			descriptor,
			expected,
			actual,
		)
	}
}
