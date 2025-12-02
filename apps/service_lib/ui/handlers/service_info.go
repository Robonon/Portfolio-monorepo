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

func ServiceInfoHandler(log *slog.Logger, client *http.Client, cfg *c.Config, tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("%s/service-info?name=%s", cfg.APIBaseURL, r.URL.Query().Get("name"))
		resp, err := client.Get(url)
		if err != nil {
			log.ErrorContext(resp.Request.Context(), err.Error(), "err", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Error(fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body)).Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		infoViewModel := &views.ServiceInfo{}
		err = json.NewDecoder(resp.Body).Decode(infoViewModel)
		if err != nil {
			log.Error(err.Error(), "err", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Debug("viewmodel", "model", infoViewModel)
		err = tpl.ExecuteTemplate(w, "oob-service-info", infoViewModel)
		if err != nil {
			log.Error(err.Error(), "err", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Info(fmt.Sprintf("%v %v %v", r.URL.Path, r.Method, http.StatusOK))
	}
}
