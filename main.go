/*
SSH based, concurrent, parallel, async, rest, cloud, social monitoring system proof of concept
*/
package main

import (
	log "code.google.com/p/log4go"
	"flag"
	"github.com/msurdi/alf/config"
	"github.com/msurdi/alf/db"
	"github.com/msurdi/alf/http"
	"os"
)

var (
	// Module level variables
	configPath = flag.String("config", "config.yml", "Configuration file")
	debug      = flag.Bool("debug", false, "Enable debugging output")
)

func main() {
	var err error
	flag.Parse()

	log.LoadConfiguration("logging.xml")

	// Load configuration
	alfConfig, err := config.NewConfig(*configPath)
	if err != nil {
		log.Error("Error parsing configuration: ", err.Error())
		os.Exit(1)
	}

	// Connect to the database
	alfDb := db.NewDB()
	err = alfDb.Connect(alfConfig.DbUrl)
	defer alfDb.Close()

	if err != nil {
		log.Error("Error connecting to DB: ", err.Error())
		os.Exit(1)
	}

	// Get all hosts, and all checks
	var hosts []db.Host
	err = alfDb.Hosts.FindAll(&hosts)
	if err != nil {
		log.Error("Failed to get hosts: ", err.Error())
		os.Exit(1)
	}
	log.Debug("Hosts: %v", hosts)

	var tasks []db.Task
	err = alfDb.Tasks.FindAll(&tasks)
	if err != nil {
		log.Error("Failed to get tasks: ", err.Error())
		os.Exit(1)
	}
	log.Debug("tasks: %v", tasks)

	// Run every task on every host
	for _, host := range hosts {
		for _, task := range tasks {
			runner := NewRunner(&host, &task, alfDb)
			// Launch runner in its own goroutine
			go runner.run()
		}

	}

	// Run HTTP API handlers
	httpService := http.NewHttpService(alfConfig.BindAddress, alfDb)
	go httpService.Start()

	// Block forever, exit with a SIGINT
	select {}
}
