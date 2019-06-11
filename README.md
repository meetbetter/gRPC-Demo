# 说明

介绍RPC作用，Go中原生和第三方RPC使用方法，环境搭建方法和材料。

[GitHub示例源码](https://github.com/meetbetter/gRPC-Demo)

# RPC

远程过程调用(Remote Procedure Call)，通俗的说，RPC可以实现跨机器、跨语言调用其他计算机的程序。举个例子，我在机器A上用C语言封装了某个功能的函数，我可以通过RPC在机器B上用GO语言调用机器A上的指定函数。
RPC为C/S模型，通常使用TCP或http协议。

# Golang官方RPC



go RPC可以利用tcp或http来传递数据，可以对要传递的数据使用多种类型的编解码方式。

## net/rpc库

Golang官方的net/rpc库可以通过tcp或http传递数据，但net/rpc库使用encoding/gob进行编解码，支持tcp或http数据传输方式，由于其他语言不支持gob编解码方式，所以使用net/rpc库实现的RPC方法没办法进行跨语言调用。

### server端代码

```go
package main

import (
	"net/rpc"
	"net"
	"fmt"
	"net/http"
)

type Chen struct {
}

//rcp方法
//func (t *T) MethodName(argType T1, replyType *T2) error
func (this *Chen) GetAdd(data int, sum *int) error {

	*sum = data + 100

	return nil
}

func main() {
	//1.对象实例化
	pd := new(Chen)
	//2. rpc注册
	rpc.Register(pd)
	//3. rpc网络
	rpc.HandleHTTP()
	//4. 监听网络
	ln, err := net.Listen("tcp", "127.0.0.1:12306")
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}
	//5. 等待连接
	http.Serve(ln, nil)
}
```

### client代码

```go
package main

import (
	"net/rpc"
	"fmt"
)

func main() {
	//1. 连接服务器
	cln, err := rpc.DialHTTP("tcp", "127.0.0.1:12306")
	if err != nil {
		fmt.Println("rpc.Dial error:", err)
		return
	}
	defer cln.Close()

	//2. 调用服务器函数
	var data int
	err = cln.Call("Chen.GetAdd",10, &data)
	if err != nil {
		fmt.Println("cln.Call error:", err)
		return
	}
	//3. 打印输出
	fmt.Println("计算结果为:", data)
}
```

### 运行结果

客户端输出：```计算结果为: 110```

## net/rpc/jsonrpc库

Go官方还提供了使用json编解码的rpc库：net/rpc/jsonrpc，但是使用tcp传递数据，不能用http。

### server代码

```go
package main

import (
	"net/rpc"
	"net"
	"fmt"
	"net/rpc/jsonrpc"
)

type Chen struct {
}

//rcp方法
//func (t *T) MethodName(argType T1, replyType *T2) error
func (this *Chen) GetAdd(data int, sum *int) error {

	*sum = data + 100

	return nil
}

func main() {
	//1.对象实例化
	pd := new(Chen)
	//2. rpc注册
	rpc.Register(pd)

	//3. 监听网络
	ln, err := net.Listen("tcp", "127.0.0.1:12306")
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}

	//4. 处理客户端请求
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		go func(conn net.Conn) {
			jsonrpc.ServeConn(conn)
		}(conn)
	}
}
```

### client代码

```go
package main

import (
	"fmt"
	"net/rpc/jsonrpc"
)

func main() {
	//1. 连接服务器
	cln, err := jsonrpc.Dial("tcp", "127.0.0.1:12306")
	if err != nil {
		fmt.Println("jsonrpc.Dial error:", err)
		return
	}
	defer cln.Close()

	//2. 调用服务器函数
	var data int
	err = cln.Call("Chen.GetAdd",10, &data)
	if err != nil {
		fmt.Println("cln.Call error:", err)
		return
	}
	//3. 打印输出
	fmt.Println("计算结果为:", data)
}
```

### 运行结果

客户端输出：```计算结果为: 110```



# gRPC

所以为了真正实现跨主机、跨语言的远程调用，需要使用第三方的RPC库，推荐使用谷歌开源的gRPC。gRPC基于HTTP/2，采用protobuf进行数据编解码，压缩和传输效率更高。可以参考本人的[Go语言protobuf入门](https://github.com/meetbetter/protocol-buffer-demo)了解Go语言protobuf的环境搭建和使用。

## gRPC安装

由于不能直接访问golang官网，所以安装gPRC和go扩展包比较麻烦，可以从本人[gRPC环境包安装](https://github.com/meetbetter/gRPC-Demo)中获取压缩包。

```shell
unzip x.zip -d /GOPATH/src/golang.org/x
unzip google.golang.org.zip -d /GOPATH/src/google.golang.org
```

## gRPC环境测试

启动服务器端，

```shell
$ cd $GOPATH/src/google.golang.org/grpc/examples/helloworld/greeter_server
$ go run main.go
```

启动客户端，

```shell
$ cd $GOPATH/src/google.golang.org/grpc/examples/helloworld/greeter_client
$ go run main.go
```

如果客户端打印出`2019/06/10 15:26:12 Greeting: Hello world`字样即表示gRPC环境正常。

## 建立proto文件

```
//版本
syntax = "proto3";

//包名
package myproto;

//服务
service Hello {
    //这儿注释才有效
    rpc GetAdd(In)returns(Out);//这儿注释无效
}

//传入
message In {
    //此处1不是赋值，而是指参数序号
    int64 num = 1;
}

//传出
message Out {
    //此处1不是赋值，而是指参数序号
    int64 size = 1;
}
```

## 生成go代码

在.proto文件所在目录执行下面的指令，

```shell
protoc --go_out=plugins=grpc:./ *.proto
```

> 生成go代码时要指定plugins=grpc表示生成的是gRPC代码。

## 服务端代码

```go
package main

import (
	"fmt"
	pb "gRPC/myproto" //给package起别名
	"context"
	"net"
	"google.golang.org/grpc"
)

//1. 结构体
type Chen struct {

}

//2. 该结构体实现HelloServer interface的方法
func (this *Chen)GetAdd(ctx context.Context, In *pb.In)(*pb.Out,error)  {
	return &pb.Out{Size:In.Num+100},nil
}

func main() {
	fmt.Println("server runing...")

	//3. 创建网络
	ln, err := net.Listen("tcp", "127.0.0.1:12345")
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}
	defer ln.Close()

	//4. 创建gRPC句柄
	srv := grpc.NewServer()

	//5. 注册server
	pb.RegisterHelloServer(srv, &Chen{})

	//6. 等待网络连接
	err = srv.Serve(ln)
	if err != nil {
		fmt.Println("srv.Serve error:", err)
		return
	}

}
```

## 客户端代码

```go
package main

import (
	"google.golang.org/grpc"
	"fmt"
	pb "gRPC/myproto"
	"context"
)

func main() {
	//1 连接服务器
	conn, err := grpc.Dial("127.0.0.1:12345",grpc.WithInsecure())//grpc.WithInsecure()指定后才不会报错
	if err != nil {
		fmt.Println("grpc.Dial error:", err)
		return
	}
	defer conn.Close()

	//2 创建客户端句柄
	cln := pb.NewHelloClient(conn)

	//3 调用服务器函数(RPC)
	out,err := cln.GetAdd(context.Background(), &pb.In{Num:10})
	if err != nil {
		fmt.Println("grpc.Dial error:", err)
		return
	}

	//4 打印
	fmt.Println("得到数据:", out.Size)
}

```

## 运行测试

先后运行服务器和客户端代码，可在客户端打印输出```得到数据: 110```，说明已经成功在客户端调用服务端程序。