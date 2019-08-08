// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package assert

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

// 获取某个pc寄存器中的函数名，并去掉函数名之前的路径信息。
func funcName(pc uintptr) string {
	if pc == 0 {
		return "<无法获取函数信息>"
	}

	name := runtime.FuncForPC(pc).Name()
	arr := strings.Split(name, "/")
	return arr[len(arr)-1]
}

// 获取调用者的信息。
//
// go test输出的错误信息中，并不包含_test.go文件中的定
// 位信息，有时候很难找到在_test.go中的具体位置，此函
// 数的作用就是定位到_test.go文件中的具体位置，并返回。
// 若测试包中的函数是嵌套调用的，则有可能不正确。
func getCallerInfo() string {
	for i := 0; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			return "<无法获取调用者信息>"
		}

		basename := path.Base(file)

		// 定位以_test.go结尾的文件，认定为起始调用的测试包。
		// 8 == len("_test.go")
		l := len(basename)
		if l < 8 || (basename[l-8:l] != "_test.go") {
			continue
		}

		return " @ " + funcName(pc) + "(" + basename + ":" + strconv.Itoa(line) + ")"
	}

	return "<无法获取调用者信息>"
}

// 格式化错误提示信息。
// 优先使用msg1中的信息，若msg1为空，则使用msg2中的内容，两者格式相同。
//
// msg*中的所有参数将依次传递给fmt.Sprintf()函数，所以第一个元素的值必
// 须为string或是可转换成string的值(如[]byte,[]rune,fmt.Stringer等)
func formatMessage(msg1 []interface{}, msg2 []interface{}) string {
	msg := msg1
	if len(msg) == 0 {
		msg = msg2
	}

	if len(msg) == 0 {
		return "<未提供任何错误信息>"
	}

	format := ""
	switch v := msg[0].(type) {
	case []byte:
		format = string(v)
	case []rune:
		format = string(v)
	case string:
		format = v
	case fmt.Stringer:
		format = v.String()
	default:
		return "<无法正确转换错误提示信息>"
	}

	return fmt.Sprintf(format, msg[1:]...)
}

// 当expr条件不成立时，输出错误信息。
//
// expr 返回结果值为bool类型的表达式；
// msg1,msg2输出的错误信息，之所以提供两组信息，是方便在用户没有提供的情况下，
// 可以使用系统内部提供的信息，优先使用msg1中的信息，若不存在，则使用msg2的内容。
func assert(t *testing.T, expr bool, msg1 []interface{}, msg2 []interface{}) {
	if !expr {
		t.Error(formatMessage(msg1, msg2) + getCallerInfo())
	}
}

// 断言表达式expr为true，否则输出错误信息。
//
// args对应fmt.Printf()函数中的参数，其中args[0]对应第一个参数format，依次类推，
// 具体可参数getCallerInfo()函数的介绍。
// 其它断言函数的args参数，功能与此相同。
func True(t *testing.T, expr bool, args ...interface{}) {
	assert(t, expr, args, []interface{}{"True失败，实际值为[%T:%v]", expr, expr})
}

// 断言表达式expr为false，否则输出错误信息
func False(t *testing.T, expr bool, args ...interface{}) {
	assert(t, !expr, args, []interface{}{"False失败，实际值为[%T:%v]", expr, expr})
}

// 断言表达式expr为nil，否则输出错误信息
func Nil(t *testing.T, expr interface{}, args ...interface{}) {
	assert(t, IsNil(expr), args, []interface{}{"Nil失败，实际值为[%T:%v]", expr, expr})
}

// 断言表达式expr为非nil值，否则输出错误信息
func NotNil(t *testing.T, expr interface{}, args ...interface{}) {
	assert(t, !IsNil(expr), args, []interface{}{"NotNil失败，实际值为[%T:%v]", expr, expr})
}

// 断言v1与v2两个值相等，否则输出错误信息
func Equal(t *testing.T, v1, v2 interface{}, args ...interface{}) {
	assert(t, IsEqual(v1, v2), args, []interface{}{"Equal失败，实际值为v1=[%T:%v];v2=[%T:%v]", v1, v1, v2, v2})
}

