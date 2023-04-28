package config

type Bootstrap struct {
	Server    Server     `yaml:"server"`
	Upstreams []Upstream `yaml:"upstreams"`
}

type Server struct {
	Name       string `yaml:"name"`
	HostPreFix string `yaml:"hostPrefix"`
	Port       int32  `yaml:"port"`
}

type Discovery struct {
	Type string `yaml:"type"`
	Host string `yaml:"host"`
	Port int32  `yaml:"port"`
	Auth string `yaml:"auth"`
}

type Register struct {
	Host string `yaml:"host"`
	Port int32  `yaml:"port"`
}

type Upstream struct {
	Name      string         `yaml:"name"`
	Discovery Discovery      `yaml:"discovery"`
	Mappings  []RouteMapping `yaml:"mappings"`
}

type RouteMapping struct {
	//Http Method
	Method string `yaml:"method"`
	// Path is the HTTP path.
	HttpPath string `yaml:"httpPath"`
	// RpcPath is the gRPC rpc method, with format of package.service/method
	RpcPath string `yaml:"rpcPath"`
}
