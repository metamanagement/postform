package postform

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestBasic(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	err := writer.WriteField("f1", "243")
	if err != nil {
		t.Errorf(err.Error())
	}
	err = writer.WriteField("f2", "abcdef")
	if err != nil {
		t.Errorf(err.Error())
	}
	err = writer.WriteField("f5", "2433")
	if err != nil {
		t.Errorf(err.Error())
	}
	err = writer.WriteField("f6", "00D67575-172D-4DA6-BCAB-F9796EA84D66")
	if err != nil {
		t.Errorf(err.Error())
	}

	testFile, err := ioutil.ReadFile("./decoder_test.go")
	if err != nil {
		t.Errorf(err.Error())
	}
	w, err := writer.CreateFormFile("f3", "arbitrary_file.txt")
	if err != nil {
		t.Errorf(err.Error())
	}
	_, err = w.Write(testFile)
	if err != nil {
		t.Errorf(err.Error())
	}
	writer.Close()

	req, err := http.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	type basic struct {
		Field1 int    `postform:"f1"`
		Field2 string `postform:"f2"`
		Field3 []byte `postform:"f3"`
		Field4 string
		Field5 float64   `postform:"f5"`
		Field6 uuid.UUID `postform:"f6"`
	}

	var b basic
	err = Decode(&b, req)
	if err != nil {
		t.Errorf(err.Error())
	}

	if b.Field1 != 243 {
		t.Errorf("Field1 failed to parse int correctly.")
	}

	if b.Field2 != "abcdef" {
		t.Errorf("Field2 failed to parse string correctly.")
	}

	if len(b.Field3) != len(testFile) {
		t.Errorf("Field3 file size does not match.")
	}

	if b.Field4 != "" {
		t.Errorf("Erroneously parsed info into Field4.")
	}

	if b.Field5 != 2433.0 {
		t.Errorf("Field5 failed to parse float correctly.")
	}

	if b.Field6.String() != strings.ToLower("00D67575-172D-4DA6-BCAB-F9796EA84D66") {
		t.Errorf("Field6 failed to parse UUID correctly.")
	}
}
