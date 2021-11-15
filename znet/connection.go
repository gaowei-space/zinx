package znet

import (
	"errors"
	"fmt"
	"io"
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

	// 告知当前链接已经退出或停止的channel
	ExitChan chan bool

	// 该链接处理的router方法
	Router ziface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")

	defer fmt.Println("ConnID = ", c.ConnID)
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中
		// buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		// _, err := c.Conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("recv buf error", err)
		// 	c.ExitChan <- true
		// 	continue
		// }

		// 创建一个拆包解包对象
		dp := NewDataPack()

		// 读取客户端的msg head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLength())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			break
		}
		// 拆包，得到msgID和msgDataLength 放在消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		// 根据 msgDataLength 再次读取data
		var data []byte
		if msg.GetLength() > 0 {
			data = make([]byte, msg.GetLength())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		// 设置 msg data
		msg.SetData(data)

		// 得到当前链接的Request
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 执行路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start... ConnID = ", c.ConnID)

	// 启动从当前链接读数据的业务
	go c.StartReader()

	for {
		select {
		case <-c.ExitChan:
			//得到退出消息，不再阻塞
			return
		}
	}
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

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送消息，先封包，再发送
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}

	// 将data封包 length|id|data
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMessagePackage(msgID, data))
	if err != nil {
		return errors.New("pack error msg")
	}

	if _, err := c.Conn.Write(binaryMsg); err != nil {
		return errors.New("connection write error")
	}

	return nil
}
