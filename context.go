/*
 * filename   : context.go
 * created at : 2013-12-15 12:04:12
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package cloak

import (
    "fmt"
    "reflect"
    "net/http"
    "io/ioutil"
    "github.com/gorilla/mux"
)

type (
    HttpHandler func (http.ResponseWriter, *http.Request)
    RequestHandler func (http.ResponseWriter, *http.Request) (map[string]string, error)
    RequestContext struct {
        mgr        *ResourceManager
        req        *http.Request
        resp       http.ResponseWriter
        body       []byte
        model      reflect.Type
        session    Session
        instance   interface{}
        extra      map[string]interface{}
    }
)

func NewRequestContext(w http.ResponseWriter, r *http.Request, model reflect.Type, mgr *ResourceManager) *RequestContext {
    /* read request body */
    body, err := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
    if err != nil {
        panic(err)
    }

    return &RequestContext {
        mgr: mgr,
        req: r,
        resp: w,
        body: body,
        model: model,
        extra: make(map[string]interface{}),
    }
}

func (context *RequestContext) Vars(name string) string {
    return mux.Vars(context.req)[name]
}

func (context *RequestContext) Set(name string, value interface{}) {
    context.extra[name] = value
}

func (context *RequestContext) Get(name string, args ...interface{}) interface{} {
    value, found := context.extra[name]
    if found {
        return value
    } else {
        if len(args) > 0 {
            return args[0]
        }
        panic(fmt.Sprintf("value %s not found", name))
    }
}

func (context *RequestContext) ModelType() reflect.Type {
    return context.model
}
