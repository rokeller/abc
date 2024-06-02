package main

import (
	"github.com/golang/glog"
	"github.com/rokeller/abc/cmd"
)

func main() {
	cmd.Execute()
	glog.Flush()
}
