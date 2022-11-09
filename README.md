helloworld/helloworld.proto 是定义好的服务

```bash
$ protoc --version
libprotoc 3.21.7
```

定义好服务之后，执行以下命令可以自动生成gRPC接口代码xxx.pb.go文件：

`protoc -I=. -I=$GOPATH/src/googleapis --go_out=./helloworld --go-grpc_out=require_unimplemented_servers=false:./helloworld --grpc-gateway_out=logtostderr=true:./helloworld/ helloworld/helloworld.proto` 

若报错protoc-gen-go: program not found or is not executable 说明缺少protoc-gen-go工具，安装之：

`go get -u github.com/golang/protobuf/proto`

`go get -u github.com/golang/protobuf/protoc-gen-go`

启动服务：
```bash
go run server.go userServer.go
```

启动gateway：
```bash
cd gateway
go run main.go
```

手动测试：
```bash
$ curl -X POST -d '{"name": "will", "names":["1","2","3"], "embed":{"param":"world"}}' 127.0.0.1:9090/helloworld
$ curl 127.0.0.1:9090/helloworld/name/embedparam  # 没传数组参数，可以跟下面的例子一样传
$ curl 127.0.0.1:9090/helloworld?name=hello\&names=1\&names=2\&names=3\&embed.param=world  # 注意转义符，和数组参数的传输方式
```