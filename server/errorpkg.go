package server

// Provides an implementation to hold information required to handle errors within am http server
// where the server may want to sanitize errors that contain sensitive information
type ErrorPkg struct {
	innerError  error // the error being packaged
	httpStatus  int   // the status code that this error should result in
	isSanitized bool  // if true, the error message will be passed thru to the http result
}

// true if an error has occured
func (epkg *ErrorPkg) Error() bool {
	return epkg.innerError != nil
}

// true if the error is considered sanitized, and its message can be displayed in the result
func (epkg *ErrorPkg) Sanitized() ErrorPkg {
	epkg.isSanitized = true
	return *epkg
}

// create a default instance of ErrorPkg which is an unsanitized, http error 500
func Error(err error) ErrorPkg {
	return ErrorPkg{
		innerError: err,
		httpStatus: 500,
	}
}

// create an instance of ErrorPkg which is not sanitized and has the 500 (Internal Server Error) error code
func ErrorWithCode(statusCode int, err error) ErrorPkg {
	return ErrorPkg{
		innerError: err,
		httpStatus: statusCode,
	}
}

// create an instance of ErrorPkg which is sanitized and has the given error code
func ErrorWithCodeSantized(statusCode int, err error) ErrorPkg {
	return ErrorPkg{
		innerError:  err,
		httpStatus:  statusCode,
		isSanitized: true,
	}
}
