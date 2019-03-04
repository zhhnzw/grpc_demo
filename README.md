helloworld/helloworld.proto 是定义好的服务


定义好服务之后，执行：protoc --go_out=plugins=grpc:. helloworld/helloworld.proto  可以自动生成gRPC接口代码xxx.pb.go文件


若报错protoc-gen-go: program not found or is not executable 说明缺少protoc-gen-go工具，安装之：


go get -u github.com/golang/protobuf/proto


go get -u github.com/golang/protobuf/protoc-gen-go


