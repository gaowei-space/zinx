package znet

import (
	"fmt"
	"net"
	"time"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	// 当前server的消息管理模块，用来绑定MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandler
	// 该Server的连接管理器
	ConnManager ziface.IConnManager
	// 该Server创建链接之后自动调用的hook方法
	OnConnStart func(conn ziface.IConnection)
	// 该Server停止链接之后自动调用的hook方法
	OnConnStop func(conn ziface.IConnection)
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandle(),
		ConnManager: NewConnManager(),
	}

	return s
}

func (s *Server) Start() {
	go func() {
		// 0 开启消息队列和worker工作池
		s.MsgHandler.StartWorkerPool()

		fmt.Printf("[Zinx] Server Name: %s, Listenner at IP: %s, Port: %d is starting \n", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)

		fmt.Printf("[Start] Resolve %s:%d \n", s.IP, s.Port)

		// 1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("Resolve TCP Addr Error: ", err)
			return
		}

		// 2 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("Listen TCP Error: ", err)
			return
		}

		cid := uint32(0)
		// 3 阻塞地等待客户端链接，处理客户端业务（读写）
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}

			// 判断是否超出最大连接数，如果超出则关闭当前连接
			if s.ConnManager.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端响应一个超出最大连接数的错误包
				fmt.Println("===> Too many conn, MaxConn is ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.MsgHandler)

			cid++

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// TODO 将服务器资源和状态停止
	fmt.Println("[STOP] Zinx server , name ", s.Name)

	s.ConnManager.ClearConn()
}

func (s *Server) Serve() {
	s.Start()

	// 阻塞防止退出
	for {
		time.Sleep(10 * time.Second)
	}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> Call OnConnStart():")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> Call OnConnStop():")
		s.OnConnStop(conn)
	}
}
