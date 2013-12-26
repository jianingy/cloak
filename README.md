[![Build Status](https://secure.travis-ci.org/jianingy/restle.png?branch=master)](http://travis-ci.org/jianingy/restle)

This document is still W.I.P

## What is Restle?

Restle is a lightweight RESTful framework writing in golang. It can be
extended by a middleware mechanism similar to Python's WSGI.

## Basic Usage

### 10-minutes tutorial by example

```go
package main

import (
    "net/http"
    "github.com/jinzhu/gorm"
    "github.com/jianingy/restle"
    _ "github.com/lib/pq"
)

type Fruit struct {
    restle.GormResource

    Name  string `validate:"presence"`
    Color string
}



func main() {
    gorm, err := gorm.Open("postgres", "user=jianingy dbname=jianingy password=123465 host=172.24.232.1")
    if err != nil {
        panic(err)
    }
    resource := restle.NewResourceManager(nil)
    resource.SetOption(map[string]interface{}{"gorm.db":gorm})
    resource.AddResource(&Fruit{})
    http.Handle("/", resource.Router())
    http.ListenAndServe(":8080", nil)
}
```

The path to rest API is derived by resource name. In the above example, the resource name is 'Fruit', therefore
its API paths is '/fruits'.

The default behavior of HTTP actions are,

| GET /fruits | List all items |
| POST /fruits | Create a new item |
| GET /fruits/[id] | Show details of an item |
| POST /furits/[id] | Update a specified item |
| DELETE /furits/[id] | Delete a specified item |



## Extend with middlewares

Restle uses a chainable middleware mechanism to process each request. Middlewares can be used for
serializing / deserializing data, authentication, session management and etc.

The order of middlewares is determined by passing a slice which contains middleware instance to function
'NewResourceManager'. For example,

```go
resource := restle.NewResourceManager([]result.Middleware{MyMiddleware})
``` 

Middlewares are essentially normal function with the following prototype,

```go
type Middleware func (ResourceHandler) ResourceHandler
```

A typical middleware may looks like,

```go
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
```

A middleware must return a closure. In the closure, normally, a middleware have to call
"resh", a Resource Handler, to get the result of other middlewares or the underly resource
handler. Then, it can apply certain modification on the result and return the revamped result.

If an error occured during the processing, the return value 'ResponseError' has to be none-nil.

