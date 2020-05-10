package postform

import (
	"bytes"
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"github.com/google/uuid"
)

var UUID_TYPE = reflect.ValueOf(uuid.UUID{}).Type()
var BYTE_SLICE_TYPE = reflect.ValueOf([]byte{}).Type()

func Decode(dst interface{}, src *http.Request) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("postform: interface must be a pointer to struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("postform")
		if tag == "" {
			continue
		}

		formString := src.PostFormValue(tag)
		if formString != "" {
			if field.Type.Kind() == reflect.String {
				v.FieldByName(field.Name).SetString(formString)
			} else if field.Type.Kind() == reflect.Int {
				i, err := strconv.Atoi(formString)
				if err != nil {
					return err
				}
				v.FieldByName(field.Name).SetInt(int64(i))
			} else if field.Type.Kind() == reflect.Float64 {
				f, err := strconv.ParseFloat(formString, 64)
				if err != nil {
					return err
				}
				v.FieldByName(field.Name).SetFloat(f)
			} else if field.Type == UUID_TYPE {
				u, err := uuid.Parse(formString)
				if err != nil {
					return err
				}
				v.FieldByName(field.Name).Set(reflect.ValueOf(u))
			} else if field.Type.Kind() == reflect.Bool {
				v.FieldByName(field.Name).SetBool(formString == "true")
			}
		}

		if field.Type == BYTE_SLICE_TYPE {
			file, _, err := src.FormFile(tag)
			if err != nil && err != http.ErrMissingFile {
				return err
			} else if err == nil {
				defer file.Close()
				buffer := bytes.NewBuffer(nil)
				if _, err := buffer.ReadFrom(file); err != nil {
					return err
				}
				v.FieldByName(field.Name).Set(reflect.ValueOf(buffer.Bytes()))
			}
		}
	}

	return nil
}
