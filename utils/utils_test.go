package utils

import (
	"bytes"
	"encoding/xml"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestGetRequestBody(t *testing.T) {
	buf := []byte(`{"msg":"Hello there!"}`)
	request, _ := http.NewRequest(http.MethodPost,
		"https://example.com/api/manage/secretreload",
		bytes.NewBuffer(buf))
	data, err := GetRequestBody(request)
	if err != nil {
		t.Errorf("TestGetRequestBody: failed with error %v", err)
	}

	if len(buf) != len(data) {
		t.Errorf("TestGetRequestBody: failed expected length %v, but got %v", len(buf), len(data))
	}

}

func TestGetRequestBodyError(t *testing.T) {
	request, _ := http.NewRequest(http.MethodPost,
		"https://example.com/api/manage/secretreload",
		nil)
	_, err := GetRequestBody(request)
	if err == nil {
		t.Errorf("TestGetRequestBodyError: expected failure, but got none")
		return
	}
	if err.Error() != "no body found in request" {
		t.Errorf("TestGetRequestBodyError: expected failure with message: \"%v\", but got \"%v\"",
			"no body found in request", err.Error())
	}

}

func TestFileExists(t *testing.T) {
	result := FileExists(os.Args[0])
	if !result {
		t.Errorf("TestFileExists success: expected %v, but got %v", true, result)
	}

	result = FileExists(os.Args[0] + "1")
	if result {
		t.Errorf("TestFileExists not found: expected %v, but got %v", false, result)
	}
}

func TestErrorToText(t *testing.T) {
	resp := &GenericResponse{
		Result:            "result",
		ResultExplanation: "result_explanation",
	}
	expectedResult := "result;result_explanation;;;"
	result := ErrorToText(*resp)
	if result != expectedResult {
		t.Errorf("Invalid result: expected[%s], got[%s]",
			expectedResult, result)
	}
}

func TestResponseWriterTextHtml(t *testing.T) {
	response := "{\"Message\":\"Hello World!\"}"
	responseWriter := httptest.NewRecorder()
	ResponseWriterTextHtml(responseWriter, response)
	if responseWriter.Header().Get("Content-Type") != ContentTypeApplicationXML {
		t.Errorf("Invalid Content-Type expected[%s], got[%s]",
			ContentTypeApplicationXML, responseWriter.Header().Get("Content-Type"))
	}
	if response != responseWriter.Body.String() {
		t.Errorf("Invalid body data: expecting[%s], got[%s]", response,
			responseWriter.Body.String())
	}
}

func TestResponseWriterXML(t *testing.T) {
	response := &struct {
		XMLName *xml.Name `xml:"root" json:",omitempty"`
		Message string    `json:"message"`
	}{
		XMLName: &xml.Name{},
		Message: "Hello World!",
	}
	expected := "<?xml version=\"1.0\" encoding=\"utf-8\"?><root><Message>Hello World!</Message></root>"
	responseWriter := httptest.NewRecorder()
	ResponseWriterXML(responseWriter, response)
	if responseWriter.Header().Get("Content-Type") != ContentTypeApplicationXML {
		t.Errorf("Invalid Content-Type expected[%s], got[%s]",
			ContentTypeApplicationXML, responseWriter.Header().Get("Content-Type"))
	}
	if expected != responseWriter.Body.String() {
		t.Errorf("Invalid response returning, expected[%s], got[%s]", expected,
			responseWriter.Body.String())
	}
}

func TestResponseWriterJSON(t *testing.T) {
	response := &struct {
		Message string `json:"message"`
	}{
		Message: "Hello World!",
	}
	expected := "{\"message\":\"Hello World!\"}\n"

	responseWriter := httptest.NewRecorder()
	ResponseWriterJSON(responseWriter, response)
	if responseWriter.Header().Get("Content-Type") != ContentTypeApplicationJSON {
		t.Errorf("Invalid Content-Type expected[%s], got[%s]",
			ContentTypeApplicationJSON, responseWriter.Header().Get("Content-Type"))
	}
	if expected != responseWriter.Body.String() {
		t.Errorf("Invalid body data: expecting[%s], got[%s]", expected,
			responseWriter.Body.String())
	}
}

func TestContains(t *testing.T) {
	search := []string{"findme", "mehere", "herefind"}
	findWhat := "mehere"
	result := Contains(search, findWhat)
	if !result {
		t.Errorf("Failed to find[%s]: expected[%v], got[%v]", findWhat, true, result)
	}
	findWhat = "here"
	result = Contains(search, findWhat)
	if result {
		t.Errorf("Failed to find[%s]: expected[%v], got[%v]", findWhat, true, result)
	}
}

func TestIsTimeType(t *testing.T) {
	valueType := reflect.TypeOf("str")
	if IsTimeType(valueType) {
		t.Errorf("Expected type inspection to be false")
	}
	valueType = reflect.TypeOf(time.Now())
	if !IsTimeType(valueType) {
		t.Errorf("Expected type inspection to be true")
	}
}

func TestIsCustomType(t *testing.T) {
	//valueType := reflect.TypeOf(syscall.AddrinfoW{})
	//result := IsCustomType(valueType)
	//if result {
	// t.Error("Invalid response, expecting[false], got[true]")
	//}

	valueType := reflect.TypeOf(reflect.Method{})
	result := IsCustomType(valueType)
	if !result {
		t.Error("Invalid response, expecting[true], got[false]")
	}

	type A struct {
		Number string
	}
	valueType = reflect.TypeOf(A{})
	result = IsCustomType(valueType)
	if !result {
		t.Error("Invalid response, expecting[true], got[false]")
	}

	valueType = reflect.TypeOf(struct{ i int }{})
	result = IsCustomType(valueType)
	if result {
		t.Error("Invalid response, expecting[true], got[false]")
	}

	valueType = reflect.TypeOf(struct{ Abc A }{})
	result = IsCustomType(valueType)
	if !result {
		t.Error("Invalid response, expecting[true], got[false]")
	}

	valueType = reflect.TypeOf(errors.New(""))
	result = IsCustomType(valueType)
	if !result {
		t.Error("Invalid response, expecting[true], got[false]")
	}

	valueType = reflect.TypeOf([]byte{})
	result = IsCustomType(valueType)
	if result {
		t.Error("Invalid response, expecting[false], got[true]")
	}

}

func TestIndirect(t *testing.T) {
	type A struct {
		Number string
	}
	result := Indirect(reflect.ValueOf(&A{}))
	if result.Kind() != reflect.Struct {
		t.Errorf("Invalid Kind got returned, expected[%s], got[%s]", reflect.Struct.String(),
			result.Kind().String())
	}

	result = Indirect(Ptr(reflect.ValueOf(&A{})))
	if result.Kind() != reflect.Struct {
		t.Errorf("Invalid Kind got returned, expected[%s], got[%s]", reflect.Struct.String(),
			result.Kind().String())
	}
}

func TestIsNil(t *testing.T) {
	type A struct {
		Number string
	}

	if !IsNil(nil) {
		t.Error("IsNil test failed, expect[true], got[false]")
	}

	if IsNil(10) {
		t.Error("IsNil test failed, expect[false], got[true]")
	}

	value := &A{}
	if IsNil(value) {
		t.Error("IsNil test failed, expect[false], got[true]")
	}

	if IsNil(*value) {
		t.Error("IsNil test failed, expect[false], got[true]")
	}
}

func TestJsonIt(t *testing.T) {
	result := JsonIt(10)
	if result != "10" {
		t.Errorf("JsonIt test failed, expect[10], got[%s]", result)
	}
	result = JsonIt(struct{ Message string }{Message: "Hello!"})
	if result != "{\"Message\":\"Hello!\"}" {
		t.Errorf("JsonIt test failed, expect[{\"Message\":\"Hello!\"}], got[%s]", result)
	}
}

func TestIsNativeKind(t *testing.T) {
	if !IsNativeKind(reflect.String) {
		t.Errorf("IsNativeKind test failed, expected[true], got[false]")
	}
	if IsNativeKind(reflect.Struct) {
		t.Errorf("IsNativeKind test failed, expected[false], got[true]")
	}
}

func TestGetenv(t *testing.T) {
	_ = os.Setenv("test01", "test01")
	result := Getenv("test00", "test00")
	if result != "test00" {
		t.Errorf("Getenv test failed, expected[test00], got[%s]", result)
	}
	result = Getenv("test01", "test00")
	if result != "test01" {
		t.Errorf("Getenv test failed, expected[test01], got[%s]", result)
	}
}

func TestSafeAtoi(t *testing.T) {
	result := SafeAtoi("123weq", 321)
	if result != 321 {
		t.Errorf("SafeAtoi test failed, expected[321], got[%d]", result)
	}
	result = SafeAtoi("123", 321)
	if result != 123 {
		t.Errorf("SafeAtoi test failed, expected[123], got[%d]", result)
	}
}

func TestStringPtr(t *testing.T) {
	value := StringPtr("hello")
	if value == nil {
		t.Error("StringPtr test failed nil pointer returned")
		return
	}
	if *value != "hello" {
		t.Errorf("StringPtr test failed, expected[hello], got[%s]", *value)
	}
}
