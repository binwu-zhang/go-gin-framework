package util

//
//InArray
//  @date: 2021-11-17 10:42:43
//  @Description: 字符串切片inArray
//  @param need string
//  @param needArr []string
//  @return bool
//
func InArray(need string, needArr []string) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

//UintInArray
//  @date: 2022-01-27 14:34:00
//  @Description: 切片inArray
//  @param need uint
//  @param needArr []uint
//  @return bool
func UintInArray(need uint, needArr []uint) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}
