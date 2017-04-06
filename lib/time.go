package lib

import (
	"time"
)

func GetNowUnix() int64 {
	return time.Now().Unix()
}

//获取今天开始的时间戳
func GetTodayStartTimeUnix() int64 {
	t, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02")+" 00:00:00")
	return t.Unix() - 8*3600
}

//获取明天开始的时间戳
func GetTomorrowStartTimeUnix() int64 {
	return GetTodayStartTimeUnix() + 86400
}
