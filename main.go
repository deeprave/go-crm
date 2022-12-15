package main

import (
	"encoding/json"
	"fmt"
	"github.com/deeprave/go-crm/api"
	"io"
	"net/http"
	"os"
	"path"
)

func main() {
	port := 4000
	host := "localhost"

	// set up our routes and possible middleware
	router := api.ApiMiddleware(api.ApiRoutes("/customers"))

	// add a way to add data to the "database" from a local file on the server
	router.HandleFunc("/load-test-data", func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()

		var (
			body []byte
			err  error
		)
		if body, err = io.ReadAll(request.Body); err == nil {
			jsmap := map[string]string{}
			err = json.Unmarshal(body, &jsmap)
			if err == nil {
				if path := jsmap["path"]; path != "" {
					if err = api.ReadCustomerData(path); err == nil {
						writer.WriteHeader(http.StatusNoContent)
						return
					}
				}
			}
		}
		api.Error(writer, err.Error(), http.StatusBadRequest)
	}).Methods(http.MethodPost)

	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		staticPath := path.Dir("./public/index.html")
		// set header
		writer.Header().Set("Content-type", "text/html")
		http.ServeFile(writer, request, staticPath)
	})

	listenAddress := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("API listening on %s\n", listenAddress)
	fmt.Fprint(os.Stderr, http.ListenAndServe(listenAddress, router))
}
