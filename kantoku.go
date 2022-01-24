package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	loadConfig()

	if os.Getenv("FIBER_PREFORK_CHILD") == "1" {
		log.SetOutput(&NopWriter{})
	} else {
		log.SetReportCaller(true)
		log.SetFormatter(&Formatter{TimestampFormat: config.Get("kantoku.logging.time_format").(string)})
	}
}

func main() {
	initializeBroker()
	initializeServer()
}
