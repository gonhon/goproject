# goproject

## env

```shell
go env -w GO111MODULE=off  
go env -w GOPATH=/home/go/goproject/

https://www.golangroadmap.com/question_bank/golang.html#%E5%88%B7%E9%A2%98%E8%AE%B0%E5%BD%95


go install  -v github.com/golang/protobuf/proto@latest
go install  -v github.com/golang/protobuf/protoc-gen-go@latest
```

## test

### 执行所有的测试函数

```go
go test xxx_test.go
```

### 执行指定测试函数

```go
go test -v -run TestXXX$ 
```

**-v      显示详细的流程**

**-run  支持正则表达式  TestXXX$  只执行TestXXX函数**