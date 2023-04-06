helloworld/helloworld.proto 是定义好的服务

```shell
$ protoc --version
libprotoc 3.21.7
```

定义好服务之后，执行以下命令可以自动生成gRPC接口代码xxx.pb.go文件：

`protoc -I=. -I=$GOPATH/src/googleapis --go_out=./helloworld --go-grpc_out=require_unimplemented_servers=false:./helloworld --grpc-gateway_out=logtostderr=true:./helloworld/ helloworld/helloworld.proto` 

若报错protoc-gen-go: program not found or is not executable 说明缺少protoc-gen-go工具，安装之：

`go get -u github.com/golang/protobuf/proto`

`go get -u github.com/golang/protobuf/protoc-gen-go`

启动服务：
```shell
$ go run server.go userServer.go
```

(可选一)启动gateway：
```shell
$ cd gateway
$ go run main.go
```

(可选二)启动负载均衡模式的gateway：
```shell
# 先多启动一组服务
$ go run server.go userServer.go -helloServicePort 50053 -userServicePort 50054
# 启动gateway
$ cd gateway-resolver
$ go run *.go
```

自动测试：
执行`server_test.go`即可

手动测试：
```shell
$ curl -X POST -d '{"name": "will", "names":["1","2","3"], "embed":{"param":"world"}}' 127.0.0.1:9090/helloworld
# 对下面2个get方法加了登录鉴权，不带token请求会报错
$ curl 127.0.0.1:9090/helloworld/name/embedparam  # 没传数组参数，可以跟下面的例子一样传
$ curl 127.0.0.1:9090/helloworld?name=hello\&names=1\&names=2\&names=3\&embed.param=world  # 注意转义符，和数组参数的传输方式
# 先登录
$ curl -X POST -d '{"username": "zw", "password": "password"}' 127.0.0.1:9090/login
# 携带token请求
$ curl -H 'Authorization:{上一步返回的token字段}' 127.0.0.1:9090/helloworld/name/embedparam
```