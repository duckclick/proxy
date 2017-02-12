package handlers

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"html/template"
	"net/http"
	"path/filepath"
)

// LoadHandler definition
type LoadHandler struct {
}

func (h *LoadHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	loadTemplatePath := filepath.Join("templates", "load.html")
	loadTemplate, err := template.New("load.html").ParseFiles(loadTemplatePath)

	if err != nil {
		log.WithError(err).Errorf("Error parsing template")
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(response, `{"message": "Error parsing template"}`)
		return
	}

	data := struct{}{}
	err = loadTemplate.Execute(response, data)

	if err != nil {
		log.WithError(err).Errorf("Error running template")
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(response, `{"message": "Error running template"}`)
	}
}
