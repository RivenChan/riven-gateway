package router

import "net/http"

var (
	EndPoints = []EndPoint
)

// EndPoint 具体的端点
type EndPoint struct {
	Host  string
	Port  int32
	path  string
	count int64
}

func init() {
	//从 etcd 获取数据 初始化Router
}

func RequestAllocate(response http.ResponseWriter, request *http.Request) {

}
