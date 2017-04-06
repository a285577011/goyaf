package lib

//对比返回在 map1 中但是不在 map2 及任何其它字典中的值
func MapDiff(map1 map[string]string, map2 map[string]string) map[string]string {
	r := make(map[string]string)
	for k, v := range map1 {
		map2v, ok := map2[k]
		if ok && map2v == v {
			continue
		}
		r[k] = v
	}
	return r
}

//用 map2 中的值替换 map1 中的值
func MapReplace(map1 map[string]string, map2 map[string]string) map[string]string {
	for k, v := range map2 {
		if map1[k] != v {
			map1[k] = v
		}
	}
	return map1
}

//取出map中的某个字段的数据
func MapColumn(m []map[string]string, column string) []string {
	a := make([]string, len(m))
	i := 0
	for _, v := range m {
		a[i] = v[column]
		i++
	}
	return a
}
