package ProcessIsolator

import (
	"errors"
	"fmt"
	"github.com/bydBoys/ProcessIsolatorSDK/config"
	"net/rpc"
)

// ProcessIsolator SDK对象，创建后需调用Init函数
type ProcessIsolator struct {
	state  int
	client *rpc.Client
}

const _version = "24.2.20"

// Init 初始化
func (app *ProcessIsolator) Init(port string) error {
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
func (app *ProcessIsolator) StartProcess(commands []string, userConfig config.UserIsolated, cgroupConfig config.CGroup) (string, error) {
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
func (app *ProcessIsolator) GetProcessLog(uuid string) (bool, []string, error) {
	if err := app.checkState(); err != nil {
		return false, nil, err
	}
	request := &config.GetProcLogRequest{
		UUID: uuid,
	}
	var response config.GetProcLogResponse
	err := app.client.Call("ProcServerImpl.GetProcLog", request, &response)
	if err != nil {
		return false, nil, fmt.Errorf("call rpc ProcServerImpl.GetProcLog error %s", err)
	}
	if response.Error != "" {
		return false, nil, errors.New(response.Error)
	}
	return response.Exist, response.Logs, nil
}

// KillProcess 尝试杀死某进程
func (app *ProcessIsolator) KillProcess(uuid string) (bool, error) {
	if err := app.checkState(); err != nil {
		return false, err
	}
	request := &config.KillProcLogRequest{
		UUID: uuid,
	}
	var response config.KillProcLogResponse
	err := app.client.Call("ProcServerImpl.KillProc", request, &response)
	if err != nil {
		return false, fmt.Errorf("call rpc ProcServerImpl.KillProc error %s", err)
	}
	if response.Error != "" {
		return false, errors.New(response.Error)
	}
	return response.Success, nil
}

// GetVersion 获取client和server的版本
func (app *ProcessIsolator) GetVersion() (string, string, error) {
	if err := app.checkState(); err != nil {
		return _version, "", err
	}

	var response string
	err := app.client.Call("ProcServerImpl.GetVersion", 1, &response)
	if err != nil {
		return _version, "", fmt.Errorf("call rpc ProcServerImpl.GetVersion error %s", err)
	}
	if response == "" {
		return _version, "", errors.New("unknown version. Check api change.")
	}
	return _version, response, nil
}

// Destroy 关闭
func (app *ProcessIsolator) Destroy() {
	app.state = closed
	_ = app.client.Close()
}
