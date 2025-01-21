package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	var patternFileName string
	var targetAddressFileName string
	flag.StringVar(&patternFileName, "i", "lowDimPatterns", "the pattern file")
	flag.StringVar(&targetAddressFileName, "o", "targetAddresses", "the target address file")
	flag.Parse()
	start := time.Now()
	generateAddress(patternFileName, targetAddressFileName)
	endTime := time.Now()
	// 计算代码执行时间
	elapsedTime := endTime.Sub(start)
	// 输出执行时间，转换为毫秒或其他单位
	fmt.Printf("代码总执行时间：%v\n", elapsedTime)
}
