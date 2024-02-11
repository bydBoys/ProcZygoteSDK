package main

import (
	"github.com/bydBoys/ProcZygoteSDK/ProZygote"
	"github.com/bydBoys/ProcZygoteSDK/config"
	"log"
	"testing"
)

func TestSDK(t *testing.T) {
	zygote := new(ProZygote.ProZygote)
	if err := zygote.Init("127.0.0.1:9963"); err != nil {
		log.Fatal(err)
		return
	}
	defer zygote.Destroy()
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
	log.Println("process uuid: ", uuid)
	exist, logs, err := zygote.GetProcessLog(uuid)
	if err != nil {
		log.Println("call GetProcessLog error ", err)
		return
	}
	log.Printf("process(%s):\n exist:%t\n logs:%s\n", uuid, exist, logs)
}
