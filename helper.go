package main

import (
	"flag"
	"fmt"
	"os"
)

func ParseArgs() {
	flag.Int64Var(&timeout, "w", 1500, "等待每次回复的超时时间(毫秒)")
	flag.IntVar(&num, "n", -1, "要发送的请求数")
	flag.IntVar(&size, "l", 32, "要发送缓冲区大小")

	flag.Parse()
}

func UsageExample() {
	fmt.Print(`
用法: ping [-t] [-n count] [-l size] [-w]
选项:
	-t             Ping 指定的主机，直到停止。
				   若要停止，请键入 Ctrl+C/Cmd+C。
	-n count       要发送的回显请求数。
	-l size        发送缓冲区大小。
	-w timeout     等待每次回复的超时时间(毫秒)。
		`)
	os.Exit(0)
}

func CheckSum(data []byte) (rt uint16) {
	var (
		sum    uint32
		length = len(data)
		index  int
	)

	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}

	if length > 0 {
		sum += uint32(data[index]) << 8
	}

	rt = uint16(sum) + uint16(sum>>16)

	return ^rt
}
