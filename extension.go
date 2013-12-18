/*
 * filename   : extension.go
 * created at : 2013-12-15 19:28:09
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package cloak

type Extension interface {
    OnAddResource(interface{})
    OnBind(map[string]interface{})
    OnSessionBegin(*RequestContext)
    OnSessionEnd(*RequestContext)
}
