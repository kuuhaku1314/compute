package main

import (
	"bufio"
	"compute"
	"fmt"
	"os"
	"strings"
)

func main() {
	engine := compute.NewComputeEngine()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("当前启动路径为:", os.Args[0])
	fmt.Println("输入你要计算的表达式，支持输入为加减乘除括号及带小数的表达式")
	var text string
	scan := scanner.Scan()
	if scan {
		text = scanner.Text()
	} else {
		fmt.Println("输入错误")
		return
	}
	err := engine.Parse(strings.Replace(text, " ", "", -1))
	if err != nil {
		fmt.Println("解释表达式出现错误")
		fmt.Println(err)
	} else {
		result, err := engine.Run()
		if err != nil {
			fmt.Println("计算表达式出现错误")
			fmt.Println(err)
		} else {
			fmt.Println("计算的结果是:", result)
		}
	}
	fmt.Println("按下ctrl + c关闭此窗口")
	scanner.Scan()
}
