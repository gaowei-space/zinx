package ziface

type IMessage interface {
	// 获取消息ID
	GetId() uint32
	// 获取消息长度
	GetLength() uint32
	// 获取消息内容
	GetData() []byte

	// 设置消息ID
	SetId(uint32)
	// 设置消息长度
	SetLength(uint32)
	// 设置消息内容
	SetData([]byte)
}
