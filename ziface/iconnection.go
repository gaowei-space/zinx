package ziface

import "net"

// 定义链接模块的抽象层
type IConnection interface {
	// 启动
	Start()

	// 停止链接，结束当前链接的工作
	Stop()

	// 获取当前链接的绑定socket conn
	GetTCPConnection() *net.TCPConn

	// 获取当前链接的链接ID
	GetConnID() uint32

	// 获取远程客户端的 TCP状态 IP port
	RemoteAddr() net.Addr

	// 发送数据，将数据发送给远程的客户端
	SendMsg(msgID uint32, data []byte) error

	// 设置链接属性
	SetProperty(key string, value interface{})

	// 获取链接属性
	GetProperty(key string) (interface{}, error)

	// 获取所有链接属性
	GetProperties() map[string]interface{}

	// 移除链接属性
	RemoveProperty(key string)
}

type HandleFunc func(*net.TCPConn, []byte, int) error
