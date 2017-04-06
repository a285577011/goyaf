package goyaf

import ()

var Config map[string]map[string]string

func SetConfig(config map[string]map[string]string) {
	Config = config
}

func GetConfigByKey(key string) string {
	config := getConfigByEnv(GetEnv())
	return config[key]
}

//根据环境参数获取配置
func getConfigByEnv(env string) map[string]string {
	envs, ok := Config[env]
	if !ok {
		return Config["common"]
	}

	commonConfig := Config["common"]
	for k, v := range commonConfig {
		_, ok := envs[k]
		if !ok {
			envs[k] = v
		}
	}
	return envs
}

//func ParseIniFile(filePath string) map[string]map[string]string {
//	file, err := os.Open(filePath)
//	if err != nil {
//		fmt.Println(err)
//	}

//	fi, err := file.Stat()
//	if err != nil {
//		fmt.Println(err)
//	}

//	data := make([]byte, fi.Size())
//	_, err = file.Read(data)
//	if err != nil {
//		fmt.Println(err)
//	}

//	dataStr := string(data)
//	a := lib.ArrayFilter(strings.Split(dataStr, "\n"))

//	//排除;开头的语句
//	a2 := make([]string, 0)
//	for _, v := range a {
//		if len(v) == 0 {
//			continue
//		}
//		if strings.Index(v, ";") == 0 {
//			continue
//		}
//		a2 = append(a2, v)
//	}

//	//生成最后数组
//	a3 := make(map[string]map[string]string)
//	section := ""
//	buffer := make(map[string]string)
//	for _, v := range a2 {
//		if strings.Index(v, "[") == 0 && strings.Index(v, "]") == len(v)-1 {
//			if len(buffer) > 0 {
//				a3[section] = buffer
//			}
//			section = v[1 : len(v)-1]
//			buffer = make(map[string]string)
//			a3[section] = make(map[string]string)
//			continue
//		}

//		vs := strings.Split(v, "=")
//		if len(vs) != 2 {
//			continue
//		}
//		buffer[strings.Trim(vs[0], " ")] = strings.Trim(vs[1], " ")
//	}
//}

func init() {
	Debug("init goyaf config")
}
