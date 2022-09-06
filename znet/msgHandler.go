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

type MsgHandler struct {
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
	// 负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作Worker池数量
	WorkPoolSize uint32
}

func NewMsgHandle() *MsgHandler {
	return &MsgHandler{
		Apis:         make(map[uint32]ziface.IRouter),
		WorkPoolSize: utils.GlobleObj.WorkerPoolSize, //从全局配置中获取
		TaskQueue:    make([]chan ziface.IRequest, utils.GlobleObj.WorkerPoolSize),
	}
}

func (m *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	// 1. 从request中找到MsgId
	handler, ok := m.Apis[request.GetMsgId()]
	if !ok {
		fmt.Printf("Api msgId ", request.GetMsgId(), " is not't register")
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	// 判断当前msg绑定的API处理方法已经存在
	if _, ok := m.Apis[msgId]; ok {
		panic("repeat api,msgId = " + strconv.Itoa(int(msgId)))
	}
	// 添加msg与api的绑定关系
	m.Apis[msgId] = router
	fmt.Print("Add api msgId = ", msgId)
}

//启动一个worker工作池
func (m *MsgHandler) StartWorkerPool() {
	// 根据workerPoolSize分别开启Worker，每个Worker用一个go来承载
	for i := 0; i < int(m.WorkPoolSize); i++ {
		// 1 当前的worker对应的channel消息队列开辟弓箭，第0个worker就用第0个channel
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobleObj.MaxWorkerTaskLen)
		// 2 启动当前的worker，阻塞等待消息从channel传递进来
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (m *MsgHandler) StartOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID= ", workerId, " is started...")
	// 不断的阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息过来，出列的就是一个客户端的request，执行当前request所绑定的方法
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

// 将消息交给TaskQueue，由worker进行处理
func (m *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 将消息平均分配给不通过的worker
	// 根据客户端建立的Connid来进行分配
	workId := request.GetConnection().GetConnID() % m.WorkPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		"request MsgID = ", request.GetMsgId(), "to WorkerID = ", workId)

	// 将消息发送给对应的消息队列
	m.TaskQueue[workId] <- request
}
