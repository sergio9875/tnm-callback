package utils

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
)

const (
	// ContentTypeApplicationANY const
	ContentTypeApplicationANY = "*/*"
	// ContentTypeApplicationJSON const
	ContentTypeApplicationJSON = "application/json"
	// ContentTypeApplicationXML const
	ContentTypeApplicationXML = "application/xml"
	ContentType               = "Content-Type"
)

type GenericResponse struct {
	XMLName           *xml.Name `xml:"API3G" json:",omitempty"`
	Result            string    `xml:"Result"`
	ResultExplanation string    `xml:"ResultExplanation"`
}

// GetRequestBody return array of byte from the request
func GetRequestBody(request *http.Request) ([]byte, error) {
	if request.ContentLength == 0 {
		return nil, errors.New("no body found in request")
	}
	rs, _ := ioutil.ReadAll(request.Body)
	_ = request.Body.Close()
	request.Body = ioutil.NopCloser(bytes.NewBuffer(rs))

	return rs, nil
}

// FileExists func
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// ErrorToText
func ErrorToText(err GenericResponse) string {
	return fmt.Sprintf("%s;%s;;;", err.Result, err.ResultExplanation)
}

// ResponseWriterTextHtml function
func ResponseWriterTextHtml(writer http.ResponseWriter, response interface{}) {
	writer.Header().Set(ContentType, ContentTypeApplicationXML)
	writer.Write([]byte(fmt.Sprint(response)))
}

// ResponseWriterXML function
func ResponseWriterXML(writer http.ResponseWriter, response interface{}) {
	writer.Header().Set(ContentType, ContentTypeApplicationXML)
	xmlHeader := "<?xml version=\"1.0\" encoding=\"utf-8\"?>"
	_, _ = writer.Write([]byte(xmlHeader))
	_ = xml.NewEncoder(writer).Encode(response)
}

// ResponseWriterJSON function
func ResponseWriterJSON(writer http.ResponseWriter, response interface{}) {
	writer.Header().Set(ContentType, ContentTypeApplicationJSON)
	//js, _ := json.Marshal(response)
	//writer.Write(js)
	_ = json.NewEncoder(writer).Encode(response)
}

// check if string slice contains a value
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// IsPtrType tester
func IsPtrType(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr
}

// IsPtrValue tester
func IsPtrValue(t reflect.Value) bool {
	return t.Kind() == reflect.Ptr
}

// is time.Time type
func IsTimeType(t reflect.Type) bool {
	return t.String() == "time.Time"
}

// isCustomType tester
func IsCustomType(t reflect.Type) bool {
	if t.PkgPath() != "" {
		if t.PkgPath() == "syscall" {
			return false
		}
		return true
	}

	if k := t.Kind(); k == reflect.Array || k == reflect.Chan || k == reflect.Map || IsPtrType(t) || k == reflect.Slice {
		return IsCustomType(t.Elem()) || k == reflect.Map && IsCustomType(t.Key())
	} else if k == reflect.Struct {
		for i := t.NumField() - 1; i >= 0; i-- {
			if IsCustomType(t.Field(i).Type) {
				return true
			}
		}
	}
	return false
}

func Indirect(reflectValue reflect.Value) reflect.Value {
	for IsPtrValue(reflectValue) {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func IndirectType(reflectType reflect.Type) reflect.Type {
	for IsPtrType(reflectType) || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Struct:
		return false
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

// ptr wraps the given value with pointer: V => *V, *V => **V, etc.
func Ptr(v reflect.Value) reflect.Value {
	pt := reflect.PtrTo(v.Type()) // create a *T type.
	pv := reflect.New(pt.Elem())  // create a reflect.Value of type *T.
	pv.Elem().Set(v)              // sets pv to point to underlying value of v.
	return pv
}

func JsonIt(a interface{}) string {
	dataType := IndirectType(reflect.TypeOf(a))
	switch dataType.Kind() {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Bool,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%v", a)
	}
	data, _ := json.Marshal(a)
	return string(data)
}

func IsNativeKind(v reflect.Kind) bool {
	switch v {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Bool,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return true
	}
	return false
}

func Getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func SafeAtoi(str string, fallback int) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		return fallback
	}
	return value
}

func StringPtr(str string) *string {
	return &str
}
