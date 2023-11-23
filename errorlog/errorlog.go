package errorlog

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
)

type AppError struct {
	Err     error
	Code    int
	Message string
}

func (e AppError) Error() string {
	return e.Err.Error()
}

type HTTPErrorResponseWriterAndLogger struct {
	Logger *log.Logger
}

const (
	runtimeCallerDepth = 2
	outputCallDepth    = 3
)

func (h *HTTPErrorResponseWriterAndLogger) LogError(w io.Writer, err AppError) {
	// Get file and line information for the caller
	_, file, line, ok := runtime.Caller(runtimeCallerDepth)

	// Prepare a concise error message with filename and line number
	var trace string
	if ok {
		trace = fmt.Sprintf("Error in %s:%d - %s", file, line, err.Error())
	} else {
		trace = fmt.Sprintf("Error: %s", err.Error())
	}

	// Log the concise error message
	h.Logger.Output(outputCallDepth, trace)
}

func (h *HTTPErrorResponseWriterAndLogger) HandleError(w http.ResponseWriter, err AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)

	errorResponse := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    err.Code,
			"message": err.Message,
		},
	}

	json.NewEncoder(w).Encode(errorResponse)
}

func (h *HTTPErrorResponseWriterAndLogger) LogAndHandleError(wlog io.Writer, wh http.ResponseWriter, err AppError) {

	h.LogError(wlog, err)
	h.HandleError(wh, err)
}

// Standard error messages

func BadRequestError(err error) AppError {
	return AppError{
		Err:     err,
		Code:    http.StatusBadRequest,
		Message: "bad request. please check your input.",
	}
}

func InternalServerError(err error) AppError {
	return AppError{
		Err:     err,
		Code:    http.StatusInternalServerError,
		Message: "internal server error. please try again later.",
	}
}

func NotFoundError(err error) AppError {
	return AppError{
		Err:     err,
		Code:    http.StatusNotFound,
		Message: "resource not found. please check the URL.",
	}
}

func UnauthorizedError(err error) AppError {
	return AppError{
		Err:     err,
		Code:    http.StatusUnauthorized,
		Message: "unauthorized access. please provide valid credentials.",
	}
}

func ForbiddenError(err error) AppError {
	return AppError{
		Err:     err,
		Code:    http.StatusForbidden,
		Message: "forbidden access. you don't have permission to access this resource.",
	}
}
