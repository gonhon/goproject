package main

import (
	"flag"
	"github.com/golang/glog"
)

// go run main.go -logtostderr
func main() {
	flag.Set("log_dir", "./logs") // 输出到控制台
	flag.Set("v", "2")            // 设置日志级别

	flag.Parse()

	glog.Info("This is an info message.")
	glog.Warning("This is a warning message.")
	glog.Error("This is an error message.")
	if glog.V(2) {
		glog.Info("Starting transaction...")
	}

	glog.V(2).Info("This is a level 2 log message.")
	glog.V(1).Info("This is a level 1 log message.")
	defer glog.Flush()
}
