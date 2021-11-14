package ziface

type IServer interface {
	// 启动
	Start()
	// 停止
	Stop()
	// 运行
	Serve()
	// 添加路由
	AddRouter(router IRouter)
}
