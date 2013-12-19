/*
 * filename   : middlewares_test.go
 * created at : 2013-12-17 17:34:32
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package cloak

import (
    "net/http"
    "testing"
)

func TestFaultWrapperMiddleware(t *testing.T) {

    f := func(ctx *RequestContext) (interface{}, *ResponseError) {
        return nil, NewResponseError(http.StatusBadRequest, "An Error Occured")
    }
	ret, err := FaultWrapperMiddleware(f)(nil)

    if err != nil {
        t.Fatalf(err.Error())
    }

    status := ret.(map[string]interface{})

    if _, found := status["status_code"]; !found {
        t.Fatalf("Faultwrapper format incorrect: missing status_code.")
    }

    if _, found := status["message"]; !found {
        t.Fatalf("Faultwrapper format incorrect: missing message.")
    }
}

func TestJSONSerializeMiddleware(t *testing.T) {
    _, _, c := makeNewRequest(t, "GET", "http://bar/foo", "", &myType{})
    f := func(*RequestContext) (interface{}, *ResponseError) {
        return map[string]string{"code": "200", "status": "ok"}, nil
    }
    ret, err := JSONSerializeMiddleware(f)(c)
    if err != nil {
        t.Fatalf(err.Error())
    }
    expected := `{"code":"200","status":"ok"}`
    if ret != expected {
        t.Fatalf("got %t, expects %s", ret, expected)
    }
}

func TestJSONDeserializeMiddlewareWithNormalJSONPUT(t *testing.T) {
    json := `{"code":200, "status":"ok"}`
    _, _, c := makeNewRequest(t, "PUT", "http://bar/foo", json, &myType{})
    f := func(ctx *RequestContext) (interface{}, *ResponseError) {
        instance := ctx.Get("decoder.instance").(*myType)
        if instance.Code != 200 || instance.Status != "ok" {
            t.Fatalf("JSON Deserialization failed. got %v, %v expected 200, ok",
                instance.Code, instance.Status)
        }
        return nil, nil
    }
    _, err := JSONDeserializeMiddleware(f)(c)
    if err != nil {
        t.Fatalf(err.Error())
    }
}

func TestJSONDeserializeMiddlewareWithNormalJSONPOST(t *testing.T) {
    json := `{"code":200, "status":"ok"}`
    _, _, c := makeNewRequest(t, "POST", "http://bar/foo", json, &myType{})
    f := func(ctx *RequestContext) (interface{}, *ResponseError) {
        instance := ctx.Get("decoder.instance").(*myType)
        if instance.Code != 200 || instance.Status != "ok" {
            t.Fatalf("JSON Deserialization failed. got %v, %v expected 200, ok",
                instance.Code, instance.Status)
        }
        return nil, nil
    }
    _, err := JSONDeserializeMiddleware(f)(c)
    if err != nil {
        t.Fatalf(err.Error())
    }
}

func TestJSONDeserializeMiddlewareWithEmptyJSONPUT(t *testing.T) {
    _, _, c := makeNewRequest(t, "POST", "http://bar/foo", "", &myType{})
    f := func(ctx *RequestContext) (interface{}, *ResponseError) {
        return nil, nil
    }
    _, err := JSONDeserializeMiddleware(f)(c)
    if err == nil {
        t.Fatalf("Deserializer didn't detect empty request body.")
    }
}

func TestJSONDeserializeMiddlewareWithEmptyJSONPOST(t *testing.T) {
    _, _, c := makeNewRequest(t, "POST", "http://bar/foo", "", &myType{})
    f := func(ctx *RequestContext) (interface{}, *ResponseError) {
        return nil, nil
    }
    _, err := JSONDeserializeMiddleware(f)(c)
    if err == nil {
        t.Fatalf("Deserializer didn't detect empty request body.")
    }
}

func TestJSONDeserializeMiddlewareWithBadJSONPUT(t *testing.T) {
    _, _, c := makeNewRequest(t, "POST", "http://bar/foo", "", &myType{})
    f := func(ctx *RequestContext) (interface{}, *ResponseError) {
        return nil, nil
    }
    _, err := JSONDeserializeMiddleware(f)(c)
    if err == nil {
        t.Fatalf("Deserializer didn't detect malformed json.")
    }
}

func TestJSONDeserializeMiddlewareWithBadJSONPOST(t *testing.T) {
    _, _, c := makeNewRequest(t, "POST", "http://bar/foo", "", &myType{})
    f := func(ctx *RequestContext) (interface{}, *ResponseError) {
        return nil, nil
    }
    _, err := JSONDeserializeMiddleware(f)(c)
    if err == nil {
        t.Fatalf("Deserializer didn't detect malformed json.")
    }
}
