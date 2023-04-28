package router

import (
	"net/http"
	"riven-gateway/config"
	"sync"
)
type endPoint struct {
	Type  string
	Host  string
	Port  int32
	Path  string
	Count int64
}

type Router struct {
	mutex sync.RWMutex
	//GrpcPath -> EndPointSlice
	endPointMap map[string][]endPoint
	//http Method#Path -> GrpcPath
	httpEndpointMap map[string]string
}

var router *Router

func buildHttpEndpointMap(bc *config.Bootstrap) (m map[string]string) {
	for _, upstream := range bc.Upstreams {
		for _, mapping := range upstream.Mappings {
			m[buildHttpEndPointMapKey(mapping.Method, mapping.HttpPath)] = mapping.RpcPath
		}
	}
	return
}

func buildHttpEndPointMapKey(method, path string) string {
	return method + "#" + path
}

func InitRouter(bc *config.Bootstrap) {
	router = &Router{
		mutex: sync.RWMutex{},
		httpEndpointMap: buildHttpEndpointMap(bc),
		endPointMap:
	}
}

func RequestAllocate(response http.ResponseWriter, request *http.Request) {

}
