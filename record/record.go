package record

import (
	"GSAutoHSProject/model"
	"GSAutoHSProject/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	// "github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

func DoRecord() {
	fmt.Println("录制中...........按下f10结束录制")
	var steps []model.Operation
	lastTime := time.Now()
	var lastOperation model.Operation

	handleDrag := func() {
		if lastOperation.Type == "mouseDrag" {
			steps = append(steps, lastOperation)
		}
	}
	handleMove := func() {
		if lastOperation.Type == "mouseMove" {
			steps = append(steps, lastOperation)
		}
	}
	/***************** 注册事件监听器 *****************/
	// 监听鼠标拖拽事件（按住移动）
	hook.Register(hook.MouseDrag, []string{}, func(event hook.Event) {
		handleMove()
		var operation model.Operation
		operation.X = int(float64(event.X) / utils.Dpi)
		operation.Y = int(float64(event.Y) / utils.Dpi)
		operation.Type = "mouseDrag"
		operation.WaitTime = time.Now().Sub(lastTime)
		lastTime = time.Now()
		lastOperation = operation
	})
	// 监听鼠标移动事件（无按键移动）
	hook.Register(hook.MouseMove, []string{}, func(event hook.Event) {
		handleDrag()
		var operation model.Operation
		operation.X = int(float64(event.X) / utils.Dpi)
		operation.Y = int(float64(event.Y) / utils.Dpi)
		operation.Type = "mouseMove"
		operation.WaitTime = time.Now().Sub(lastTime)
		lastOperation = operation
		if operation.WaitTime > time.Millisecond*500 {
			lastTime = time.Now()
			steps = append(steps, operation)
		}
	})
	// 监听鼠标按下事件（左/右键点击）
	hook.Register(hook.MouseDown, []string{}, func(event hook.Event) {
		handleDrag()
		handleMove()
		var operation model.Operation
		if event.Button == 1 {
			operation.MouseType = "left"
		} else {
			operation.MouseType = "right"
		}
		operation.X = int(float64(event.X) / utils.Dpi)
		operation.Y = int(float64(event.Y) / utils.Dpi)
		operation.Type = "mouseClick"
		operation.WaitTime = time.Now().Sub(lastTime)
		lastTime = time.Now()
		steps = append(steps, operation)
		lastOperation = operation
	})
	// 监听F10按键（结束录制）
	hook.Register(hook.KeyHold, []string{"f10"}, func(event hook.Event) {
		hook.End() // 终止事件监听
	})
	// // 监听Ctrl组合键
	// hook.Register(hook.KeyHold, []string{"ctrl"}, func(event hook.Event) {
	// 	if event.Rawcode == 162 {
	// 		return
	// 	}
	// 	handleDrag()
	// 	handleMove()
	// 	var operation model.Operation
	// 	operation.Type = "keyboardDownWithCtrl"
	// 	operation.Key = string(event.Rawcode)
	// 	operation.WaitTime = time.Now().Sub(lastTime)
	// 	lastTime = time.Now()
	// 	steps = append(steps, operation)
	// 	lastOperation = operation
	// })

	// 监听Alt组合键
	hook.Register(hook.KeyHold, []string{"alt"}, func(event hook.Event) {
		handleDrag()
		handleMove()
		var operation model.Operation
		operation.Type = "keyboardDownWithAlt"
		operation.Key = string(event.Rawcode)
		operation.WaitTime = time.Now().Sub(lastTime)
		lastTime = time.Now()
		steps = append(steps, operation)
		lastOperation = operation
	})
	hook.Register(hook.KeyHold, []string{"shift"}, func(event hook.Event) {
		handleDrag()
		handleMove()
		var operation model.Operation
		operation.Type = "KeyboardDownWithShift"
		operation.Key = string(event.Rawcode)
		operation.WaitTime = time.Now().Sub(lastTime)
		lastTime = time.Now()
		steps = append(steps, operation)
		lastOperation = operation
	})
	// 监听普通键盘按下
	// hook.Register(hook.KeyDown, []string{}, func(event hook.Event) {
	// 	if lastOperation.Type == "keyboardDownWithCtrl" || lastOperation.Type == "keyboardDownWithAlt" {
	// 		return
	// 	}
	// 	handleDrag()
	// 	handleMove()
	// 	var operation model.Operation
	// 	// operation.Type = "keyboardDown"
	// 	operation.Type = "InputStr"
	// 	operation.Key = string(event.Keychar)
	// 	operation.WaitTime = time.Now().Sub(lastTime)
	// 	lastTime = time.Now()
	// 	steps = append(steps, operation)
	// 	lastOperation = operation
	// })

	/***************** 启动事件监听 *****************/
	s := hook.Start() // 开始捕获系统事件
	<-hook.Process(s) // 阻塞等待事件处理完成（直到调用hook.End()）

	ok := hook.AddEvents("f10")
	if ok && len(steps) > 0 {
		marshal, _ := json.Marshal(steps)
		ioutil.WriteFile("./script.txt", marshal, 0666)
		fmt.Println("脚本录制完毕并写入script.txt文件")
		fmt.Println("===========================================================")
		fmt.Printf("\n\n\n")
	} else {
		fmt.Println("脚本为空")
		fmt.Println("===========================================================")
		fmt.Printf("\n\n\n")
	}
}
