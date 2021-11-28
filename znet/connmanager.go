package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

// 链接管理模块抽象层
type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理的连接集合
	connLock    sync.RWMutex                  // 保护连接集合的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加链接
func (cm *ConnManager) Add(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connID = ", conn.GetConnID(), " add to ConnManager succ; conn num = ", cm.Len())
}

// 删除链接
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	delete(cm.connections, conn.GetConnID())
	fmt.Println("connID = ", conn.GetConnID(), " remove from ConnManager succ; conn num = ", cm.Len())
}

// 根据ConnID获取链接
func (cm *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	// 保护共享资源 map， 加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connId]; ok {
		return conn, nil
	} else {
		fmt.Println("conn not found!")
		return nil, errors.New("conn not found!")
	}
}

// 得到当前连接总数
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

// 清除并终止所有的连接
func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	for connID, conn := range cm.connections {
		// 停止
		conn.Stop()

		// 删除
		delete(cm.connections, connID)
	}

	fmt.Println("Clear ALL connections, conn num = ", cm.Len())
}
