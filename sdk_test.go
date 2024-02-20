package main

import (
	"github.com/bydBoys/ProcessIsolatorSDK/ProcessIsolator"
	"github.com/bydBoys/ProcessIsolatorSDK/config"
	"log"
	"testing"
)

func TestSDK(t *testing.T) {
	// 初始化Zygote实例
	zygote := new(ProcessIsolator.ProcessIsolator)
	// new完第一步必须init
	if err := zygote.Init("127.0.0.1:9963"); err != nil {
		log.Fatal(err)
		return
	}
	// 不使用时记得Destroy, Destroy操作不会停止server端，只是关闭了client与server的通信
	defer zygote.Destroy()
	// 获取客户端和服务端的版本
	clientVersion, serverVersion, err := zygote.GetVersion()
	if err != nil {
		log.Fatal("call GetVersion error ", err)
		return
	}
	log.Println("client version: ", clientVersion)
	log.Println("server version: ", serverVersion)

	// 启动一个进程，执行touch ./fuckfuck命令
	uuid, err := zygote.StartProcess([]string{"touch", "./fuckfuck"}, config.UserIsolated{Enable: false}, config.CGroup{
		Enable:      false,
		CpuShare:    "",
		CpuSet:      "",
		MemoryLimit: "",
	})
	if err != nil {
		log.Println("call StartProcess error ", err)
		return
	}
	// 为了与本地的pid做区分，Zygote运行的进程使用uuid作为唯一标识
	log.Println("process uuid: ", uuid)
	// 获取某uuid的日志和是否运行结束（也可能是不存在该uuid）
	exist, logs, err := zygote.GetProcessLog(uuid)
	if err != nil {
		log.Println("call GetProcessLog error ", err)
		return
	}
	log.Printf("process(%s):\n exist:%t\n logs:%s\n", uuid, exist, logs)
	if exist {
		log.Println("still exist, try kill ", uuid)
		// 尝试杀死某进程
		success, err := zygote.KillProcess(uuid)
		if success {
			log.Println("kill success")
		}
		if err != nil {
			log.Println("kill fail: ", err.Error())
		}
	}

}
