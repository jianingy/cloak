/*
 * filename   : hood.go
 * created at : 2013-12-14 20:30:52
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package cloak

import (
    "log"
    "reflect"
    "net/http"
    "github.com/eaigner/hood"
)

type (
    HoodExtension struct {
        hd *hood.Hood
        context *RequestContext
    }

    HoodSessionT interface {
        Add(interface {})
        Delete(interface{})
    }

    HoodSession struct {
        hd *hood.Hood
        tx *hood.Hood
    }

    HoodResource struct {}
)

func (ext *HoodExtension) OnAddResource(model interface{}) {
}

func (ext *HoodExtension) OnBind(opts map[string]interface{}) {
    if hd, found := opts["exts.hood"]; found {
        ext.hd = hd.(*hood.Hood)
        log.Printf("Bind Hood instance: %v\n", ext.hd)
    }
}

func (ext *HoodExtension) OnSessionBegin(context *RequestContext) {
    context.BindSession(NewHoodSession(ext.hd))
}

func (ext *HoodExtension) OnSessionEnd(context *RequestContext) {
}

func NewHoodSession(hd *hood.Hood) *HoodSession {
    return &HoodSession{hd:hd, tx: nil}
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



func (*HoodResource) List(context *RequestContext) (interface{}, *ResponseError) {
    return nil, NewResponseError(http.StatusNotImplemented,"Not Implement")
}

func (*HoodResource) Create(context *RequestContext) (interface{}, *ResponseError) {
    session := context.Session().(HoodSessionT)
    session.Add(context.Instance())
    return context.Instance(), nil
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
