package form

import ()

//表单接口
type Form interface {
	Validate() bool
	GetFieldsFirstValue() map[string]string
	GetFieldsErrmsg() map[string]string
}
