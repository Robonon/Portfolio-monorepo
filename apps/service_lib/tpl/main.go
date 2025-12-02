package main

import (
	l "api/logger"
	"fmt"
	"net/http"
	"os"
	c "tpl/configs"
	h "tpl/internal/handlers"
)

func main() {
	log := l.NewLogger("TPL")
	cfg := c.NewConfig(log)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello world"))
		log.Info(fmt.Sprintf("%v %v %v", r.URL.Path, r.Method, r.Response.StatusCode))
	})

	http.HandleFunc("/service", h.ServiceHandler(log))
	http.HandleFunc("/services", h.ServicesHandler(log))

	err := http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
