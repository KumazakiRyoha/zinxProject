package znet

import (
	"errors"
	"fmt"
	"github.com/KumazakiRyoha/zinxProject/utils"
	"github.com/KumazakiRyoha/zinxProject/ziface"
	"io"
	"net"
	"sync"
)

type Connection struct {

	// 当前conn隶属于哪个server
	TcpServer ziface.IServer
	// 当前链接的socket TCP 套接字
	Conn *net.TCPConn
	// 链接的IP
	ConnID uint32
	// 当前的链接状态
	isClosed bool

	// 告知当前链接已经退出 channel
	ExitChan chan bool

	// 无缓冲管道，用于读写Goroutine之间的消息通信
	msgChan chan []byte
	// 当前server的消息管理模块，用来绑定MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandler
	// 连接属性集合
	property map[string]interface{}
	// 保护连接树形的锁
	propertyLock sync.RWMutex
}

func (c *Connection) Start() {
	fmt.Println("Conn Start... ConnID = ", c.ConnID)
	//TODO 启动从当前链接的读业务
	go c.StartReader()
	//TODO 启动从当前链接的写业务
	go c.StartWrite()

	// 按照开发者传递进来的，创建练级之后需要调用的处理业务，执行对应函数
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit,remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {

		// 从当前conn数据的Request请求数据
		dp := NewDataPack()
		// 读取客户端的Msg Head 二进制流 8字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			break
		}
		// 拆包，得到msgId和msgDataLen放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}

		//根据dataLen，再次读取Data，放在msgData中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		msg.SetData(data)

		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobleObj.WorkerPoolSize > 0 {
			// 已经开启了工作池机制，将消息发送给worker工作池即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 从路由中，找到注册绑定的conn对应的router
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// 写消息Goroutine，专门给客户端写消息
func (c *Connection) StartWrite() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), " [conn Writer exit!]")

	// 不断地阻塞等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Sending data false,", err)
				return
			}
		case <-c.ExitChan: // 没有阻塞，代表有数据可读
			return
		}
	}
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)

	// 如果链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 调用开发者注册的销毁链接之前需要执行的业务hook函数
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket链接
	c.Conn.Close()

	// 告知writer管道已经关闭
	c.ExitChan <- true

	// 将当前链接从绑定的connMgr中删除
	c.TcpServer.GetConnMgr().Remove(c)

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)

}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 提供一个SendMsg方法，将我们需要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}

	// 将data进行封包
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id=", msgId)
		return errors.New("Pack error msg")
	}
	// 将数据放入管道中
	c.msgChan <- binaryMsg
	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	//添加一个连接树形
	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	// 删除属性
	delete(c.property, key)
}

// 初始化链接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		property:   make(map[string]interface{}),
	}
	// 将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}
