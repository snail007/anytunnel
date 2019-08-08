// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// assert是对testing包的一些简单包装。方便在测试包里少写一点代码。
//
// 提供了两种操作方式：直接调用包函数；或是使用Assertion对象。
// 两种方式完全等价，可以根据自己需要，选择一种。
//  func TestAssert(t *testing.T) {
//      var v interface{} = 5
//
//      // 直接调用包函数
//      assert.True(t, v == 5, "v的值[%v]不等于5", v)
//      assert.Equal(t, 5, v, "v的值[%v]不等于5", v)
//      assert.Nil(t, v)
//
//      // 以Assertion对象方式使用
//      a := assert.New(t)
//      a.True(v==5, "v的值[%v]不等于5", v)
//      a.Equal(5, v, "v的值[%v]不等于5", v)
//      a.Nil(v)
//      a.T().Log("success")
//
//      // 以函数链的形式调用Assertion对象的方法
//      a.True(false).Equal(5,6)
//  }
package assert

// 当前库的版本号
const Version = "0.5.11.141118"
