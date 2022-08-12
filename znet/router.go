package znet

import "github.com/KumazakiRyoha/zinxProject/ziface"

// 实现router时，先嵌入BaseRouter基类，
type BaseRouter struct {
}

// 这里baseRouter的方法都为空，是因为有的Router不希望有PreHandle
// PostHandle这两个业务。所以Router全部继承BaseRouter的好处是，不需要实现
//全部Handle

func (b *BaseRouter) PreHandle(request ziface.IRequest) {}

func (b *BaseRouter) Handle(request ziface.IRequest) {}

func (b *BaseRouter) PostHandle(request ziface.IRequest) {}
