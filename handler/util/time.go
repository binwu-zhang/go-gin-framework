package util

import (
	"strconv"
	"time"
)

//  时区
const (
	TIME_LOCATION_NEWYORK = "America/New_York"    //纽约
	TIME_LOCATION_MEXICO  = "America/Mexico_City" //墨西哥
	TIME_LOCATION_CHINA   = "Asia/Shanghai"       //中国
)

const (
	TIME_FORMAT_YM          = "200601"
	TIME_FORMAT_YM_         = "2006/01"
	TIME_FORMAT_YMD_        = "2006/01/02"
	TIME_FORMAT_Y_M_D       = "2006-01-02"
	TIME_FORMAT_YMD         = "20060102"
	TIME_FORMAT_Y_M_D_H_I_S = "2006-01-02 15:04:05"
	TIME_FORMAT_Y_M_D_H_I   = "2006-01-02 15:04"
	TIME_FORMAT_M_D         = "01-02"
)

//
//GetNowDate
//  @date: 2021-11-17 10:46:01
//  @Description:
//  @param location string 时区
//  @param format string  时间字符串格式
//  @param times ...int  times[0]:增加或减少的年 times[1]:增加或减少的月 times[2]:增加或减少的日
//  @return string
//
func GetNowDate(location, format string, times ...int) string {
	if format == "" {
		format = TIME_FORMAT_Y_M_D_H_I_S
	}

	if location == "" {
		location = TIME_LOCATION_CHINA
	}
	var local, _ = time.LoadLocation(location)

	switch len(times) {
	case 1:
		return time.Now().In(local).AddDate(times[0], 0, 0).Format(format)
	case 2:
		return time.Now().In(local).AddDate(times[0], times[1], 0).Format(format)
	case 3:
		return time.Now().In(local).AddDate(times[0], times[1], times[2]).Format(format)
	default:
		return time.Now().In(local).Format(format)
	}
}

//TimestampFormat
//  @date: 2021-11-17 10:45:10
//  @Description:
//  @param timestamp int64 时间戳 10位和13位
//  @param location string  时区
//  @param format string  时间字符串格式
//  @return string
//
func TimestampFormat(timestamp int64, location, format string) string {
	if format == "" {
		format = TIME_FORMAT_Y_M_D_H_I_S
	}

	if location == "" {
		location = TIME_LOCATION_CHINA
	}

	var local, _ = time.LoadLocation(location)

	var localTime time.Time
	if len(strconv.Itoa(int(timestamp))) == 10 {
		localTime = time.Unix(timestamp, 0).In(local)
	} else if len(strconv.Itoa(int(timestamp))) == 19 {
		localTime = time.Unix(0, timestamp).In(local)
	} else {
		return ""
	}

	return localTime.Format(format)
}

//ParseDate
//  @date: 2021-12-03 17:08:01
//  @Description:
//  @param format string
//  @param location string
//  @param date string
//  @return int64
func ParseDate(format, location, date string) int64 {
	if format == "" {
		format = TIME_FORMAT_Y_M_D_H_I_S
	}

	if location == "" {
		location = TIME_LOCATION_CHINA
	}

	var local, _ = time.LoadLocation(location)

	t, _ := time.ParseInLocation(format, date, local)
	return t.Unix()
}
