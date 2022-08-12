package znet

import (
	"errors"
	"fmt"
	"github.com/KumazakiRyoha/zinxProject/ziface"
	"net"
)

type Connection struct {

	// 当前链接的socket TCP 套接字
	Conn *net.TCPConn
	// 链接的IP
	ConnID uint32
	// 当前的链接状态
	isClosed bool

	// 告知当前链接已经退出 channel
	ExitChan chan bool
	//该链接处理的方法Router
	Router ziface.IRouter
}

func (c *Connection) Start() {
	fmt.Println("Conn Start... ConnID = ", c.ConnID)
	//TODO 启动从当前链接的读业务
	go c.StartReader()
	//TODO 启动从当前链接的写业务
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit,remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中，最大字节512
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			break
		}

		// 从当前conn数据的Request请求数据
		req := Request{
			conn: c,
			data: buf,
		}
		// 从路由中，找到注册绑定的conn对应的router
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)

	// 如果链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	// 关闭socket链接
	c.Conn.Close()
	// 回收资源
	close(c.ExitChan)

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

func (c *Connection) Send(data []byte) error {
	return errors.New("")
}

// 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
}
