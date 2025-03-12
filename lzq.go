package main

import (
	// "GSAutoHSProject/replay"
	"GSAutoHSProject/model"
	"GSAutoHSProject/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"path/filepath"

	hook "github.com/robotn/gohook"
)

func main() {
	stopFlag := false
	go func() {
		hook.Register(hook.KeyHold, []string{"esc"}, func(event hook.Event) {
			hook.End()
			stopFlag = true
		})
		s := hook.Start()
		<-hook.Process(s)
	}()
	// 读取所有KML文件
	kmlFiles, _ := filepath.Glob("E:\\尹文萍\\IJDE\\Hurricane_IAN_Google3D_bulid_dataset\\Hurricane_IAN_Google3D_bulid_dataset\\HurricaneIAN_Google3D_kml\\*.kml")

	//保存的视频路径
	// 创建视频保存目录
	videoDir := "E:\\尹文萍\\IJDE\\Hurricane_IAN_Google3D_bulid_dataset\\video"

	// 加载基础脚本
	baseScript, _ := ioutil.ReadFile("./script.txt")
	var steps []model.Operation
	err := json.Unmarshal(baseScript, &steps)
	if err != nil {
		fmt.Println("脚本反序列化失败！")
		return
	}
	if len(steps) < 1 {
		fmt.Println("录制内容为空")
	}

	// 数字递增计数器
	index := 1
	for _, kmlPath := range kmlFiles {
		if index < 1000 {
			continue
		}
		// 动态修改脚本中的文件路径操作
		modifiedSteps := modifyScriptSteps(steps, kmlPath, videoDir, index)
		// 计数器递增
		index++
		for _, step := range modifiedSteps {
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
			// case "keyboardDown":
			// 	utils.KeyboardDown(step)
			// case "keyboardDownWithCtrl":
			// 	utils.KeyboardDownWithCtrl(step)
			case "keyboardDownWithAlt":
				utils.KeyboardDownWithAlt(step)
			case "KeyboardDownWithShift":
				utils.KeyboardDownWithShift(step)
				// case "inputStr":
				// 	utils.InputStr(step.InputStr)
			}
		}
	}

}

// 修改脚本中的关键操作（如文件路径输入）

func modifyScriptSteps(steps []model.Operation, kmlPath string, videoDir string, index int) []model.Operation {
	videoName := fmt.Sprintf("%d.m4v", index) // 生成数字文件名
	videoPath := filepath.Join(videoDir, videoName)

	isKML := true
	isVideo := true
	for i := range steps {
		if steps[i].Type == "keyboardDownWithAlt" {
			// 替换为当前KML路径
			if isKML {
				steps[i].InputStr = kmlPath
				isKML = false
			}
		}
		if steps[i].Type == "KeyboardDownWithShift" {
			// 替换为保存视频路径
			if isVideo {
				steps[i].InputStr = videoPath
				fmt.Println(videoPath)
				isVideo = false
			}
		}
	}
	return steps
}
