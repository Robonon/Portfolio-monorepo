package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"text/template"
	c "ui/configs"
	views "ui/view_models"
)

func ServiceLibHandler(log *slog.Logger, client *http.Client, cfg *c.Config, tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("%s/service-lib", cfg.APIBaseURL)
		resp, err := client.Get(url)
		if err != nil {
			log.Error(err.Error(), "err", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Error(fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body)).Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		libViewModel := &views.ServiceLib{}
		err = json.NewDecoder(resp.Body).Decode(libViewModel)
		if err != nil {
			log.Error(err.Error(), "err", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		err = tpl.ExecuteTemplate(w, "oob-service-lib", libViewModel)
		if err != nil {
			log.Error(err.Error(), "err", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Info(fmt.Sprintf("%v %v %v", r.URL.Path, r.Method, http.StatusOK))
	}
}
