// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package validator

// 我国现行的身份证号码有两种标准：GB11643-1989、GB11643-1999：
//
// GB11643-1989为一代身份证，从左至右分别为：
//  ------------------------------------------------------------
//  | 6位行政区域代码 | 6位出生年日期（不含世纪数）| 3位顺序码 |
//  ------------------------------------------------------------
//
// GB11643-1999为二代身份证，从左至右分别为：
//  ------------------------------------------------------------
//  | 6位行政区域代码 |  8位出生日期 |  3位顺序码 |  1位检验码 |
//  ------------------------------------------------------------

var (
	// 校验位对应的规则。
	gb11643Map = []byte{'1', '0', 'x', '9', '8', '7', '6', '5', '4', '3', '2'}

	// 前17位号码对应的权值，为一个固定数组。可由gb11643_test.getWeight()计算得到。
	gb11643Weight = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
)

// 判断一个身份证是否符合gb11643标准。
// 若是15位则当作一代身份证，仅简单地判断各位是否都是数字；
// 若是18位则当作二代身份证，会计算校验位是否正确；
// 其它位数都返回false。
func IsGb11643(val interface{}) bool {
	switch v := val.(type) {
	case string:
		return IsGb11643Bytes([]byte(v))
	case []byte:
		return IsGb11643Bytes(v)
	case []rune:
		return IsGb11643Bytes([]byte(string(v)))
	default:
		return false
	}
}

// 判断一个身份证是否符合gb11643标准。
func IsGb11643Bytes(val []byte) bool {
	if len(val) == 15 {
		// 15位，只检测是否包含非数字字符。
		for i := 0; i < 15; i++ {
			if val[i] < '0' || val[i] > '9' {
				return false
			}
		} // end for
		return true
	}

	if len(val) != 18 {
		return false
	}

	sum := 0
	for i := 0; i < 17; i++ {
		sum += (gb11643Weight[i] * int((val[i] - '0')))
	}
	if val[17] == 'X' {
		val[17] = 'x'
	}
	return gb11643Map[sum%11] == val[17]
}
