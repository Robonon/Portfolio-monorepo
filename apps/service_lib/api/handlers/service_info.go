package handlers

import (
	c "api/configs"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	tpl "tpl/api"
	views "ui/view_models"
)

func ServiceInfoHandler(log *slog.Logger, cfg *c.Config, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// if multiple client calls are made, use channels later?
		service, err := getServiceDefinition(client, cfg, r.URL.Query().Get("name"))
		if err != nil || service == nil {
			log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Debug("response", "service", service)
		infoViewModel := &views.ServiceInfo{
			Description: service.Info,
			PodCount:    0,                               // This will be fetched from another service later
			Logs:        []string{"log 1, log 2, log 3"}, // This will be fetched from another service later
		}
		err = json.NewEncoder(w).Encode(infoViewModel)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Info(fmt.Sprintf("%v %v %v", r.URL.Path, r.Method, http.StatusOK))
	}
}

func getServiceDefinition(client *http.Client, cfg *c.Config, name string) (*tpl.ServiceDefinition, error) {
	url := fmt.Sprintf("%s/service?name=%s", cfg.TplServiceBaseUrl, name)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	serviceInfo := &tpl.ServiceDefinition{}
	err = json.NewDecoder(resp.Body).Decode(serviceInfo)
	if err != nil {
		return nil, err
	}
	return serviceInfo, nil
}
