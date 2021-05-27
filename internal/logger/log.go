package logger

import (
	"log"
	"runtime"
	"shellbin/internal/config"
)

func Write(mess ...interface{}) {
	if !config.Conf.Debug {
		return
	}
	_, fn, line, _ := runtime.Caller(1)
	log.Printf("[+] %s:%d %s", fn, line, mess)
}
