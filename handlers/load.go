package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"html/template"
	"net/http"
	"path/filepath"
)

// LoadHandler definition
type LoadHandler struct {
}

func (h *LoadHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Headers", "content-type")

	var err error
	if request.URL.Path == "/" {
		err = serveTemplate(response)

		if err != nil {
			log.WithError(err).Errorf("Error serving template")
			response.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(response, `Error serving template`)
		}

		return
	}

	cookie, err := request.Cookie(ConfigureCookieName)
	if err != nil || cookie.Value == "" {
		log.WithError(err).Errorf("Failed to read cookie")
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(response, `Error reading configuration cookie`)
		return
	}

	configurationEntry, err := decodeBase64Cookie(cookie.Value)
	proxy := NewProxy(configurationEntry, request)
	proxy.Handle(response)
}

func serveTemplate(response http.ResponseWriter) error {
	loadTemplatePath := filepath.Join("templates", "load.html")
	loadTemplate, err := template.New("load.html").ParseFiles(loadTemplatePath)

	if err != nil {
		return errors.Wrap(err, "Failed to load template")
	}

	data := struct{}{}
	err = loadTemplate.Execute(response, data)

	if err != nil {
		return errors.Wrap(err, "Failed to run template")
	}

	return nil
}

func decodeBase64Cookie(base64Value string) (ConfigurationEntry, error) {
	var entry ConfigurationEntry
	value, err := base64.StdEncoding.DecodeString(base64Value)

	if err != nil {
		return entry, err
	}

	err = json.Unmarshal(value, &entry)
	return entry, err
}
