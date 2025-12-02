package handlers

import (
	c "api/configs"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	tpl "tpl/api"
	views "ui/view_models"
)

func ServiceLibHandler(log *slog.Logger, cfg *c.Config, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("%s/services", cfg.TplServiceBaseUrl)
		resp, err := client.Get(url)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		services := &[]tpl.ServiceDefinition{}
		err = json.NewDecoder(resp.Body).Decode(services)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		viewModel := &views.ServiceLib{}
		for _, s := range *services {
			viewModel.Services = append(viewModel.Services, s.Name)
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(viewModel)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Info(fmt.Sprintf("%v %v %v", r.URL.Path, r.Method, http.StatusOK))
	}
}
