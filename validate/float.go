package validate

import (
	"git.oschina.net/pbaapp/goyaf/lib"
)

type Float struct {
	options  map[string]interface{}
	value    float64
	minIsSet bool
	min      float64
	maxIsSet bool
	max      float64
}

//校验方法
func (this *Float) Validate() bool {
	if !this.analyseOptions() {
		return false
	}

	if this.minIsSet && this.value < this.min {
		return false
	}
	if this.maxIsSet && this.value > this.max {
		return false
	}

	return true
}

//分析options参数
func (this *Float) analyseOptions() bool {
	if len(this.options) == 0 {
		return false
	}
	var err error
	this.value, err = lib.InterfaceToFloat64(this.options["value"])
	if err != nil {
		return false
	}

	var ok bool
	_, ok = this.options["min"]
	if ok {
		this.min, err = lib.InterfaceToFloat64(this.options["min"])
		if err != nil {
			return false
		}
	}

	_, ok = this.options["max"]
	if ok {
		this.max, err = lib.InterfaceToFloat64(this.options["max"])
		if err != nil {
			return false
		}
	}

	return true
}

//实例化float校验器
func NewFloat(options map[string]interface{}) Float {
	f := Float{}
	f.options = options
	return f
}
