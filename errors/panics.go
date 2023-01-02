package errors

import (
	"fmt"
	"github.com/evilsocket/islazy/log"
)

func CheckErrorPanic(err error) {
	if err == nil {
		return
	}
	error := fmt.Errorf("%s", err)
	NotifyError("Wedding", error.Error(), "Panic triggered")
	panic(err)
}
func RecoverError() {
	if err := recover(); err != nil {
		error := fmt.Errorf("Recovered from %s", err)
		log.Info(error.Error())
		//LogToFile(error.Error())
	}
}
func NotifyError(source, err, message string) {
	// TODO
}
