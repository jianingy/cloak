/*
 * filename   : gorm.go
 * created at : 2013-12-19 21:41:50
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package restle

import (
    _ "log"
    _ "reflect"
    _ "net/http"
    "github.com/jinzhu/gorm"
)

type (
    GormSession struct {
        db *gorm.DB
        tx *gorm.DB
    }
    GormResource struct {
        Resource

        Id       int64
    }
)

func GormSessionMiddleware(resh ResourceHandler) ResourceHandler {
    return func (context *RequestContext) (interface{}, *ResponseError) {
        db := context.mgr.Option("gorm.db")
        session := NewGormSession(db)
        context.Set("gorm.session", session)
        session.Begin()
        resp, err := resh(context)
        session.Commit()
        return resp, err
    }
}

func NewGormSession(db interface{}) *GormSession {
    return &GormSession{db: db.(*gorm.DB)}
}

func (session *GormSession) Begin() {
    session.tx = session.db.Begin()
}

func (session *GormSession) Commit() {
    session.tx.Commit()
}

func (session *GormSession) Rollback() {
    session.tx.Rollback()
}

func (session *GormSession) Transaction() *gorm.DB {
    return session.tx
}

func (session *GormSession) Add(instance interface{}) {
    session.tx.Save(instance)
}

func (session *GormSession) Delete(instance interface{}) {
    session.tx.Delete(instance)
}

func (*GormResource) Session(context *RequestContext) (*GormSession) {
    return context.Get("gorm.session").(*GormSession)
}

func (g *GormResource) List(context *RequestContext) (interface{}, *ResponseError) {
    session := g.Session(context)
    result := createSliceByType(context.ModelType())
    session.Transaction().Find(result)
    return result, nil
}

func (g *GormResource) Create(context *RequestContext) (interface{}, *ResponseError) {
    instance := g.Instance(context)
    g.Session(context).Add(instance)
    return instance, nil
}

func (g *GormResource) Show(context *RequestContext) (interface{}, *ResponseError) {
    id := context.Vars("id")
    session := g.Session(context)
    result := createInstanceByType(context.ModelType())
    session.Transaction().First(result, id)
    return result, nil
}

func (g *GormResource) Update(context *RequestContext) (interface{}, *ResponseError) {
    id := context.Vars("id")
    instance := g.Instance(context)
    session := g.Session(context)
    result := createInstanceByType(context.ModelType())
    session.Transaction().First(result, id)
    session.Transaction().Model(result).Update(instance)
    return instance, nil
}

func (g *GormResource) Delete(context *RequestContext) (interface{}, *ResponseError) {
    id := context.Vars("id")
    session := g.Session(context)
    result := createInstanceByType(context.ModelType())
    session.Transaction().First(result, id)
    session.Transaction().Delete(result)
    return result, nil
}
