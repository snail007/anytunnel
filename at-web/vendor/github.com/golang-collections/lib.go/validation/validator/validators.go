// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package validator

import (
	"github.com/caixw/lib.go/assert"
)

// 判断一个[]byte变量是否可以转换成数值。
func IsNumberBytes(bs []byte) bool {
	// 首位字符可以是[0-9.+-]
	b := bs[0]
	if (b < '0' && b > '9') && b != '+' && b != '-' && b != '.' {
		return false
	}

	hasDot := (b == '.')
	for _, b = range bs[1:] {
		if (b < '0' || b > '9') && b != '.' {
			return false
		}

		if b == '.' {
			if hasDot { // 多个小数点
				return false
			}
			hasDot = true
		}
	}
	return true
}

// 判断一个值是否可转换为数值。不支持全角数值的判断。
// 允许带一个小数点及起始的正负号，但不允许千位分隔符什么的。
func IsNumber(val interface{}) bool {
	switch v := val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	case []byte:
		return IsNumberBytes(v)
	case string:
		return IsNumberBytes([]byte(v))
	default:
		return false
	}
}

// 判断是否为空值。
// 具体参照githbu.com/caixw/lib.go/assert.IsEmtpy()
func IsEmpty(val interface{}) bool {
	return assert.IsEmpty(val)
}

// 判断是否为Nil。
// 具体参照githbu.com/caixw/lib.go/assert.IsNil()
func IsNil(val interface{}) bool {
	return assert.IsNil(val)
}