// 断言v1与v2两个值不相等，否则输出错误信息
func NotEqual(t *testing.T, v1, v2 interface{}, args ...interface{}) {
	assert(t, !IsEqual(v1, v2), args, []interface{}{"NotEqual失败，实际值为v1=[%T:%v];v2=[%T:%v]", v1, v1, v2, v2})
}

// 断言expr的值为空(nil,"",0,false)，否则输出错误信息
func Empty(t *testing.T, expr interface{}, args ...interface{}) {
	assert(t, IsEmpty(expr), args, []interface{}{"Empty失败，实际值为[%T:%v]", expr, expr})
}

// 断言expr的值为非空(除nil,"",0,false之外)，否则输出错误信息
func NotEmpty(t *testing.T, expr interface{}, args ...interface{}) {
	assert(t, !IsEmpty(expr), args, []interface{}{"NotEmpty失败，实际值为[%T:%v]", expr, expr})
}

// 断言有错误发生，否则输出错误信息
// 传递未初始化的error值(var err error = nil)，将断言失败
func Error(t *testing.T, expr interface{}, args ...interface{}) {
	if IsNil(expr) { // 空值，必定没有错误
		assert(t, false, args, []interface{}{"Error失败，实际类型为[%T]", expr})
	} else {
		_, ok := expr.(error)
		assert(t, ok, args, []interface{}{"Error失败，实际类型为[%T]", expr})
	}
}

// 断言没有错误发生，否则输出错误信息
func NotError(t *testing.T, expr interface{}, args ...interface{}) {
	if IsNil(expr) { // 空值必定没有错误
		assert(t, true, args, []interface{}{"NotError失败，实际类型为[%T]", expr})
	} else {
		err, ok := expr.(error)
		assert(t, !ok, args, []interface{}{"NotError失败，错误信息为[%v]", err})
	}
}

// 断言文件存在，否则输出错误信息
func FileExists(t *testing.T, path string, args ...interface{}) {
	_, err := os.Stat(path)

	if err != nil && !os.IsExist(err) {
		assert(t, false, args, []interface{}{"FileExists发生以下错误：%v", err.Error()})
	}
}

// 断言文件不存在，否则输出错误信息
func FileNotExists(t *testing.T, path string, args ...interface{}) {
	_, err := os.Stat(path)
	assert(t, os.IsNotExist(err), args, []interface{}{"FileExists发生以下错误：%v", err.Error()})
}

// 断言函数会发生panic，否则输出错误信息。
func Panic(t *testing.T, fn func(), args ...interface{}) {
	has, _ := HasPanic(fn)
	assert(t, has, args, []interface{}{"并未发生panic"})
}

// 断言函数会发生panic，否则输出错误信息。
func NotPanic(t *testing.T, fn func(), args ...interface{}) {
	has, msg := HasPanic(fn)
	assert(t, !has, args, []interface{}{"发生了panic，其信息为[%]", msg})
}

// 断言container包含item的或是包含item中的所有项
// 具体函数说明可参考IsContains()
func Contains(t *testing.T, container, item interface{}, args ...interface{}) {
	assert(t, IsContains(container, item), args,
		[]interface{}{"container:[%v]并未包含item[%v]", container, item})
}

// 断言container不包含item的或是不包含item中的所有项
func NotContains(t *testing.T, container, item interface{}, args ...interface{}) {
	assert(t, !IsContains(container, item), args,
		[]interface{}{"container:[%v]包含item[%v]", container, item})
}

// 判断两个字符串相等。
//
// StringEqual()与Equal()的不同之处在于：
// StringEqual()可以以相对宽松的条件来比较字符串是否相等，
// 比如忽略大小写；忽略多余的空格等，比较方式由style参数指定。
// 若style值指定为StyleStrit，则和Equal()完全相等。
func StringEqual(t *testing.T, s1, s2 string, style int, args ...interface{}) {
	assert(t, StringIsEqual(s1, s2, style), args,
		[]interface{}{"在[%v]比较方式中s1[%v] != s2[%v]", styleString(style), s1, s2})
}

// 判断两个字符串不相等。
func StringNotEqual(t *testing.T, s1, s2 string, style int, args ...interface{}) {
	assert(t, !StringIsEqual(s1, s2, style), args,
		[]interface{}{"在[%v]比较方式中s1[%v] == s2[%v]", styleString(style), s1, s2})
}
