package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"os"
	"time"
)

var (
	RequestMaxNum = 9999
	BufferByteMax = 65535
	TotalTime     int
	FailTimes     int
	MinTime       int
	MaxTime       int
	SuccessTimes  int
	size, num     int
	timeout       int64
)

type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func main() {
	// 解析终端中输入值
	ParseArgs()

	// 当小于 2 时，即认为没有输入值，则返回使用示例
	args := os.Args
	if len(args) < 2 {
		UsageExample()
	}

	// 请求的 host
	requestHost := args[len(args)-1]

	conn, _ := net.DialTimeout("ip:icmp", requestHost, time.Duration(timeout)*time.Millisecond)

	// 执行完毕后将连接关闭
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			panic(err.Error())
		}
	}(conn)

	// 设置 ICMP 头部信息
	icmp := ICMP{8, 0, 0, 0, 0}

	fmt.Printf("\n正在 ping %s 具有 %d 字节的数据:\n", requestHost, size)

	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.BigEndian, icmp)

	data := make([]byte, size)
	buffer.Write(data)
	data = buffer.Bytes()

	if num == -1 {
		num = RequestMaxNum
	}

	for i := 0; i < num; i++ {
		icmp.SequenceNum = uint16(1) // 检验和设为0
		data[2], data[3] = byte(0), byte(0)

		data[6], data[7] = byte(icmp.SequenceNum>>8), byte(icmp.SequenceNum)

		icmp.Checksum = CheckSum(data)
		data[2], data[3] = byte(icmp.Checksum>>8), byte(icmp.Checksum) // 开始时间

		tmpTimeNow := time.Now()
		_ = conn.SetDeadline(tmpTimeNow.Add(time.Duration(time.Duration(timeout) * time.Millisecond)))

		_, err := conn.Write(data)
		if err != nil {
			fmt.Println(err.Error())
		}

		buf := make([]byte, BufferByteMax)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("请求超时了！！！")
			FailTimes++
			continue
		}

		endTime := int(time.Since(tmpTimeNow) / 1000000)

		// 处理最小、最大响应时间
		HandleEndTime(endTime)

		// 格式化本次请求信息
		fmt.Printf("来自 %s 的回复: 字节=%d 时间=%dms TTL=%d\n", requestHost, len(buf[28:n]), endTime, buf[8])

		time.Sleep(1 * time.Second)
	}

	fmt.Printf("\n%s 的 Ping 统计信息:\n", requestHost)
	fmt.Printf("    数据包: 已发送 = %d，已接收 = %d，丢失 = %d (%.2f%% 丢失)，\n", SuccessTimes+FailTimes, SuccessTimes, FailTimes, float64(FailTimes*100)/float64(SuccessTimes+FailTimes))

	if MaxTime != 0 && MinTime != int(math.MaxInt32) {
		fmt.Printf("往返行程的估计时间(以毫秒为单位):\n")
		fmt.Printf("    最短 = %dms，最长 = %dms，平均 = %dms\n", MinTime, MaxTime, TotalTime/SuccessTimes)
	}
}

func HandleEndTime(endTime int) {
	if MinTime == 0 || MinTime > endTime {
		MinTime = endTime
	}

	if MaxTime < endTime {
		MaxTime = endTime
	}

	TotalTime += endTime

	SuccessTimes++
}
