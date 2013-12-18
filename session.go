/*
 * filename   : session.go
 * created at : 2013-12-16 20:07:52
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package cloak

type Session interface {
    Begin()
    Commit()
    Rollback()
}
