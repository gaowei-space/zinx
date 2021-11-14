package znet

type Message struct {
	Id     uint32 // 消息ID
	Length uint32 // 长度
	Data   []byte // 内容
}

// 获取消息ID
func (m *Message) GetId() uint32 {
	return m.Id
}

// 获取消息长度
func (m *Message) GetLength() uint32 {
	return m.Length
}

// 获取消息内容
func (m *Message) GetData() []byte {
	return m.Data
}

// 设置消息ID
func (m *Message) SetId(id uint32) {
	m.Id = id
}

// 设置消息长度
func (m *Message) SetLength(length uint32) {
	m.Length = length
}

// 设置消息内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}
