helloworld/helloworld.proto 是定义好的服务


然后执行：protoc --go_out=plugins=grpc:. helloworld/helloworld.proto  生成gRPC接口代码pb文件


若报错protoc-gen-go: program not found or is not executable 说明缺少protoc-gen-go工具，安装之：


go get -u github.com/golang/protobuf/proto


go get -u github.com/golang/protobuf/protoc-gen-go


