package flash

import (
	"encoding/gob"
	"goblog/pkg/session"
)

// Flashes Flash 消息数组类型，用以在会话中存储 map
type Flashes map[string]interface{}

// 存入会话数据里的 key
var flashKey = "_flashes"

func init() {

	gob.Register(Flashes{})
}

// Info 添加 Info 类型的消息提示
func Info(message string) {
	addFlash("info", message)
}

// Warning 添加 Warning 类型的消息提示
func Warning(message string) {
	addFlash("warning", message)
}

// Success 添加 Success 类型的消息提示
func Success(message string) {
	addFlash("success", message)
}

// Danger 添加 Danger 类型的消息提示
func Danger(message string) {
	addFlash("danger", message)
}

// All 获取所有消息
func All() Flashes {

	val := session.Get(flashKey)
	// 读取时必须做类型检测
	flashMessages, ok := val.(Flashes)
	if !ok {
		return nil
	}
	session.Forget(flashKey)
	return flashMessages
}

// 私有方法，新增一条提示
func addFlash(key string, message string) {
	flashs := Flashes{}
	flashs[key] = message
	session.Put(flashKey, flashs)
	session.Save()

}
