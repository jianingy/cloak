/*
 * filename   : middlewares.go
 * created at : 2013-12-15 12:21:10
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package cloak

import (
    "net/http"
    "encoding/json"
)

type (
    Middleware func (ResourceHandler) ResourceHandler
)

func FaultWrapperMiddleware(resh ResourceHandler) ResourceHandler {
    return func (context *RequestContext) (interface{}, *ResponseError) {
        resp, err := resh(context)
        if err == nil {
            return resp, err
        }
        error := map[string]interface{}{
            "status_code": err.StatusCode,
            "message": err.StatusText,
        }
        return error, nil
    }
}

func JSONDeserializeMiddleware(resh ResourceHandler) ResourceHandler {
    return func (context *RequestContext) (interface{}, *ResponseError) {
        if context.req.Method == "PUT" || context.req.Method == "POST" {
            if len(context.body) > 0 {
                instance := createInstanceByType(context.model)
                err := json.Unmarshal(context.body, instance)
                if err != nil {
                    return nil, NewResponseError(http.StatusBadRequest,
                        err.Error())
                }
                context.Set("decoder.instance", instance)
            } else {
                return nil, NewResponseError(http.StatusBadRequest,
                    "Empty request body")
            }
        }
        resp, err := resh(context)
        return resp, err
    }
}

func JSONSerializeMiddleware(resh ResourceHandler) ResourceHandler {
    return func (context *RequestContext) (interface{}, *ResponseError) {
        resp, err := resh(context)
        if err != nil {
            return resp, err
        }
        marshalled, error := json.Marshal(resp)
        if err != nil {
            return resp, NewResponseError(http.StatusBadRequest, error.Error())
        }
        context.resp.Header().Set("Content-Type", "application/json")
        return string(marshalled), err
    }
}
