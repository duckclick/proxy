package handlers

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Proxy definition
type Proxy struct {
	target             *url.URL
	configurationEntry ConfigurationEntry
}

// NewProxy - proxy factory
func NewProxy(configurationEntry ConfigurationEntry, request *http.Request) *Proxy {
	targetURL, _ := url.Parse(configurationEntry.URL)
	proxyURL := &url.URL{
		Scheme: targetURL.Scheme,
		Host:   targetURL.Host,
		Path:   configurationEntry.CurrentPath + request.RequestURI,
	}

	return &Proxy{
		target:             proxyURL,
		configurationEntry: configurationEntry,
	}
}

// Handle request
func (p *Proxy) Handle(response http.ResponseWriter) {
	request, err := http.NewRequest("GET", p.target.String(), nil)
	if err != nil {
		log.WithError(err).Errorf("Error creating request")
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(response, "Error creating request")
		return
	}

	request.Header.Add("Host", p.configurationEntry.Host)
	request.Header.Add("Accept", "*/*")
	request.Header.Add("User-Agent", "duckclick/proxy")

	client := &http.Client{}
	proxyResponse, err := client.Do(request)
	if err != nil {
		log.WithError(err).Errorf("Error making request")
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(response, "Error creating request %s %+v", request.Method, p.target)
		return
	}

	defer proxyResponse.Body.Close()
	body, err := ioutil.ReadAll(proxyResponse.Body)
	if err != nil {
		log.WithError(err).Errorf("Error reading response")
		response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(response, "Error reading response %s %+v %s", request.Method, p.target, proxyResponse.Status)
		return
	}

	originalContentType := proxyResponse.Header["Content-Type"][0]

	log.Infof(
		"proxy %s %+v, content-type: %s, %s",
		request.Method,
		p.target,
		originalContentType,
		proxyResponse.Status,
	)

	response.Header().Set("Content-Type", originalContentType)
	response.WriteHeader(http.StatusOK)
	fmt.Fprint(response, string(body))
}
