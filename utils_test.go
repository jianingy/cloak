/*
 * filename   : utils_test.go
 * created at : 2013-12-17 10:39:46
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package restle

import (
    _ "log"
    "bytes"
    "testing"
    "reflect"
    "net/http"
    "net/http/httptest"
)

type myType struct {
    Code   int
    Status string
}

func makeNewRequest(t *testing.T, method, url, body string, model interface{}) (http.ResponseWriter, *http.Request, *RequestContext) {
    buf := bytes.NewBufferString(body)
	r, err := http.NewRequest(method, url, buf)
	if err != nil {
		t.Fatalf(err.Error())
	}
    mgr := NewResourceManager([]Middleware{})
	w := httptest.NewRecorder()
    c := NewRequestContext(w, r, reflect.Indirect(reflect.ValueOf(model)).Type(), mgr)
    return w, r, c
}

func TestGetTypeName(t *testing.T) {
    /* test with pointer */
    if ret, expected := getTypeName(&myType{}), "myType"; ret != expected {
        t.Fatalf("got %t, expects %s", ret, expected)
    }

    /* test with instance */
    if ret, expected := getTypeName(myType{}), "myType"; ret != expected {
        t.Fatalf("got %t, expects %s", ret, expected)
    }
}

func TestCreateInstanceByType(t *testing.T) {

    expected := reflect.Indirect(reflect.ValueOf(&myType{})).Type()
    ret := createInstanceByType(expected)
    converted := reflect.Indirect(reflect.ValueOf(ret))
    if converted.Type() != expected {
        t.Fatalf("got %t, expects %s", reflect.TypeOf(converted), expected)
    }

}
