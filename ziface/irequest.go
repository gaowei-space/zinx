package ziface

type IRequest interface {
	// 获取链接
	GetConnection()

	// 获取数据
	GetData()
}