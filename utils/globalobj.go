package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

/*
 *	存储框架全局参数，供其他模块调用
 */
type GlobalObj struct {
	TcpServer ziface.IServer // 当前zinx全局的Server对象
	Host      string         // 服务器监听的IP
	TcpPort   int            // 服务器端口
	Name      string         // 服务器名称

	// Zinx
	Version          string // 当前zinx的版本号
	MaxConn          int    // 当前zinx框架允许的最大连接数
	MaxPackageSize   uint32 // 当前zinx框架数据包最大值
	WorkerPoolSize   uint32 // worker数量
	MaxWorkerTaskLen uint32 // 每个worker对应的消息队列的任务的最大值
}

var GlobalObject *GlobalObj

// 从zinx.json加载自定义参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 提供一个init方法，初始化当前的GlobalObject
func init() {
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "v0.8",
		Host:             "0.0.0.0",
		TcpPort:          8999,
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	// 加载用户自定义参数
	GlobalObject.Reload()
}
