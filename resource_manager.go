/*
 * filename   : resource_manager.go
 * created at : 2013-12-14 12:43:16
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package cloak

import (
    "fmt"
    "log"
    "reflect"
    "strings"
    "github.com/gorilla/mux"
)

type ResourceManager struct {
    router *mux.Router
    middlewares []Middleware
    extensions []Extension
}

func NewResourceManager(mdws []Middleware, exts []Extension) *ResourceManager {
    mgr := &ResourceManager{router:mux.NewRouter()}
    if mdws == nil {
        mgr.middlewares = []Middleware{
            JSONDeserializeMiddleware,
            SessionMiddleware,
            FaultWrapperMiddleware,
            JSONSerializeMiddleware,
        }
    }
    if exts == nil {
        mgr.extensions = []Extension{
            &HoodExtension{},
        }
    }
    return mgr
}

func (mgr *ResourceManager) AddResource(model interface{}) {
    if reflect.ValueOf(model).Kind() != reflect.Ptr ||
        reflect.Indirect(reflect.ValueOf(model)).Kind() != reflect.Struct {
        panic("argument model has to be a pointer pointing to a struct")
    }

    for _, extension := range mgr.extensions {
        extension.OnAddResource(model)
    }

    r := mgr.router
    resourceT := model.(ResourceInterface)
    plural := strings.ToLower(getTypeName(model)) + "s"

    handle := func (p string, m string, h ResourceHandler) {
        r.Path(p).Methods(m).HandlerFunc(NewResourceHandler(model, h, mgr))
    }

    collection := fmt.Sprintf("/%s/", plural)
    handle(collection, "GET", resourceT.List)
    handle(collection, "POST", resourceT.Create)

    object := fmt.Sprintf("/%s/{id:[0-9a-zA-Z-]+}/", plural)
    handle(object, "GET", resourceT.Show)
    handle(object, "PUT", resourceT.Update)
    handle(object, "DELETE", resourceT.Delete)

    log.Printf("Resource `%s' been added to router\n", plural)
}

func (mgr *ResourceManager) Bind(opts map[string]interface{}) {
    for _, extension := range mgr.extensions {
        extension.OnBind(opts)
    }
}

func (mgr *ResourceManager) Router() *mux.Router {
    return mgr.router
}
