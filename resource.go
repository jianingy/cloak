/*
 * filename   : resource.go
 * created at : 2013-12-15 11:42:03
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package cloak

import (
    "fmt"
    "net/http"
    "reflect"
)

type (
    ResourceHandler func (*RequestContext) (interface{}, *ResponseError)
    Resource struct {}
    ResourceAction interface {
        List(*RequestContext)   (interface{}, *ResponseError)
        Create(*RequestContext) (interface{}, *ResponseError)
        Show(*RequestContext)   (interface{}, *ResponseError)
        Update(*RequestContext) (interface{}, *ResponseError)
        Delete(*RequestContext) (interface{}, *ResponseError)
    }
)

func NewResourceHandler(model interface{}, app ResourceHandler, mgr *ResourceManager) HttpHandler {

    /* wrap middlewares */
    for _, mdw := range mgr.middlewares {
        app = mdw(app)
    }

    /* find the type of resource model */
    t := reflect.Indirect(reflect.ValueOf(model)).Type()

    return func(w http.ResponseWriter, r *http.Request) {
        /* Start a request Life-Cycle */

        /* Create a new RequestContext per each request */
        context := NewRequestContext(w, r, t, mgr)

        /* call application */
        resp, err := app(context)
        if err != nil {
            panic(err)
        }

        /* generate response */
        switch (reflect.ValueOf(resp).Kind()) {
        case reflect.String:
            fmt.Fprintf(w, "%s", resp)
        default:
            fmt.Fprintf(w, "%v", resp)
        }

        /* Request ends here */
    }
}

func (*Resource) Instance(context *RequestContext) interface{} {
    return context.Get("decoder.instance")
}

func (*Resource) List(context *RequestContext) (interface{}, *ResponseError) {
    return nil, NewResponseError(http.StatusNotImplemented,"Not Implement")
}

func (*Resource) Create(context *RequestContext) (interface{}, *ResponseError) {
    return nil, NewResponseError(http.StatusNotImplemented,"Not Implement")
}

func (*Resource) Show(context *RequestContext) (interface{}, *ResponseError) {
    return nil, NewResponseError(http.StatusNotImplemented,"Not Implement")
}

func (*Resource) Update(context *RequestContext) (interface{}, *ResponseError) {
    return nil, NewResponseError(http.StatusNotImplemented,"Not Implement")
}

func (*Resource) Delete(context *RequestContext) (interface{}, *ResponseError) {
    return nil, NewResponseError(http.StatusNotImplemented,"Not Implement")
}
