package main

import (
	"fmt"
	"github.com/KumazakiRyoha/zinxProject/znet"
	"io"
	"net"
	"time"
)

/**
模拟客户端
*/

func main() {
	fmt.Println("client start...")
	time.Sleep(1 * time.Second)

	// 1.直接链接远程服务器，得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err,exit!")
		return
	}
	for {
		// 发送封包消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(1, []byte("ZinxV0.6 client Test Message")))
		if err != nil {
			fmt.Println("Pack error:", err)
			return
		}
		_, err = conn.Write(binaryMsg)
		if err != nil {
			fmt.Println("write error", err)
			return
		}
		// 服务器回复msg
		// 先读取流中的head部分，得到ID和dataLen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("client unpack msgHead error", err)
			break
		}
		// 再根据dataLen进行二次读取，将data读取出来
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error", err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			// msg中有数据
			// msgHead（IMessage）转换为Message
			message := msgHead.(*znet.Message)
			message.Data = make([]byte, message.GetMsgLen())

			if _, err := io.ReadFull(conn, message.Data); err != nil {
				fmt.Println("read msg data error", err)
				return
			}

			fmt.Println("---> Recv Server Msg : ID=", message.Id, ", len=",
				message.DataLen, " ,dtat=", string(message.Data))
		}
		// 阻塞一下CPU
		time.Sleep(1 * time.Second)
	}

}
