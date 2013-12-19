/*
 * filename   : resource_test.go
 * created at : 2013-12-18 10:57:22
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package cloak

import (
    "bytes"
    "testing"
    "net/http"
    "net/http/httptest"
)

func TestNewResourceHandler(t *testing.T) {
    expected := "hello, world"
    mgr := NewResourceManager([]Middleware{})
    h := func (context *RequestContext) (interface{}, *ResponseError) {
        return expected, nil
    }
    f := NewResourceHandler(&myType{}, h, mgr)
	r, err := http.NewRequest("GET", "http://foo/bar", new(bytes.Buffer))
	if err != nil {
		t.Fatalf(err.Error())
	}
	w := httptest.NewRecorder()
    f(w, r)
    if w.Body.String() != expected {
        t.Fatalf("got %v expected %v\n", w.Body.String(), expected)
    }
}

func TestNewResourceHandlerWithError(t *testing.T) {
    expected := "some error"
    mgr := NewResourceManager([]Middleware{})
    h := func (context *RequestContext) (interface{}, *ResponseError) {
        return "", NewResponseError(http.StatusBadRequest, expected)
    }
    f := NewResourceHandler(&myType{}, h, mgr)
	r, err := http.NewRequest("GET", "http://foo/bar", new(bytes.Buffer))
	if err != nil {
		t.Fatalf(err.Error())
	}
	w := httptest.NewRecorder()
    defer func() {
        if r := recover(); r != nil {
            if r.(error).Error() != expected {
                t.Fatalf("got %v expected %v\n", r, expected)
            }
        }
    }()
    f(w, r)
}
