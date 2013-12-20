/*
 * filename   : error.go
 * created at : 2013-12-18 11:29:48
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package restle

type ResponseError struct {
    StatusCode int
    StatusText string
}

func NewResponseError(code int, text string) *ResponseError {
    return &ResponseError{StatusCode:code, StatusText:text}
}

func (e *ResponseError) Error() string {
    return  e.StatusText
}
