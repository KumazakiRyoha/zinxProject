package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/KumazakiRyoha/zinxProject/utils"
	"github.com/KumazakiRyoha/zinxProject/ziface"
)

// 封包、拆包具体模块

type DataPack struct{}

// 拆包封包的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	// DataLen uint32（4字节） ID uint32（4字节）
	return 8
}

// 封包方法
func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放byte字节的缓冲
	dataBuffer := bytes.NewBuffer([]byte{})
	// 将dataLen写进dataBuff中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	// 将msgId写进dataBuff中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// 将data数据写进databuff中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuffer.Bytes(), nil

}

// 拆包方法
func (d *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个存放byte字节的缓冲
	dataBuff := bytes.NewReader(binaryData)
	// 只解压head信息，得到dataLen和MsgId
	msg := &Message{}
	// 都dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断dataLen是否已经超过最大长度
	if utils.GlobleObj.MaxPackageSize > 0 && msg.DataLen > utils.GlobleObj.MaxPackageSize {
		return nil, errors.New("too large msg data recv!")
	}
	return msg, nil
}
