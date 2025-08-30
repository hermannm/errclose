# errclose

A tiny Go package for handling errors when closing resources.

Run `go get hermannm.dev/errclose` to add it to your project!

**Docs:** [pkg.go.dev/hermannm.dev/errclose](https://pkg.go.dev/hermannm.dev/errclose)

**Contents:**

- [Motivation](#motivation)
- [Usage](#usage)
- [Developer's guide](#developers-guide)

## Motivation

In Go, it is easy to forget handling errors returned by calling `Close` on some resource, since one
typically calls it in a `defer` statement, like this:

<!-- @formatter:off -->
```go
func example() error {
	file, err := os.Open("/some/path")
	if err != nil {
		return err
	}
	// Close error not handled!
	defer file.Close()

	// Use file
}
```
<!-- @formatter:on -->

To properly handle the close error here, one would have to do something like this:

<!-- @formatter:off -->
```go
// Use named error return value, so we can update it in the defer
func example() (returnedErr error) {
	file, err := os.Open("/some/path")
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Check if we're already returning another error: in that case,
			// we want to combine the existing error with the close error
			if returnedErr != nil {
				returnedErr = fmt.Errorf("%w (and failed to close file: %w)", returnedErr, closeErr)
			} else {
				returnedErr = fmt.Errorf("failed to close file: %w", closeErr)
			}
		}
	}()

	// Use file
}
```
<!-- @formatter:on -->

`errclose` gives you a function to achieve this in a single line (see [Usage](#usage)).

## Usage

`errclose.Close` is the central function provided by this package. You'll typically call it in a
`defer` statement, giving it a resource (something implementing `Close() error`) and a pointer to a
named error return value. When your function exits, `errclose.Close` closes the given resource, and
if there's an error, it will update the error returned by your function to include the close error.

<!-- @formatter:off -->
```go
import (
	"os"

	"hermannm.dev/errclose"
)

// Use a named error return value, so we can pass a pointer to `errclose.Close`
func example() (returnedErr error) {
	file, err := os.Open("/some/path")
	if err != nil {
		return err
	}
	defer errclose.Close(file, &returnedErr, "file")

	// Use file
}
```
<!-- @formatter:on -->

It's recommended to give the error returned by your function a unique name (e.g. `returnedErr`), so
you don't accidentally give `errclose.Close` a pointer to a local error.

The third argument to `errclose.Close` (`"file"` in the example above) is the resource name, which
is used to format the error message in case of a close error. Errors are formatted like this:

```
failed to close <resource name>: <close error>
```

If there is already an error returned by your function when a close error is encountered (i.e.,
`returnedErr` points to a non-nil error), then the close error is combined with the existing error,
on the following format:

```
<existing error> (and failed to close <resource name>: <close error>)
```

The error formatting uses [`fmt.Errorf`](https://pkg.go.dev/fmt#Errorf) with the `%w` verb, so that
the underlying errors can still be checked with [`errors.Is`](https://pkg.go.dev/errors#Is) and
[`errors.As`](https://pkg.go.dev/errors#As).

If you want to format the resource name, you can use `errclose.Closef`, which takes a format string
and args instead of just a plain string for the resource name. The formatting is only performed if
there is a close error, so this is more efficient than calling `fmt.Sprintf` yourself and passing
the result to `errclose.Close`, as that will perform formatting even when there's no error.

<!-- @formatter:off -->
```go
func example(filePath string) (returnedErr error) {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	// If there's a close error, returnedErr will be set to:
	// failed to close file at path <filePath>: <close error>
	defer errclose.Closef(file, &returnedErr, "file at path %s", filePath)

	// Use file
}
```
<!-- @formatter:on -->

## Developer's guide

When publishing a new release:

- Run tests and linter ([`golangci-lint`](https://golangci-lint.run/)):
  ```
  go test ./... && golangci-lint run
  ```
- Add an entry to `CHANGELOG.md` (with the current date)
    - Remember to update the link section, and bump the version for the `[Unreleased]` link
- Create commit and tag for the release (update `TAG` variable in below command):
  ```
  TAG=vX.Y.Z && git commit -m "Release ${TAG}" && git tag -a "${TAG}" -m "Release ${TAG}" && git log --oneline -2
  ```
- Push the commit and tag:
  ```
  git push && git push --tags
  ```
    - Our release workflow will then create a GitHub release with the pushed tag's changelog entry
