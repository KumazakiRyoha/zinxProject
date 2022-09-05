package znet

import (
	"fmt"
	"github.com/KumazakiRyoha/zinxProject/utils"
	"github.com/KumazakiRyoha/zinxProject/ziface"
	"strconv"
)

/**
消息处理模块的实现
*/

type MsgHandle struct {
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
	// 负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作Worker池数量
	WorkPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:         make(map[uint32]ziface.IRouter),
		WorkPoolSize: utils.GlobleObj.WorkerPoolSize, //从全局配置中获取
		TaskQueue:    make([]chan ziface.IRequest, utils.GlobleObj.WorkerPoolSize),
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
	// 添加msg与api的绑定关系
	m.Apis[msgId] = router
	fmt.Print("Add api msgId = ", msgId)
}
