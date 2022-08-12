package ziface

/**
IRequest接口：实际上是把客户端请求的连接信息和请求的数据包装到同一个request中
*/

type IRequest interface {
	// 得到当前的链接
	GetConnection() IConnection
	// 得到请求的消息数据
	GetData() []byte
}
