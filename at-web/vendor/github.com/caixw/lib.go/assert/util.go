// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package assert

import (
	"bytes"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 判断一个值是否为空(0, "", false, 空数组等)。
// []string{""}空数组里套一个空字符串，不会被判断为空。
func IsEmpty(expr interface{}) bool {
	if expr == nil {
		return true
	}

	switch v := expr.(type) {
	case bool:
		return false == v
	case int:
		return 0 == v
	case int8:
		return 0 == v
	case int16:
		return 0 == v
	case int32:
		return 0 == v
	case int64:
		return 0 == v
	case uint:
		return 0 == v
	case uint8:
		return 0 == v
	case uint16:
		return 0 == v
	case uint32:
		return 0 == v
	case uint64:
		return 0 == v
	case string:
		return "" == v
	case time.Time:
		return v.IsZero()
	case *time.Time:
		return v.IsZero()
	}

	// 符合IsNil条件的，都为Empty
	ret := IsNil(expr)
	if ret {
		return true
	}

	v := reflect.ValueOf(expr)
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.Chan:
		return 0 == v.Len()
	case reflect.Ptr:
		return false
	}

	return false
}

// 判断一个值是否为nil。
// 当特定类型的变量，已经声明，但还未赋值时，也将返回true
func IsNil(expr interface{}) bool {
	if nil == expr {
		return true
	}

	v := reflect.ValueOf(expr)
	k := v.Kind()

	if (k == reflect.Chan ||
		k == reflect.Func ||
		k == reflect.Interface ||
		k == reflect.Map ||
		k == reflect.Ptr ||
		k == reflect.Slice) &&
		v.IsNil() {
		return true
	}

	return false
}

// 判断两个值是否相等。
//
// 除了通过reflect.DeepEqual()判断值是否相等之外，一些类似
// 可转换的数值也能正确判断，比如以下值也将会被判断为相等：
//  int8(5)                     == int(5)
//  []int{1,2}                  == []int8{1,2}
//  []int{1,2}                  == [2]int8{1,2}
//  []int{1,2}                  == []float32{1,2}
//  map[string]int{"1":"2":2}   == map[string]int8{"1":1,"2":2}
//
//  // map的键值不同，即使可相互转换也判断不相等。
//  map[int]int{1:1,2:2}        <> map[int8]int{1:1,2:2}
func IsEqual(v1, v2 interface{}) bool {
	if reflect.DeepEqual(v1, v2) {
		return true
	}

	vv1 := reflect.ValueOf(v1)
	vv2 := reflect.ValueOf(v2)

	// NOTE: 这里返回false，而不是true
	if !vv1.IsValid() || !vv2.IsValid() {
		return false
	}

	if vv1 == vv2 {
		return true
	}

	vv1Type := vv1.Type()
	vv2Type := vv2.Type()

	// 过滤掉已经在reflect.DeepEqual()进行处理的类型
	switch vv1Type.Kind() {
	case reflect.Struct, reflect.Ptr, reflect.Func, reflect.Interface:
		return false
	case reflect.Slice, reflect.Array:
		// vv2.Kind()与vv1的不相同
		if vv2.Kind() != reflect.Slice && vv2.Kind() != reflect.Array {
			// 虽然类型不同，但可以相互转换成vv1的，如：vv2是string，vv2是[]byte，
			if vv2Type.ConvertibleTo(vv1Type) {
				return IsEqual(vv1.Interface(), vv2.Convert(vv1Type).Interface())
			}
			return false
		}

		// reflect.DeepEqual()未考虑类型不同但是类型可转换的情况，比如：
		// []int{8,9} == []int8{8,9}，此处重新对slice和array做比较处理。
		if vv1.Len() != vv2.Len() {
			return false
		}

		for i := 0; i < vv1.Len(); i++ {
			if !IsEqual(vv1.Index(i).Interface(), vv2.Index(i).Interface()) {
				return false
			}
		}
		return true // for中所有的值比较都相等，返回true
	case reflect.Map:
		if vv2.Kind() != reflect.Map {
			return false
		}

		if vv1.IsNil() != vv2.IsNil() {
			return false
		}
		if vv1.Len() != vv2.Len() {
			return false
		}
		if vv1.Pointer() == vv2.Pointer() {
			return true
		}

		// 两个map的键名类型不同
		if vv2Type.Key().Kind() != vv1Type.Key().Kind() {
			return false
		}

		for _, index := range vv1.MapKeys() {
			if !IsEqual(vv1.MapIndex(index).Interface(), vv2.MapIndex(index).Interface()) {
				return false
			}
		}
		return true // for中所有的值比较都相等，返回true
	case reflect.String:
		if vv2.Kind() == reflect.String {
			return vv1.String() == vv2.String()
		}
		if vv2Type.ConvertibleTo(vv1Type) { // 考虑v1是string，v2是[]byte的情况
			return IsEqual(vv1.Interface(), vv2.Convert(vv1Type).Interface())
		}

		return false
	}

	if vv1Type.ConvertibleTo(vv2Type) {
		return vv2.Interface() == vv1.Convert(vv2Type).Interface()
	} else if vv2Type.ConvertibleTo(vv1Type) {
		return vv1.Interface() == vv2.Convert(vv1Type).Interface()
	}

	return false
}

