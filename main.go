package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/duckclick/proxy/handlers"
	"net/http"
	"os"
)

func main() {
	port := getPort()
	host := fmt.Sprintf(":%s", port)
	fileServer := http.FileServer(http.Dir("public"))

	log.Infof("Starting proxy at port %s", port)

	http.Handle("/", &handlers.LoadHandler{})
	http.Handle("/__duckclick__/assets/", http.StripPrefix("/__duckclick__/assets/", fileServer))
	http.ListenAndServe(host, nil)
}

func getPort() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "7275"
	}

	return port
}
