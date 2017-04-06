package form

//字段定义
type Field struct {
	Name     string
	Values   []string
	Required bool
	Errmsg   string
	Validate []Validate
}
