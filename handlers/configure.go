package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

// cookie name
const ConfigureCookieName = "duckclick.proxy.configure"

// expiration time
const ConfigureCookieExpiration = 1 * time.Hour

// ConfigureHandler definition
type ConfigureHandler struct {
}

// ConfigurationEntry definition
type ConfigurationEntry struct {
	URL         string `json:"url"`
	Host        string `json:"host"`
	CurrentPath string `json:"current_path"`
}

func (h *ConfigureHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	if request.Method != "POST" {
		response.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(response, `{"message": "Method Not Allowed"}`)
		return
	}

	// validate if the body is correct
	configurationEntry, err := decodeJSON(request)
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(response, `{"message": "Invalid JSON payload"}`)
		return
	}

	err = createCookie(configurationEntry, response)
	if err != nil {
		response.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(response, `{"message": "Failed to create cookie"}`)
		return
	}

	log.Infof("Configured with: %+v", configurationEntry)
	response.WriteHeader(http.StatusCreated)
	fmt.Fprint(response, `{"configured": true}`)
}

func createCookie(entry ConfigurationEntry, response http.ResponseWriter) error {
	value, err := json.Marshal(entry)

	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     ConfigureCookieName,
		Value:    base64.StdEncoding.EncodeToString(value),
		Expires:  time.Now().Add(ConfigureCookieExpiration),
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(response, cookie)
	return nil
}

func decodeJSON(request *http.Request) (ConfigurationEntry, error) {
	var entry ConfigurationEntry
	error := json.Unmarshal(streamToByte(request.Body), &entry)
	return entry, error
}

func streamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
