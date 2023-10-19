package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/uptrace/bunrouter"
)

// ShouldBindQuery is a shortcut for bunrouter.ParamsFromContext(); params.ByName() on passed obj every field with tag "uri".
// For now, those fields with tag "uri" must be either string or int.
// PLEASE ONLY PASS POINTER OF STRUCT!
func ShouldBindUri(obj any, bunReq bunrouter.Request) error {
	rv := reflect.ValueOf(obj)
	if rv.Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}

	t := rv.Type()
	params := bunrouter.ParamsFromContext(bunReq.Context())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		val, ok := field.Tag.Lookup("uri")
		if !ok {
			continue
		}

		extractedParam := params.ByName(val)

		rvField := rv.FieldByName(field.Name)
		switch rvField.Kind() {
		case reflect.Int:
			parsedInt, err := strconv.Atoi(extractedParam)
			if err != nil {
				return err
			}
			rvField.Set(reflect.ValueOf(parsedInt))
		case reflect.String:
			rvField.Set(reflect.ValueOf(extractedParam))
		default:
			return fmt.Errorf("can not set %v", rvField.Kind())
		}
	}

	return nil
}

// ShouldBindQuery is a shortcut for bunReq.URL.Query().Get() on passed obj every field with tag "form".
// If you want to parse multipart/form-data, then just use bunReq.ParseMultipartForm(); bunReq.MultipartForm.
// PLEASE ONLY PASS POINTER OF STRUCT!
func ShouldBindQuery(obj any, bunReq bunrouter.Request) error {
	rv := reflect.ValueOf(obj)
	if rv.Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}

	t := rv.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		val, ok := field.Tag.Lookup("form")
		if !ok {
			continue
		}

		rv.FieldByName(field.Name).Set(reflect.ValueOf(
			bunReq.URL.Query().Get(val),
		))
	}

	return nil
}

// ShouldBindJSON is a shortcut io.ReadAll(); json.Unmarshal().
// PLEASE ONLY PASS POINTER OF STRUCT!
func ShouldBindJSON(obj any, bunReq bunrouter.Request) error {
	body, err := io.ReadAll(bunReq.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, obj)
	if err != nil {
		return err
	}

	return nil
}
