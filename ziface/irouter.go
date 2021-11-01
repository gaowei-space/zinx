package ziface

type IRouter interface {
	// 处理业务前hook
	PreHandle(request IRequest)

	// 处理业务主方法
	Handle(request IRequest)

	// 处理业务后hook
	PostHandle(request IRequest)
}