// 判断fn函数是否会发生panic
// 若发生了panic，将把msg一起返回。
func HasPanic(fn func()) (has bool, msg interface{}) {
	defer func() {
		if msg = recover(); msg != nil {
			has = true
		}
	}()
	fn()

	return
}

// 判断container是否包含了item的内容。若是指针，会判断指针指向的内容，
// 但是不支持多重指针。
//
// 若container是字符串(string、[]byte和[]rune，不包含fmt.Stringer接口)，
// 都将会以字符串的形式判断其是否包含item。
// 若container是个列表(array、slice、map)则判断其元素中是否包含item中的
// 的所有项，或是item本身就是container中的一个元素。
func IsContains(container, item interface{}) bool {
	if container == nil { // nil不包含任何东西
		return false
	}

	cv := reflect.ValueOf(container)
	iv := reflect.ValueOf(item)

	if cv.Kind() == reflect.Ptr {
		cv = cv.Elem()
	}

	if iv.Kind() == reflect.Ptr {
		iv = iv.Elem()
	}

	if IsEqual(container, item) {
		return true
	}

	// 判断是字符串的情况
	switch c := cv.Interface().(type) {
	case string:
		switch i := iv.Interface().(type) {
		case string:
			return strings.Contains(c, i)
		case []byte:
			return strings.Contains(c, string(i))
		case []rune:
			return strings.Contains(c, string(i))
		case byte:
			return bytes.IndexByte([]byte(c), i) != -1
		case rune:
			return bytes.IndexRune([]byte(c), i) != -1
		}
	case []byte:
		switch i := iv.Interface().(type) {
		case string:
			return bytes.Contains(c, []byte(i))
		case []byte:
			return bytes.Contains(c, i)
		case []rune:
			return strings.Contains(string(c), string(i))
		case byte:
			return bytes.IndexByte(c, i) != -1
		case rune:
			return bytes.IndexRune(c, i) != -1
		}
	case []rune:
		switch i := iv.Interface().(type) {
		case string:
			return strings.Contains(string(c), string(i))
		case []byte:
			return strings.Contains(string(c), string(i))
		case []rune:
			return strings.Contains(string(c), string(i))
		case byte:
			return strings.IndexByte(string(c), i) != -1
		case rune:
			return strings.IndexRune(string(c), i) != -1
		}
	}

	if (cv.Kind() == reflect.Slice) || (cv.Kind() == reflect.Array) {
		if !cv.IsValid() || cv.Len() == 0 { // 空的，就不算包含另一个，即使另一个也是空值。
			return false
		}

		if !iv.IsValid() {
			return false
		}

		// item是container的一个元素
		for i := 0; i < cv.Len(); i++ {
			if IsEqual(cv.Index(i).Interface(), iv.Interface()) {
				return true
			}
		}

		// 开始判断item的元素是否与container中的元素相等。

		// 若item的长度为0，表示不包含
		if (iv.Kind() != reflect.Slice) || (iv.Len() == 0) {
			return false
		}

		// item的元素比container的元素多，必须在判断完item不是container中的一个元素之
		if iv.Len() > cv.Len() {
			return false
		}

		// 依次比较item的各个子元素是否都存在于container，且下标都相同
		ivIndex := 0
		for i := 0; i < cv.Len(); i++ {
			if IsEqual(cv.Index(i).Interface(), iv.Index(ivIndex).Interface()) {
				if (ivIndex == 0) && (i+iv.Len() > cv.Len()) {
					return false
				}
				ivIndex++
				if ivIndex == iv.Len() { // 已经遍历完iv
					return true
				}
			} else if ivIndex > 0 {
				return false
			}
		}
		return false
	} // end cv.Kind == reflect.Slice and reflect.Array

	if cv.Kind() == reflect.Map {
		if cv.Len() == 0 {
			return false
		}

		if (iv.Kind() != reflect.Map) || (iv.Len() == 0) {
			return false
		}

		if iv.Len() > cv.Len() {
			return false
		}

		// 判断所有item的项都存在于container中
		for _, key := range iv.MapKeys() {
			cvItem := iv.MapIndex(key)
			if !cvItem.IsValid() { // container中不包含该值。
				return false
			}
			if !IsEqual(cvItem.Interface(), iv.MapIndex(key).Interface()) {
				return false
			}
		}
		// for中的所有判断都成立，返回true
		return true
	}

	return false
}

