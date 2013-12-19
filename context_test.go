/*
 * filename   : context_test.go
 * created at : 2013-12-17 19:58:02
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package cloak

import (
    "bytes"
    "testing"
    "reflect"
    "net/http"
    "net/http/httptest"
)

func TestNewRequestContext(t *testing.T) {
    body := "sample text"
    model := reflect.Indirect(reflect.ValueOf(&myType{})).Type()
    buf := bytes.NewBufferString(body)
	r, err := http.NewRequest("GET", "http://bar/foo", buf)
	if err != nil {
		t.Fatalf(err.Error())
	}
    mgr := NewResourceManager([]Middleware{})
	w := httptest.NewRecorder()
    c := NewRequestContext(w, r, model, mgr)
    if string(c.body) != body || c.req != r ||
        c.resp != w || c.model != model ||
        c.session != nil || c.instance !=nil {
        t.Fatalf("Make request context failed")
    }
}

func TestSetGet(t *testing.T) {
    _, _, c := makeNewRequest(t, "GET", "http://bar/foo", "", &myType{})
    c.Set("dummy", "1234")
    ret := c.Get("dummy")
    if  ret != "1234"{
        t.Fatalf("Set/Get failed: incorrect value")
    }
}

func TestSetGetDefaultValue(t *testing.T) {
    _, _, c := makeNewRequest(t, "GET", "http://bar/foo", "", &myType{})
    ret := c.Get("dummy", "1234")
    if  ret != "1234"{
        t.Fatalf("Get failed with default value")
    }
}

func TestGetNotExists(t *testing.T) {
    _, _, c := makeNewRequest(t, "GET", "http://bar/foo", "", &myType{})
    defer func() {
        if r := recover(); r == nil {
            t.Fatalf("Get failed: error should be raised when value not exist.")
        }
    }()
    _ = c.Get("dummy")
}
