package ToolsPDF

import (
	"bytes"
	"log"
	"os/exec"
	"runtime"
)

/**
libreoffice转PDF
1. 采用libreoffice转PDF方案
2. 采用命令行方式转换

注意，必须安装libreoffice，如果是windows系统，需要配置环境变量；如果是linux系统，需要安装libreoffice软件
*/

// ConvertToPDF
// @Description: 转换文件为pdf
// @param filePath 需要转换的文件
// @param outPath 转换后的PDF文件存放目录
// @return string
func ConvertToPDF(filePath string, outPath string) bool {
	// 1、拼接执行转换的命令
	commandName := ""
	var params []string
	if runtime.GOOS == "windows" {
		commandName = "cmd"
		params = []string{"/c", "soffice", "--headless", "--invisible", "--convert-to", "pdf", filePath, "--outdir", outPath}
	} else if runtime.GOOS == "linux" {
		commandName = "libreoffice"
		params = []string{"--invisible", "--headless", "--convert-to", "pdf", filePath, "--outdir", outPath}
	}
	// 开始执行转换
	if _, ok := interactiveToexec(commandName, params); ok {
		return true
	} else {
		return false
	}
}

// interactiveToexec
// @Description: 执行指定命令
// @param commandName 命令名称
// @param params 命令参数
// @return string 执行结果返回信息
// @return bool 是否执行成功
func interactiveToexec(commandName string, params []string) (string, bool) {
	cmd := exec.Command(commandName, params...)
	buf, err := cmd.Output()
	w := bytes.NewBuffer(nil)
	cmd.Stderr = w
	if err != nil {
		log.Println("Error: <", err, "> when exec command read out buffer")
		return "", false
	} else {
		return string(buf), true
	}
}
