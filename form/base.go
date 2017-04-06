package form

import (
	"git.oschina.net/pbaapp/goyaf"
	"net/url"
	"strconv"
)

//表单基类
type Base struct {
	Fields map[string]*Field
}

//校验表单
func (this *Base) Validate() bool {
	result := true
	for k, f := range this.Fields {
		//如果该值必须，则检测是否有值
		if f.Required && len(this.GetFieldFirstValue(k)) == 0 {
			result = false
			this.Fields[k].Errmsg = k + "值不能为空"
			continue
		}
		//如果该值非必须，则如果没有值直接返回true，不进入校验
		if !f.Required && len(this.GetFieldFirstValue(k)) == 0 {
			continue
		}
		//如果没有校验器，则直接返回
		if f.Validate == nil || len(f.Validate) == 0 {
			continue
		}

		//没有传递值
		if len(f.Values) == 0 {
			result = false
			this.Fields[k].Errmsg = f.Validate[0].Errmsg
			continue
		}

		for _, value := range f.Values {
			for _, validate := range f.Validate {
				var validateResult bool

				switch validate.Type {
				case "int": //校验整形
					v, err := strconv.Atoi(value)
					if err != nil {
						validateResult = false
					} else {
						options := map[string]int{"value": v}
						if validate.Min["isset"] == 1 {
							options["min"] = validate.Min["value"]
						}
						if validate.Max["isset"] == 1 {
							options["max"] = validate.Max["value"]
						}
						validateResult = this.ValidateInt(options)
					}
				case "string":
					options := map[string]string{"value": value}
					if validate.Min["isset"] == 1 {
						options["min"] = strconv.Itoa(validate.Min["value"])
					}
					if validate.Max["isset"] == 1 {
						options["max"] = strconv.Itoa(validate.Max["value"])
					}
					validateResult = this.ValidateString(options)
				case "float":
					v, err := strconv.ParseFloat(value, 64)
					if err != nil {
						validateResult = false
					} else {
						options := map[string]float64{"value": v}
						if validate.FloatMin["isset"] == 1.0 {
							options["min"] = validate.FloatMin["value"]
						}
						if validate.FloatMax["isset"] == 1.0 {
							options["max"] = validate.FloatMax["value"]
						}
						validateResult = this.ValidateFloat(options)
					}
				case "func":
					validateResult = validate.Func(value)
				}
				if !validateResult {
					result = false
					this.Fields[k].Errmsg = validate.Errmsg
				}
			}
		}
	}
	return result
}

//获取所有字段第一个值
func (this *Base) GetFieldFirstValue(name string) string {
	field, ok := this.Fields[name]
	if ok && len(field.Values) > 0 {
		return field.Values[0]
	}
	return ""
}

//获取所有字段的第一个值
func (this *Base) GetFieldsFirstValue() map[string]string {
	fieldsFirstValue := make(map[string]string)
	for k, f := range this.Fields {
		if len(f.Values) < 1 {
			fieldsFirstValue[k] = ""
			continue
		}
		fieldsFirstValue[k] = f.Values[0]
	}

	return fieldsFirstValue
}

//获取所有字段的错误消息
func (this *Base) GetFieldsErrmsg() map[string]string {
	fieldsErrmsg := make(map[string]string)
	for k, f := range this.Fields {
		if f.Errmsg == "" {
			continue
		}
		fieldsErrmsg[k] = f.Errmsg
	}

	return fieldsErrmsg
}

//根据传入的值进行设置表单字段值
func (this *Base) SetFieldsValues(data url.Values) {
	for k, values := range data {
		_, ok := this.Fields[k]
		if !ok {
			continue
		}

		this.Fields[k].Values = values
	}
}

//校验整形
func (this *Base) ValidateInt(options map[string]int) bool {
	var ok bool
	var value int
	var min int
	var max int

	value, ok = options["value"]
	if !ok {
		return false
	}

	min, ok = options["min"]
	if ok && value < min {
		return false
	}

	max, ok = options["max"]
	if ok && value > max {
		return false
	}

	return true
}

//校验整形
func (this *Base) ValidateFloat(options map[string]float64) bool {
	var ok bool
	var value float64
	var min float64
	var max float64

	value, ok = options["value"]
	if !ok {
		return false
	}

	min, ok = options["min"]
	if ok && value < min {
		return false
	}

	max, ok = options["max"]
	if ok && value > max {
		return false
	}

	return true
}

//校验字符串
func (this *Base) ValidateString(options map[string]string) bool {
	var ok bool
	var value string
	var min string
	var max string

	value, ok = options["value"]
	if !ok {
		return false
	}

	length := len([]rune(value))

	min, ok = options["min"]
	minInt, _ := strconv.Atoi(min)
	if ok && length < minInt {
		return false
	}

	max, ok = options["max"]
	maxInt, _ := strconv.Atoi(max)
	if ok && length > maxInt {
		return false
	}

	return true
}

func init() {
	goyaf.Debug("init goyaf form base")
}
