package main

import log "github.com/sirupsen/logrus"

func init() {
	loadConfig()

	log.SetReportCaller(true)
	log.SetFormatter(&Formatter{TimestampFormat: config.Get("kantoku.logging.time_format").(string)})
}

func main() {
	initializeBroker()
	initializeServer()
}
