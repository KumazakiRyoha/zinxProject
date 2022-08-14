package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	/**
	模拟服务器
	*/
	// 1.创建socket TCP
	listen, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}
	// 创建一个go 承载从客户端处理业务
	go func() {
		// 2.从客户端读取数据，拆包处理
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("server accept error", err)
		}
		go func(conn net.Conn) {
			// 处理客户端的请求
			// 拆包过程
			// 1. 第一次从conn读，把包的head读出来
			dp := NewDataPack()
			for {
				headData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData)
				if err != nil {
					fmt.Println("read head error")
					break
				}
				msgHead, err := dp.Unpack(headData)
				if err != nil {
					fmt.Println("server unpacke err", err)
					return
				}
				if msgHead.GetMsgLen() > 0 {
					// msg是有数据的，需要进行第二次读取
					msg := msgHead.(*Message)
					msg.Data = make([]byte, msg.GetMsgLen())
					// 根据dataLen的长度再次从io流中读取
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("server unpack data err: ", err)
						return
					}
					// 完整的一个消息已经读取完毕
					fmt.Println("--->Revc MsgId:", msg.Id, " dataLen:", msg.DataLen,
						" data:", string(msg.Data))
				}
			}
		}(conn)
	}()

	/**
	模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial error", err)
		return
	}

	// 创建一个封包对象dp
	dp := NewDataPack()
	// 封装两个包一起发送
	// 封装第一个包
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error", err)
		return
	}
	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'n', 'i', 'h', 'a', 'a', 'y', 'h'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg1 error", err)
		return
	}
	// 将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)

	// 一次性发送给服务端
	conn.Write(sendData1)

	// 客户端阻塞
	select {}
}
