package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	// 当前Conn隶属于哪个server
	TcpServer ziface.IServer

	// 当前链接的socket，TCP套接字
	Conn *net.TCPConn

	// 链接ID
	ConnID uint32

	// 是否关闭
	isClosed bool

	// 告知当前链接已经退出或停止的channel
	ExitChan chan bool

	// 无缓冲的管道，用于读写Goroutine之间的消息通信
	msgChan chan []byte

	// 消息的管理MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandler

	// 链接属性集合
	property map[string]interface{}

	// 保护链接属性的锁
	propertyLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		MsgHandler: msgHandler,
	}

	c.TcpServer.GetConnManager().Add(c)

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running...]")

	defer fmt.Println("[Reader Goroutine is exit.] ConnID =", c.ConnID)
	defer c.Stop()

	for {
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

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已配置工作池机制，写入woker处理消息
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 开启Goroutine处理消息
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// 写消息的 Goroutine,  专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")

	defer fmt.Println("[Writer Goroutine is exit.] ConnID =", c.ConnID)

	// 不断的阻塞的等待 channel 的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error,", err)
				return
			}

		case <-c.ExitChan:
			// 代表Reader已经退出，此时Writer也要退出
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start... ConnID = ", c.ConnID)

	// 启动从当前链接读数据的业务
	go c.StartReader()

	// 启动从当前链接回写的模块
	go c.StartWriter()

	// 调用创建链接之后的hook方法
	c.TcpServer.CallOnConnStart(c)

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

	// 调用销毁链接之前的hook方法
	c.TcpServer.CallOnConnStop(c)

	// 关闭链接
	c.Conn.Close()

	// 告知 Writer 关闭
	c.ExitChan <- true

	// 将当前连接从ConnManager中移除
	c.TcpServer.GetConnManager().Remove(c)

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
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

	// 写入 msg channel
	c.msgChan <- binaryMsg

	return nil
}

// 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

// 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
