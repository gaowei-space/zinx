package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

/**
 * 消息处理模块的实现
 */
type MsgHandler struct {
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandle() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
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
