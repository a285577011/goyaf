package lib

//排除空字符串
func ArrayFilter(array []string) []string {
	r := make([]string, 0)
	for _, v := range array {
		if len(v) == 0 {
			continue
		}
		r = append(r, v)
	}
	return r
}

//对比返回在 array1 中但是不在 array2 及任何其它参数数组中的值
func ArrayDiff(array1 []string, array2 []string) []string {
	r := make([]string, 0)
	for _, v := range array1 {
		if InArray(array2, v) {
			continue
		}
		r = append(r, v)
	}
	return r
}

//在 array 中搜索 needle
func InArray(array []string, needle string) bool {
	for _, v := range array {
		if v == needle {
			return true
		}
	}
	return false
}

//排序方法
func Sort(array []int) []int {
	for i := 0; i < len(array); i++ {
		for j := 0; j < len(array)-i-1; j++ {
			if array[j] < array[j+1] {
				array[j], array[j+1] = array[j+1], array[j]
			}
		}
	}
	return array
}
