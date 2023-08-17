package main

import (
	"docker-cleaner/cmd"
)

func main() {
	// Connect do docker daemon
	//dClient := docker.DockerClient()

	cmd.Execute()

	//ch := make(chan int, len(cfg.Watch))
	//for _, wFile := range cfg.Watch {
	//	log.Println(`Spawn inotify watcher for file:`, wFile.Path)
	//	go watcher.Watch(wFile, dClient, false)
	//}
	//
	//<-ch
}
