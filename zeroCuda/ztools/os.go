package ztools

import "time"

func GetNowUnixNanoTime() int64 {
	return time.Now().UnixNano() / 1e6
}

//获取区间时间
func GetDistrictUnixNanoTime(t time.Time) int64 {
	return time.Now().UnixNano()/1e6 - t.UnixNano()/1e6
}

func GetDistrictUnixTime(t time.Time) int64 {
	return time.Now().Unix() - t.Unix()
}
