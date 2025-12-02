package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	tpl "tpl/api"
)

func ServiceHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			name := r.URL.Query().Get("name")
			if name == "" {
				log.Error("Missing name", "status code", http.StatusBadRequest)
				http.Error(w, "Missing name", http.StatusBadRequest)
				return
			}
			service, err := getService(name)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(service)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		case http.MethodPut:
			service := &tpl.ServiceDefinition{}
			err := json.NewDecoder(r.Body).Decode(service)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			// parse body
			err = putService(service)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			log.Info(fmt.Sprintf("/service PUT %v", http.StatusOK), "name", service.Name)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		log.Info(fmt.Sprintf("%v %v %v", r.URL.Path, r.Method, http.StatusOK))
	}
}

func ServicesHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services, err := getServices()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(services)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Info(fmt.Sprintf("%v %v %v", r.URL.Path, r.Method, http.StatusOK))
	}
}

func getServices() (*[]tpl.ServiceDefinition, error) {
	return &[]tpl.ServiceDefinition{
		{Name: "a", Info: "Some info"},
		{Name: "b", Info: "Some info"},
		{Name: "c", Info: "Some info"},
		{Name: "d", Info: "Some info"},
	}, nil
}

func getService(name string) (*tpl.ServiceDefinition, error) {

	service := &tpl.ServiceDefinition{
		Name: name,
		Info: "This is some info about the service",
	}
	return service, nil
}

func putService(service *tpl.ServiceDefinition) error {
	return nil
}
