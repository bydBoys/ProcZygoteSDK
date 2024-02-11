package ProZygote

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

func (app *ProZygote) checkState() error {
	switch app.state {
	case nonInit:
		return errors.New("pls call Init")
	case initSuccess:
		return nil
	case initFail:
		return errors.New("call Init occur error, pls check")
	case closed:
		return errors.New("you are using an closed ProZygote")
	default:
		return fmt.Errorf("unknown state %d", app.state)
	}
}
