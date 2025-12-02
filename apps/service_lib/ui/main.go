package main

import (
	l "api/logger"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"text/template"
	c "ui/configs"
	h "ui/handlers"
)

func main() {
	client := http.DefaultClient
	log := l.NewLogger("UI")
	slog.SetDefault(log)
	cfg := c.NewConfig(log)

	tpl := template.Must(template.ParseGlob("views/*.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		err := tpl.ExecuteTemplate(w, "index", "")
		if err != nil {
			log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Info(fmt.Sprintf("%v %v %v", r.URL.Path, r.Method, http.StatusOK))
	})

	http.HandleFunc("/service-lib", h.ServiceLibHandler(log, client, cfg, tpl))
	http.HandleFunc("/service-info", h.ServiceInfoHandler(log, client, cfg, tpl))
	http.Handle("/", http.FileServer(http.Dir("static")))
	log.Info(fmt.Sprintf("Starting server on :%v", cfg.Port), "port", cfg.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), nil)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
