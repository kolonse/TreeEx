// TreeEx project main.go
package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var dir = flag.String("d", "./", "-d=./ 操作目录 默认当前目录")
var ofile = flag.String("o", "out.txt", "-o=out.txt 输出文件 out.txt")
var reg = flag.String("e", "", "-e=xx 正则匹配 如果配置将会只输出指定匹配的文件列表")
var rreg = flag.String("re", "", "-re=xx 正则匹配，如果文件匹配规则那么不会将文件输出，不配置该参数或者配置为空，该参数无效，优先级比 -e参数高")
var style = flag.String("s", "linux", "-s=linux 配置输出文件路径风格，有none,linux,windows 三种，none表示按照当前系统路径格式")
var pipe = flag.Bool("p", false, "-p=false 指定pipe时会直接输出到控制台,这时文件将无效")

func formatStyle(sty, s string) string {
	if sty == "none" {
		return s
	} else if sty == "linux" {
		return strings.Replace(s, "\\", "/", len(s))
	} else if sty == "windows" {
		return strings.Replace(s, "/", "\\", len(s))
	} else {
		panic(errors.New("-s 应该只能填写三种类型，none,linux,windows"))
	}
}

func main() {
	flag.Parse()
	var fileList []string
	var regex, rregex *regexp.Regexp
	if len(*reg) != 0 {
		regex = regexp.MustCompile(*reg)
	}
	if len(*rreg) != 0 {
		rregex = regexp.MustCompile(*rreg)
	}
	absDir, err := filepath.Abs(*dir)
	if err != nil {
		panic(err)
	}
	filepath.Walk(absDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		path = path[len(absDir):]
		isMatch := true
		if regex != nil {
			isMatch = regex.MatchString(path)
		}
		if isMatch && rregex != nil {
			isMatch = !rregex.MatchString(path)
		}
		if isMatch {
			fileList = append(fileList, path)
		}
		return nil
	})
	if *pipe {
		w := bufio.NewWriter(os.Stdout)
		for _, v := range fileList {
			s := formatStyle(*style, v) + "\n"
			r := strings.NewReader(s)
			io.Copy(w, r)
		}
		w.Flush()
	} else {
		file, err := os.Create(*ofile)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		for _, v := range fileList {
			file.WriteString(formatStyle(*style, v))
			file.WriteString("\n")
		}
	}
}
