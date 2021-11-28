package ziface

/**
 * 消息管理抽象层
 */
type IMsgHandler interface {
	// 调度、执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)
	// 为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)
	// 启动工作池
	StartWorkerPool()
	// 将消息交给消息队列，由woker处理
	SendMsgToTaskQueue(request IRequest)
}
