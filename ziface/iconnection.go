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
	// Send(data []byte) error
}

type HandleFunc func(*net.TCPConn, []byte, int) error
