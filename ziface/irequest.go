package ziface

type IRequest interface {
	// 获取链接
	GetConnection() IConnection

	// 获取数据
	GetData() []byte
}
