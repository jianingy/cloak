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
	w := httptest.NewRecorder()
    c := NewRequestContext(w, r, model)
    if string(c.body) != body || c.req != r ||
        c.resp != w || c.model != model ||
        c.session != nil || c.instance !=nil {
        t.Fatalf("Make request context failed")
    }
}

func TestSetGet(t *testing.T) {
    _, _, c := makeNewRequest(t, "GET", "http://bar/foo", "", &myType{})
    c.Set("dummy", "1234")
    ret, found := c.Get("dummy")
    if found && ret != "1234"{
        t.Fatalf("Set/Get failed: incorrect value")
    } else if found == false {
        t.Fatalf("Set/Get failed: value not found")
    }
}

func TestGetNotExists(t *testing.T) {
    _, _, c := makeNewRequest(t, "GET", "http://bar/foo", "", &myType{})
    _, found := c.Get("dummy")
    if found {
        t.Fatalf("Get failed: value should not exists.")
    }
}

func TestGetBindInstance(t *testing.T) {
    _, _, c := makeNewRequest(t, "GET", "http://bar/foo", "", &myType{})
    c.BindInstance(&myType{Code:1234, Status: "ok"})
    instance := c.Instance().(*myType)
    if instance.Code != 1234 || instance.Status != "ok" {
        t.Fatalf("Get instance failed: got %v, %v expected 1234, ok.",
            instance.Code, instance.Status)
    }
}

func TestGetInstanceUnbind(t *testing.T) {
    _, _, c := makeNewRequest(t, "GET", "http://bar/foo", "", &myType{})
    defer func() {
        if r := recover(); r != nil {
            /* ok, we got error */
        } else {
            /* not ok */
            t.Fatalf("Unbind instance should not raise error.")
        }
    }()
    _ = c.Instance()
}

func TestGetBindSession(t *testing.T) {
    _, _, c := makeNewRequest(t, "GET", "http://bar/foo", "", &myType{})
    c.BindSession(&mySession{})
    _ = c.Session()
}

func TestGetSessionUnbind(t *testing.T) {
    _, _, c := makeNewRequest(t, "GET", "http://bar/foo", "", &myType{})

    defer func() {
        if r := recover(); r != nil {
            /* ok, we got error */
        } else {
            /* not ok */
            t.Fatalf("Unbind session error didn't be generated.")
        }
    }()
    _ = c.Session()
}
