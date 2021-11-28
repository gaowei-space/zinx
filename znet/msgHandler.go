package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

/**
 * 消息处理模块的实现
 */
type MsgHandler struct {
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter

	// 消息队列
	TaskQueue []chan ziface.IRequest

	// Worker数量
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// 调度、执行对应的Router消息处理方法
func (mh *MsgHandler) DoMsgHandler(req ziface.IRequest) {
	// 1. 从req中获取msgID
	handler, ok := mh.Apis[req.GetMsgID()]
	if !ok {
		fmt.Println("API msgID = ", req.GetMsgID(), " is NOT FOUND! Need Register!")
	}
	// 2. 根据msgID,调度对应的router业务即可
	handler.PreHandle(req)
	handler.Handle(req)
	handler.PostHandle(req)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1. 判断当前msg绑定的API处理方法中是否存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}
	// 2. 添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add API msgID = ", msgID, " succ!")
}

// 启动一个worker工作池
func (mh *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 当前的worker对应的channel消息队列，开辟空间，第0个worker，就用第0个channel ...
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)

		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个Wordker工作流程
func (mh *MsgHandler) startOneWorker(workID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workID, " is started...")

	// 不断的阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息过来，出列的就是一个客户端的Request，执行当前Request所绑定的业务 handler
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1. 将消息平均分配给不同的 worker
	// 根据客户端简历的ConnID来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(), " request MsgID = ", request.GetMsgID(), " to WorkerID = ", workerID)

	// 2. 将消息发送给对应的worker的taskQueue即可
	mh.TaskQueue[workerID] <- request
}
