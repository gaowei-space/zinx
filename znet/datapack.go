package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

/**
 * 封包，拆包 模块
 * 直接面向TCP连接中的数据流，用于处理TCP粘包问题
 */
type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包的长度
func (dp *DataPack) GetHeadLength() uint32 {
	// DataLength uint32(4个字节) + ID uint32(4个字节)
	return 8
}

// 封包(Length|ID|Data)
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放byte字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将dataLen写入dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetLength()); err != nil {
		return nil, err
	}
	// 将Id写入dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetId()); err != nil {
		return nil, err
	}
	// 将data写入dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包(只需要将Head信息读出来，之后再根据head信息里的data的长度，再进行一次读)
func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 只解压head信息，得到length和ID
	msg := &Message{}

	// 读 length
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Length); err != nil {
		return nil, err
	}

	// 读 id
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断length是否已经超出了最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.Length > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv!")
	}

	return msg, nil
}