const (
	StyleStrit = 1 << iota // 严格的字符串比较，会忽略其它方式
	StyleTrim              // 去掉首尾空格
	StyleSpace             // 缩减所有的空格为一个
	StyleCase              // 不区分大小写
	styleAll   = StyleTrim | StyleSpace | StyleCase
)

// 将StringIsEqual()中的Style参数转换为字符串
func styleString(style int) (ret string) {
	if style > styleAll {
		return "<invalid style:" + strconv.Itoa(style) + ">"
	}

	if (style & StyleStrit) == StyleStrit {
		return "StyleStrit"
	}

	if (style & StyleTrim) == StyleTrim {
		ret += " | StyleTrim"
	}

	if (style & StyleSpace) == StyleSpace {
		ret += " | StyleSpace"
	}

	if (style & StyleCase) == StyleCase {
		ret += " | StyleCase"
	}

	return ret[3:] // 去掉第一个|
}

var spaceReplaceRegexp = regexp.MustCompile("\\s+")

// 比较两个字符串是否相等。
// 根据第三个参数style指定比较方式，style值可以是：
//  - StyleStrit
//  - StyleTrim
//  - StyleSpace
//  - StyleCase
func StringIsEqual(s1, s2 string, style int) (ret bool) {
	// 若存在StyleStrit，则忽略其它比较属性。
	if (style & StyleStrit) == StyleStrit {
		return s1 == s2
	}

	if (style & StyleTrim) == StyleTrim {
		s1 = strings.TrimSpace(s1)
		s2 = strings.TrimSpace(s2)
	}

	if (style & StyleSpace) == StyleSpace {
		s1 = spaceReplaceRegexp.ReplaceAllString(s1, " ")
		s2 = spaceReplaceRegexp.ReplaceAllString(s2, " ")
	}

	if (style & StyleCase) == StyleCase {
		s1 = strings.ToLower(s1)
		s2 = strings.ToLower(s2)
	}

	return s1 == s2
}
