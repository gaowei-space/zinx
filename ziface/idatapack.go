package ziface

/**
 * 封包，拆包 模块
 * 直接面向TCP连接中的数据流，用于处理TCP粘包问题
 */
type IDataPack interface {
	// 获取包的长度
	GetHeadLength() uint32

	// 封包
	Pack(msg IMessage) ([]byte, error)

	// 拆包
	UnPack([]byte) (IMessage, error)
}