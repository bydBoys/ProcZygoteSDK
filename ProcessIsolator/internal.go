package ProcessIsolator

import (
	"errors"
	"fmt"
)

const (
	nonInit = iota
	initSuccess
	initFail
	closed
)

func (app *ProcessIsolator) checkState() error {
	switch app.state {
	case nonInit:
		return errors.New("please call ProcessIsolator.Init")
	case initSuccess:
		return nil
	case initFail:
		return errors.New("call Init occurs error, please check")
	case closed:
		return errors.New("you are using an closed ProcessIsolator")
	default:
		return fmt.Errorf("unknown internal state %d", app.state)
	}
}
