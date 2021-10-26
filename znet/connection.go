package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

type Connection struct {
	// 当前链接的socket，TCP套接字
	Conn *net.TCPConn

	// 链接ID
	ConnID uint32

	// 是否关闭
	isClosed bool

	// 当前链接所绑定的处理业务的方法API
	handleAPI ziface.HandleFunc

	// 告知当前链接已经退出或停止的channel
	ExitChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, callbackAPI ziface.HandleFunc) * Connection {
	c := &Connection{
		Conn: conn,
		ConnID: connID,
		handleAPI: callbackAPI,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")

	defer fmt.Println("ConnID = ", c.ConnID, " Reader is exit, remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf重，最大512字节
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf error", err)
			continue
		}

		// 调用当前链接所绑定的HandleAPI
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("ConnID =", c.ConnID, " handle is error", err)
			break
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start... ConnID = ", c.ConnID)

	// 启动从当前链接读数据的业务
	go c.StartReader()
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop... ConnID = ", c.ConnID)

	if c.isClosed {
		return
	}

	c.isClosed = true

	c.Conn.Close()

	close(c.ExitChan)
}

func (c *Connection) GetConnID() {

}

func (c *Connection) GetTCPConnection() {

}

func (c *Connection) RemoteAddr() {

}