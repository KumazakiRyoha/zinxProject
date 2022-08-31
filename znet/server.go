package znet

import (
	"fmt"
	"github.com/KumazakiRyoha/zinxProject/utils"
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
	// 当前server的消息管理模块，用来绑定MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandler
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	//TODO implement me
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println(" Add Router Success...")
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name : %s, listenner at IP : %s, Port:%d is starting\n", utils.GlobleObj.Name,
		utils.GlobleObj.Host, utils.GlobleObj.TcpPort)

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
			dealConn := NewConnection(conn, cid, s.MsgHandler)
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
func NewServer() ziface.IServer {
	return &Server{
		Name:       utils.GlobleObj.Name,
		IPServer:   "tcp4",
		IP:         utils.GlobleObj.Host,
		Port:       utils.GlobleObj.TcpPort,
		MsgHandler: NewMsgHandle(),
	}
}
