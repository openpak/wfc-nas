package main

import (
	"openpak/wfc-common/common"
	"openpak/wfc-common/logging"
	"openpak/wfc-nas/nas"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := common.LoadTokenKey(); err != nil {
		panic(err)
	}
	config := common.GetConfig()
	logging.SetLevel(*config.LogLevel)
	config.RegisterWebhooks()
	common.ReadGameList()
	nas.StartServer(false)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	nas.Shutdown()
}
