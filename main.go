package main

import (
	"cyolo-efs/api"
	"cyolo-efs/conf"
	logger "cyolo-efs/logging"
	"cyolo-efs/service"
	"cyolo-efs/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Module base module interface
type Module interface {
	Start()
	Stop()
	Title() string
}

func main() {
	cfg, err := conf.InitConf()

	if err != nil {
		log.Printf("cannot decode config: %s", err.Error())
		return
	}

	// Init storage
	fs := storage.New(cfg)

	// Init service
	srv := service.New(cfg, fs)

	// api module
	apiModel := api.New(cfg, srv)

	RunModules(apiModel, fs)
}

// RunModules runs each of the modules in a separate goroutine.
func RunModules(modules ...Module) {
	defer func() {
		for _, m := range modules {
			logger.Log.Infof("Stopping module %s", m.Title())
			m.Stop()
		}
		logger.Log.Infof("Stopped all modules")
	}()

	for _, m := range modules {
		logger.Log.Infof("Starting module %s", m.Title())
		go m.Start()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
