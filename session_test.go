/*
 * filename   : session_test.go
 * created at : 2013-12-17 20:13:28
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package cloak

type mySession struct {}

func (*mySession) Begin() {}
func (*mySession) Commit() {}
func (*mySession) Rollback() {}
