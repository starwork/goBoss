package utils

import (
	"testing"
	"fmt"
)

func TestMail_Send(t *testing.T) {
	d := &Mail{
		Subject: "自动回复Boss消息成功!",
		Attach:  "/Users/wuranxu/go/src/goBoss/picture/哈罗单车_吴冉旭_2018_07_11_22_57_44.png",
		Content: fmt.Sprintf(`<h4>内容: %s, 接受者公司: %s, 接受者: %s</h4>`, "我想进你你你你", "111", "222"),
	}
	d.Send()
}