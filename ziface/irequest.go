package ziface

type IRequest interface {
	// 获取链接
	GetConnection() IConnection

	// 获取消息数据
	GetData() []byte

	// 获取消息ID
	GetMsgID() uint32
}
