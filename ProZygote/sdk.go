package ProZygote

import (
	"errors"
	"fmt"
	"github.com/bydBoys/ProcZygoteSDK/config"
	"net/rpc"
)

// ProZygote SDK对象，创建后需调用Init函数
type ProZygote struct {
	state  int
	client *rpc.Client
}

// Init 初始化
func (app *ProZygote) Init(port string) error {
	client, err := rpc.DialHTTP("tcp", port)
	if err != nil {
		app.state = initFail
		return err
	}
	app.client = client
	app.state = initSuccess
	return nil
}

// StartProcess 启动程序
func (app *ProZygote) StartProcess(commands []string, userConfig config.UserIsolated, cgroupConfig config.CGroup) (string, error) {
	if err := app.checkState(); err != nil {
		return "", err
	}
	request := &config.StartProcRequest{
		Commands:     commands,
		UserIsolated: userConfig,
		CGroup:       cgroupConfig,
	}
	var response config.StartProcResponse
	err := app.client.Call("ProcServerImpl.StartProc", request, &response)
	if err != nil {
		return "", fmt.Errorf("call rpc ProcServerImpl.StartProc error %s", err)
	}
	if response.Error != "" {
		return "", errors.New(response.Error)
	}
	return response.UUID, nil
}

// GetProcessLog 查看某个程序是否结束，以及他的日志
func (app *ProZygote) GetProcessLog(uuid string) (bool, []string, error) {
	if err := app.checkState(); err != nil {
		return false, nil, err
	}
	request := &config.GetProcLogRequest{
		UUID: uuid,
	}
	var response config.GetProcLogResponse
	err := app.client.Call("ProcServerImpl.GetProcLog", request, &response)
	if err != nil {
		return false, nil, fmt.Errorf("call rpc  error %s", err)
	}
	if response.Error != "" {
		return false, nil, errors.New(response.Error)
	}
	return response.Exist, response.Logs, nil
}

// Destroy 关闭
func (app *ProZygote) Destroy() {
	app.state = closed
	_ = app.client.Close()
}
