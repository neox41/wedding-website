package security

import (
	"fmt"
	"github.com/evilsocket/islazy/log"
	"mrwebsites/internal"
	"os"
	"sync"
	"time"
)
var(
	LogFileLock sync.RWMutex
)
func LogToFile(message string) {
	if internal.LogFile ==""{
		return
	}
		LogFileLock.Lock()
		defer LogFileLock.Unlock()
		f, err := os.OpenFile(internal.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != err {
			log.Error("Error opening log file")
			return
		}
		defer f.Close()
		if _, err = f.WriteString(fmt.Sprintf("%s: %s\n", time.Now().Format("2006-01-02 03:04:05"), message)); err != nil {
			log.Error("Error writing log file")
			return
		}
}
