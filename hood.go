/*
 * filename   : hood.go
 * created at : 2013-12-14 20:30:52
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package restle

import (
    _ "log"
    "reflect"
    "net/http"
    "github.com/eaigner/hood"
)

type (
    HoodSession struct {
        hd *hood.Hood
        tx *hood.Hood
    }
    HoodResource struct {
        Resource
    }
)

func HoodSessionMiddleware(resh ResourceHandler) ResourceHandler {
    return func (context *RequestContext) (interface{}, *ResponseError) {
        hd := context.mgr.Option("hood.hd")
        session := NewHoodSession(hd)
        context.Set("hood.session", session)
        session.Begin()
        resp, err := resh(context)
        session.Commit()
        return resp, err
    }
}

func NewHoodSession(hd interface{}) *HoodSession {
    return &HoodSession{hd: hd.(*hood.Hood)}
}

func (session *HoodSession) Begin() {
    session.tx = session.hd.Begin()
}

func (session *HoodSession) Commit() {
    session.tx.Commit()
}

func (session *HoodSession) Rollback() {
    session.tx.Rollback()
}

func (session *HoodSession) Transaction() *hood.Hood {
    return session.tx
}

func (session *HoodSession) Add(instance interface{}) {
    if reflect.ValueOf(instance).Kind() == reflect.Slice {
        session.tx.SaveAll(instance)
    } else {
        session.tx.Save(instance)
    }
}

func (session *HoodSession) Delete(instance interface{}) {
    if reflect.ValueOf(instance).Kind() == reflect.Slice {
        session.tx.DeleteAll(instance)
    } else {
        session.tx.Delete(instance)
    }
}

func (*HoodResource) Session(context *RequestContext) (*HoodSession) {
    return context.Get("hood.session").(*HoodSession)
}

func (h *HoodResource) List(context *RequestContext) (interface{}, *ResponseError) {
    session := h.Session(context)
    result := createInstanceByType(context.ModelType())
    session.Transaction().Find(result)
    return result, nil
}

func (h *HoodResource) Create(context *RequestContext) (interface{}, *ResponseError) {
    instance := h.Instance(context)
    h.Session(context).Add(instance)
    return instance, nil
}

func (self *HoodResource) Show(context *RequestContext) (interface{}, *ResponseError) {
     return context.Vars("id"), nil
}

func (*HoodResource) Update(context *RequestContext) (interface{}, *ResponseError) {
    return nil, NewResponseError(http.StatusNotImplemented,"Not Implement")
}

func (*HoodResource) Delete(context *RequestContext) (interface{}, *ResponseError) {
    return nil, NewResponseError(http.StatusNotImplemented,"Not Implement")
}
