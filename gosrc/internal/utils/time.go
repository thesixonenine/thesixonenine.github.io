package utils

import "time"

func TsToTime(ts int64) string {
	location, _ := time.LoadLocation("Asia/Shanghai")
	return time.Unix(ts, 0).In(location).Format("2006-01-02 15:04:05")
}
