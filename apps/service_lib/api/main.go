package main

import (
	c "api/configs"
	h "api/handlers"
	l "api/logger"
	"net/http"
	"os"
)

func main() {
	log := l.NewLogger("API")
	cfg := c.NewConfig(log)

	client := http.DefaultClient
	http.HandleFunc("/service-lib", h.ServiceLibHandler(log, cfg, client))
	http.HandleFunc("/service-info", h.ServiceInfoHandler(log, cfg, client))

	err := http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
