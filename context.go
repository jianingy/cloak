/*
 * filename   : context.go
 * created at : 2013-12-15 12:04:12
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package cloak

import (
    "reflect"
    "net/http"
    "io/ioutil"
    "github.com/gorilla/mux"
)

type (
    HttpHandler func (http.ResponseWriter, *http.Request)
    RequestHandler func (http.ResponseWriter, *http.Request) (map[string]string, error)
    RequestContext struct {
        req        *http.Request
        resp       http.ResponseWriter
        body       []byte
        model      reflect.Type
        session    Session
        instance   interface{}
        extra      map[string]interface{}
    }
)

func NewRequestContext(w http.ResponseWriter, r *http.Request, model reflect.Type) *RequestContext {
    /* read request body */
    body, err := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
    if err != nil {
        panic(err)
    }

    return &RequestContext{
        req: r,
        resp: w,
        body: body,
        model: model,
        session: nil,
        instance: nil,
        extra: make(map[string]interface{}),
    }
}

func (context *RequestContext) Vars(name string) string {
    return mux.Vars(context.req)[name]
}

func (context *RequestContext) Set(name string, value interface{}) {
    context.extra[name] = value
}

func (context *RequestContext) Get(name string) (interface{}, bool) {
    value, found := context.extra[name]
    return value, found
}

func (context *RequestContext) BindSession(session Session) {
    if context.session != nil {
        panic("session has already been bind")
    } else {
        context.session = session
    }
}

func (context *RequestContext) Session() Session {
    if context.session != nil {
        return context.session
    } else {
        panic("session has not been bind")
    }
}

func (context *RequestContext) BindInstance(instance interface{}) {
    if context.instance != nil {
        panic("instance has already been bind")
    } else {
        context.instance = instance
    }
}

func (context *RequestContext) Instance() interface{} {
    if context.instance != nil {
        return context.instance
    } else {
        panic("instance has not been bind")
    }
}
