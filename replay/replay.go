package replay

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

func DoReplay() {
	fmt.Println("正在执行中...(您可以按 Esc 退出)")
	stopFlag := false
	go func() {
		hook.Register(hook.KeyHold, []string{"esc"}, func(event hook.Event) {
			hook.End()
			stopFlag = true
		})
		s := hook.Start()
		<-hook.Process(s)
	}()
	bytes, _ := ioutil.ReadFile("./script.txt")
	var steps []model.Operation
	err := json.Unmarshal(bytes, &steps)
	if err != nil {
		fmt.Println("脚本反序列化失败！")
		return
	}
	if len(steps) < 1 {
		fmt.Println("录制内容为空")
	}
	for _, step := range steps {
		if stopFlag {
			return
		}
		time.Sleep(step.WaitTime)
		switch step.Type {
		case "mouseMove":
			utils.MouseMove(step)
		case "mouseDrag":
			utils.MouseDrag(step)
		case "mouseClick":
			utils.MouseClick(step)
		case "keyboardDown":
			utils.KeyboardDown(step)
		case "keyboardDownWithCtrl":
			utils.KeyboardDownWithCtrl(step)
		case "keyboardDownWithAlt":
			utils.KeyboardDownWithAlt(step)
		case "inputStr":
			utils.InputStr(step.InputStr)
			fmt.Println(step.InputStr)
		}

	}
	fmt.Println("============================================================")
}
