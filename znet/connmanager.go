package znet

import (
	"errors"
	"fmt"
	"github.com/KumazakiRyoha/zinxProject/ziface"
	"sync"
)

/**
链接管理模块
*/
type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex // 保护连接集合的读写锁
}

// 创建当前连接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加链接
func (c *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 将conn加入mao中
	c.connections[conn.GetConnID()] = conn
	fmt.Println("connId = ", conn.GetConnID(), " connnection add to ConnManager successful:conn num = ", c.Len())
}

// 删除链接
func (c *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 删除链接信息
	delete(c.connections, conn.GetConnID())
	fmt.Println("connId = ", conn.GetConnID(), " connnection remove from ConnManager successful:conn num = ", c.Len())

}

// 根据链接ID获取链接
func (c *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源map，加读锁
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not FOUND!")
	}
}

// 得到当前链接总数
func (c *ConnManager) Len() int {
	return len(c.connections)
}

// 清除并终止所有链接
func (c *ConnManager) ClearConn() {
	// 保护共享资源map，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 删除conn并停止conn
	for connId, conn := range c.connections {
		conn.Stop()
		delete(c.connections, connId)
	}
	fmt.Println("Clear all connections successful conn num = ", c.Len())
}
