helloworld/helloworld.proto 是定义好的服务

```bash
$ protoc --version
libprotoc 3.21.7
```

定义好服务之后，执行以下命令可以自动生成gRPC接口代码xxx.pb.go文件：

`protoc -I=. -I=$GOPATH/src/googleapis --go_out=./helloworld --go-grpc_out=require_unimplemented_servers=false:./helloworld helloworld/helloworld.proto` 

若报错protoc-gen-go: program not found or is not executable 说明缺少protoc-gen-go工具，安装之：

`go get -u github.com/golang/protobuf/proto`

`go get -u github.com/golang/protobuf/protoc-gen-go`

生成gRPC-Gateway源码:

`protoc -I=./helloworld -I=$GOPATH/src/googleapis --grpc-gateway_out=logtostderr=true:./helloworld/ helloworld/helloworld.proto`

启动gateway后测试：

`curl -X POST -d '{"name": "will"}' 127.0.0.1:9090/helloworld`