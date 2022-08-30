package znet

import (
	"fmt"
	"github.com/KumazakiRyoha/zinxProject/ziface"
	"strconv"
)

/**
消息处理模块的实现
*/

type MsgHandle struct {
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 1. 从request中找到MsgId
	handler, ok := m.Apis[request.GetMsgId()]
	if !ok {
		fmt.Printf("Api msgId ", request.GetMsgId(), " is not't register")
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	// 判断当前msg绑定的API处理方法已经存在
	if _, ok := m.Apis[msgId]; ok {
		panic("repeat api,msgId = " + strconv.Itoa(int(msgId)))
	}
}
