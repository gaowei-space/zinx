package znet

import (
	"errors"
	"fmt"
	"net"
	"zinx/ziface"
)

type Server struct {
	Name string
	IPVersion string
	IP string
	Port int
}

func CallbackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn Handle] CallbackToClient...")

	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf error", err)
		return errors.New("CallbackToClient Error")
	}
}

func (s *Server) Start() {

	go func () {
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
		
		var cid uint32
		cid = 0
		// 3 阻塞地等待客户端链接，处理客户端业务（读写）
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
                continue
			}

			dealConn := NewConnection(conn, cid, CallbackToClient)

			cid++

			go dealConn.Start()

			/* go func () {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("Read err ", err)
                		continue
					}

					fmt.Printf("client send: %s, cnt = %d \n", buf, cnt)

					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("Write err ", err)
                		continue
					}
				}	
			}() */
		}

	}()
}

func (s *Server) Stop() {
	// TODO 将服务器资源和状态停止
	fmt.Println("[STOP] Zinx server , name " , s.Name)
}

func (s *Server) Serve() {
	s.Start()

	// 阻塞防止退出
	select{}
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name: name,
		IPVersion: "tcp4",
		IP: "0.0.0.0",
		Port: 9999,
	}

	return s
}