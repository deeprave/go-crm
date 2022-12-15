package api

import (
	"encoding/json"
	"fmt"
	"github.com/deeprave/go-crm/crm"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

var customers crm.CustomerTable

// generic utils

func setJson(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}

func Error(writer http.ResponseWriter, errString string, status int) {
	setJson(writer)
	writer.WriteHeader(status)
	errorMessage := map[string]string{"message": errString}
	_ = json.NewEncoder(writer).Encode(errorMessage)
}

// API handlers

func getCustomers(writer http.ResponseWriter, _ *http.Request) {
	setJson(writer)
	data, _ := customers.GetAllCustomers().ToJSON()
	_, _ = writer.Write([]byte(data))
}

func getCustomer(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	if idString, ok := params["id"]; !ok {
		Error(writer, fmt.Sprintf("bad request"), http.StatusBadRequest)
	} else {
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			Error(writer, fmt.Sprintf("invalid id %s", idString), http.StatusBadRequest)
		} else {
			c := customers.GetCustomerById(id)
			if c == nil {
				Error(writer, fmt.Sprintf("unknown id %d", id), http.StatusNotFound)
			} else {
				setJson(writer)
				data, _ := c.ToJSON()
				_, _ = writer.Write([]byte(data))
			}
		}
	}
}

func addCustomer(writer http.ResponseWriter, request *http.Request) {
	var (
		err  error
		body []byte
		c    crm.Customer
	)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(request.Body)

	if body, err = io.ReadAll(request.Body); err == nil {
		if err = c.FromJSON(body); err == nil {
			n := customers.NewCustomer(c.Name, c.Role, c.Email, c.Phone)
			setJson(writer)
			data, _ := n.ToJSON()
			writer.WriteHeader(http.StatusCreated)
			_, _ = writer.Write([]byte(data))
			return
		}
	}
	Error(writer, err.Error(), http.StatusBadRequest)
}

func updateCustomer(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := fmt.Errorf("bad request")
	if idString, ok := params["id"]; ok {
		id, err := strconv.ParseInt(idString, 10, 64)
		if err == nil {
			var customer = &crm.Customer{}
			if body, err := io.ReadAll(request.Body); err == nil {
				if err = customer.FromJSON(body); err == nil {
					if customer, err = customers.UpdateCustomerById(id, customer); err == nil {
						setJson(writer)
						writer.WriteHeader(http.StatusOK)
						data, _ := customer.ToJSON()
						_, _ = writer.Write([]byte(data))
						return
					}
				}
			}
		}
	}
	Error(writer, err.Error(), http.StatusBadRequest)
}

func deleteCustomer(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := fmt.Errorf("invalid request")
	if idString, ok := params["id"]; ok {
		id, err := strconv.ParseInt(idString, 10, 64)
		if err == nil {
			customer, err := customers.DeleteCustomerById(id)
			if err == nil {
				setJson(writer)
				data, _ := customer.ToJSON()
				// 202=not yet enacted, likely to succeed, 204=no return data
				// we are returning the deleted customer however
				writer.WriteHeader(http.StatusOK)
				_, _ = writer.Write([]byte(data))
				return
			} else {
				Error(writer, err.Error(), http.StatusNotFound)
				return
			}
		}
	}
	Error(writer, err.Error(), http.StatusBadRequest)
}

func ReadCustomerData(filename string) error {
	return customers.ReadCustomerData(filename)
}

func ApiRoutes(basePath string) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc(basePath, getCustomers).Methods(http.MethodGet)
	router.HandleFunc(basePath+"/{id}", getCustomer).Methods(http.MethodGet)
	router.HandleFunc(basePath, addCustomer).Methods(http.MethodPost)
	router.HandleFunc(basePath+"/{id}", updateCustomer).Methods(http.MethodPatch, http.MethodPut)
	router.HandleFunc(basePath+"/{id}", deleteCustomer).Methods(http.MethodDelete)
	return router
}

func ApiMiddleware(router *mux.Router) *mux.Router {
	return router
}
