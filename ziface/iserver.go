package ziface

type IServer interface {
	// 启动
	Start()
	// 停止
	Stop()
	// 运行
	Serve()
	// 添加路由
	AddRouter(msgID uint32, router IRouter)
	// 获取连接管理器
	GetConnManager() IConnManager

	// 设置链接开始时的hook
	SetOnConnStart(func(IConnection))

	// 设置链接关闭时的hook
	SetOnConnStop(func(IConnection))

	// 调用链接开始时的hook
	CallOnConnStart(IConnection)

	// 调用链接关闭时的hook
	CallOnConnStop(IConnection)
}
