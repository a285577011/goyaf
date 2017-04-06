package form

//校验器
type Validate struct {
	Type     string
	Errmsg   string
	Min      map[string]int
	Max      map[string]int
	FloatMin map[string]float64 //针对小数类型的校验
	FloatMax map[string]float64 //针对小数类型的校验
	Func     func(string) bool
}
