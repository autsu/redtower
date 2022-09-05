package logs

import (
	"github.com/autsu/redtower/conf"
	"log"
)

var L = &Logs{}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Logs struct{}

// Printf_IfOpenDebug 如果开启了 openDebug，则打印日志
func (l *Logs) Printf_IfOpenDebug(format string, v ...any) {
	if conf.GlobalConf.OpenDebug {
		log.Printf(format, v...)
	}
}

func (l *Logs) Printf(format string, v ...any) {
	log.Printf(format, v)
}
