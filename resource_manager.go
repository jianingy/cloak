/*
 * filename   : resource_manager.go
 * created at : 2013-12-14 12:43:16
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package restle

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
    options map[string]interface{}
}

func NewResourceManager(mdws []Middleware) *ResourceManager {
    mgr := &ResourceManager{router:mux.NewRouter()}
    if mdws == nil {
        mgr.middlewares = []Middleware{
            JSONDeserializeMiddleware,
            //HoodSessionMiddleware,
            GormSessionMiddleware,
            FaultWrapperMiddleware,
            JSONSerializeMiddleware,
        }
    }
    mgr.options = make(map[string]interface{})
    return mgr
}

func (mgr *ResourceManager) AddResource(model interface{}) {
    if reflect.ValueOf(model).Kind() != reflect.Ptr ||
        reflect.Indirect(reflect.ValueOf(model)).Kind() != reflect.Struct {
        panic("argument model has to be a pointer pointing to a struct")
    }

    r := mgr.router
    resourceT := model.(ResourceAction)
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

func (mgr *ResourceManager) SetOption(opts map[string]interface{}) {
    for key, val := range opts {
        mgr.options[key] = val
    }
}

func (mgr *ResourceManager) Option(name string, args ...interface{}) interface{} {
    val, found := mgr.options[name]
    if found {
        return val
    } else {
        if len(args) > 0 {
            return args[0]
        }
        panic(fmt.Sprintf("value %s not found", name))
    }
}

func (mgr *ResourceManager) Router() *mux.Router {
    return mgr.router
}
