package znet

import (
	"fmt"
	ziface "github.com/KumazakiRyoha/zinxProject/ziface"
	"net"
)

// IServer接口的实现，定义一个服务器实现
type Server struct {

	// 服务器名称
	Name string
	// 服务器绑定的ip版本
	IPServer string
	// 服务器监听的IP
	IP string
	// 服务器监听的端口
	Port int
	// 给当前server添加一个router，server注册的链接对应的业务处理
	Router ziface.IRouter
}

func (s *Server) AddRouter(router ziface.IRouter) {
	//TODO implement me
	s.Router = router
	fmt.Println("Add Router Success...")
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP :%s, Port %d, is starting\n", s.IP, s.Port)

	go func() {
		// 1. 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPServer, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}
		// 2. 监听服务器的地址
		listenner, err := net.ListenTCP(s.IPServer, addr)
		if err != nil {
			fmt.Println("listen: ", s.IPServer, " err: ", err)
			return
		}
		fmt.Println("start Zinx server success, ", s.Name, " success,Listening...")
		var cid uint32
		cid = 0
		// 3. 阻塞等待的客户端链接，处理客户端链接业务
		for {
			// 如果有客户端链接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			//已经与客户端建立链接 ,做一些业务。此处做一个最基本的最大512字节长度的回显业务
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			//启动当前的链接业务处理
			go dealConn.Start()
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	// TODO 将一些服务器资源、状态或者已经开辟的连接信息进行停止或者回收

}

//
func (s *Server) Serve() {
	s.Start()

	// TODO 做一些启动服务器之后的额外业务

	// 阻塞状态
	select {}
}

// 初始化Server方法
func NewServer(name string) ziface.IServer {
	return &Server{
		Name:     name,
		IPServer: "tcp4",
		IP:       "0.0.0.0",
		Port:     8999,
		Router:   nil,
	}
}
