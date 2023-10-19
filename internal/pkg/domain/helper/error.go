package helper

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/uptrace/bunrouter"
)

// ------------ Use these errors below for USECASE ------------

func NewUnauthorizedError(format string, a ...any) StatusError {
	internalCode, format := findInternalCode(format)
	return StatusError{
		InternalCode: internalCode,
		HTTPCode:     http.StatusUnauthorized,
		Err:          fmt.Errorf(format, a...),
	}
}

func NewBadRequestError(format string, a ...any) StatusError {
	internalCode, format := findInternalCode(format)
	return StatusError{
		InternalCode: internalCode,
		HTTPCode:     http.StatusBadRequest,
		Err:          fmt.Errorf(format, a...),
	}
}

func NewConflictError(format string, a ...any) StatusError {
	internalCode, format := findInternalCode(format)
	return StatusError{
		InternalCode: internalCode,
		HTTPCode:     http.StatusConflict,
		Err:          fmt.Errorf(format, a...),
	}
}

type StatusError struct {
	HTTPCode     int
	InternalCode string
	Err          error
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

// ------------ Use these errors below for HANDLER ------------

// findInternalCode finds a parseable int before a ":" char,
// example: "007: unable to find something: %v" will return
// internalCode = "007" and format = "unable to find something: %v"
func findInternalCode(s string) (internalCode, format string) {
	internalCode = "007" // this is an arbitrary number, inspired by james bond wkwk

	strstr := strings.Split(s, ":")
	for i, str := range strstr {
		_, err := strconv.Atoi(str)
		if i == 0 && err == nil {
			internalCode = str
			continue
		}

		if format == "" {
			format = strings.TrimLeft(str, " ")
			continue
		}

		format = fmt.Sprintf("%v:%v", format, str)
	}

	return // don't be surprised for this naked return bcs it uses named return value ;)
	// but make sure to avoid naked returns + named return values in this code base bcs it's confusing
}

// NewErrRes is a shorthand for w.WriteHeader(); bunrouter.JSON().
// The variadic params will only take the first two values if provided,
// an error will be placed at message and a string will be the the internal code.
// PLEASE USE IT WISELY!
func NewErrRes(w http.ResponseWriter, httpStatus int, params ...any) {
	res := bunrouter.H{ // arbitrary default response body
		"code":    "007", // this is an arbitrary number, inspired by james bond wkwk
		"message": "Uh oh no, something went wrong :(",
	}

	for _, p := range params {
		switch v := p.(type) {
		case string:
			res["code"] = p
		case error:
			res["message"] = v.Error()
		default:
			log.Println("THE HELL ARE YOU DOING?!")
		}
	}

	w.WriteHeader(httpStatus)
	bunrouter.JSON(w, res)
}
