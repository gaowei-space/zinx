package main

import (
	"fmt"
	"io"
	"net"
	"testing"
	"zinx/znet"
)

// 只是负责datapack拆包和封包
func TestDataPack(t *testing.T) {
	// ---- 模拟的服务器 ----

	// 1.创建socketTCP
	listenner, err := net.Listen("tcp", "127.0.0.1:8686")
	if err != nil {
		fmt.Println("server listen err", err)
		return
	}
	go func() {
		// 2.从客户的读取数据，拆包处理
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
			}

			go func(conn net.Conn) {
				// 处理客户端请求
				// ---拆包过程 start---
				dp := znet.NewDataPack()
				for {
					// a. 读head，id和lenght
					headData := make([]byte, dp.GetHeadLength())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						break
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack err", err)
						return
					}
					if msgHead.GetLength() > 0 {
						// msg是有数据的，需要进行二次读取
						// b. 读data，根据length获取data内容
						msg := msgHead.(*znet.Message)
						msg.Data = make([]byte, msg.GetLength())
						// 根据length，再次从io流中读取data
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err", err)
							return
						}

						// 完整的一个消息已经读取完毕
						fmt.Printf("---> Reve ID:%d, Length:%d, Data:%v \n", msg.Id, msg.Length, string(msg.Data))
					}

				}
				// ---拆包过程 end---
			}(conn)
		}
	}()

	// ---- 模拟客户端 ----
	conn, err := net.Dial("tcp", "127.0.0.1:8686")
	if err != nil {
		fmt.Println("client dial err", err)
		return
	}

	dp := znet.NewDataPack()

	// 模拟粘包过程，封装两个msg一同发送
	// 封装第一个msg包
	msg1 := &znet.Message{
		Id:     1,
		Length: 4,
		Data:   []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err", err)
		return
	}

	// 封装第二个msg包
	msg2 := &znet.Message{
		Id:     2,
		Length: 5,
		Data:   []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 err", err)
		return
	}

	// 讲两个包粘在一起
	sendData1 = append(sendData1, sendData2...)

	// 一次性发送给服务端
	conn.Write(sendData1)

	// 客户端阻塞
	select {}
}
