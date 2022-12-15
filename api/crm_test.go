package api

import (
	"encoding/json"
	"github.com/deeprave/go-crm/crm"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Tests for API handlers
// start by reading in a test dataset
func setupData(t *testing.T) {
	if err := ReadCustomerData("../crm/data/customers.json"); err != nil {
		t.Fatal(err)
	}
}

// some infrastructuret to help testing customer record against a fixture

func testCustomerValue(t *testing.T, customer *crm.Customer, field string, value any) {
	var fieldVal any
	switch field {
	case "id":
		fieldVal = customer.Id
	case "name":
		fieldVal = customer.Name
	case "role":
		fieldVal = customer.Role
	case "email":
		fieldVal = customer.Email
	case "phone":
		fieldVal = customer.Phone
	case "contacted":
		fieldVal = customer.Contacted
	default:
		t.Fatalf("test for unknown field %s, got %v", field, value)
	}
	if fieldVal != value {
		t.Errorf("field %s: expected %v, got %v", field, fieldVal, value)
	}
}

func testCustomerValues(t *testing.T, customer *crm.Customer, c crm.Customer) {
	testCustomerValue(t, customer, "id", c.Id)
	testCustomerValue(t, customer, "name", c.Name)
	testCustomerValue(t, customer, "role", c.Role)
	testCustomerValue(t, customer, "email", c.Email)
	testCustomerValue(t, customer, "phone", c.Phone)
	testCustomerValue(t, customer, "contacted", c.Contacted)
}

func TestGetCustomers(t *testing.T) {
	setupData(t)
	request := httptest.NewRequest(http.MethodGet, "/customer", nil)
	writer := httptest.NewRecorder()
	//
	getCustomers(writer, request)
	//
	result := writer.Result()
	defer result.Body.Close()
	if data, err := io.ReadAll(result.Body); err != nil {
		t.Errorf("unexpected read error: %v", err)
	} else {
		customers := &crm.Customers{}
		if err = json.Unmarshal(data, customers); err != nil {
			t.Errorf("unexpected json error: %v", err)
		} else {
			custLength := len(*customers)
			expected := 14
			if custLength != expected {
				t.Errorf("expected %d records, got %d", expected, custLength)
			}
		}
	}
}

func TestGetCustomer(t *testing.T) {
	setupData(t)

	urlVars := map[string]string{"id": "5"}
	request := httptest.NewRequest(http.MethodGet, "/customers/{id}", nil)
	request = mux.SetURLVars(request, urlVars)
	writer := httptest.NewRecorder()
	//
	getCustomer(writer, request)
	//
	result := writer.Result()
	defer result.Body.Close()
	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", result.StatusCode, http.StatusOK)
	}
	if data, err := io.ReadAll(result.Body); err != nil {
		t.Errorf("unexpected read error: %v", err)
	} else {
		customer := &crm.Customer{}
		if err = json.Unmarshal(data, customer); err != nil {
			t.Errorf("unexpected json error: %v", err)
		} else {
			fixture := crm.Customer{5, "Bianca Bruxner", "student", "bbruxner@dayrep.com", "(07) 4938 5904", false}
			testCustomerValues(t, customer, fixture)
		}
	}
}

func TestAddCustomer(t *testing.T) {
	setupData(t)

	reader := strings.NewReader("{\"name\":\"Bill Gates\",\"role\":\"teacher\",\"email\":\"bill.gates@microsoft.com\",\"phone\":\"(555) 555 5555\"}\n")
	request := httptest.NewRequest(http.MethodPost, "/customers/{id}", reader)
	request.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()
	//
	addCustomer(writer, request)
	//
	result := writer.Result()
	defer result.Body.Close()
	if result.StatusCode != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", result.StatusCode, http.StatusCreated)
	}
	if data, err := io.ReadAll(result.Body); err != nil {
		t.Errorf("unexpected read error: %v", err)
	} else {
		customer := &crm.Customer{}
		if err = json.Unmarshal(data, customer); err != nil {
			t.Errorf("unexpected json error: %v", err)
		} else {
			// cheat here, steal the id from the created record
			fixture := crm.Customer{customer.Id, "Bill Gates", "teacher", "bill.gates@microsoft.com", "(555) 555 5555", false}
			testCustomerValues(t, customer, fixture)
		}
	}
}

func TestUpdateCustomer(t *testing.T) {
	setupData(t)

	reader := strings.NewReader("{\"name\":\"Bill Gates\",\"role\":\"teacher\",\"email\":\"bill.gates@microsoft.com\",\"phone\":\"(555) 555 5555\"}\n")
	urlVars := map[string]string{"id": "5"}
	request := httptest.NewRequest(http.MethodPut, "/customers/{id}", reader)
	request = mux.SetURLVars(request, urlVars)
	request.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()
	//
	updateCustomer(writer, request)
	//
	result := writer.Result()
	defer result.Body.Close()
	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", result.StatusCode, http.StatusOK)
	}
	if data, err := io.ReadAll(result.Body); err != nil {
		t.Errorf("unexpected read error: %v", err)
	} else {
		customer := &crm.Customer{}
		if err = json.Unmarshal(data, customer); err != nil {
			t.Errorf("unexpected json error: %v", err)
		} else {
			fixture := crm.Customer{5, "Bill Gates", "teacher", "bill.gates@microsoft.com", "(555) 555 5555", false}
			testCustomerValues(t, customer, fixture)
		}
	}
}

func TestDeleteCustomer(t *testing.T) {
	setupData(t)

	urlVars := map[string]string{"id": "5"}
	request := httptest.NewRequest(http.MethodDelete, "/customers/{id}", nil)
	request = mux.SetURLVars(request, urlVars)
	writer := httptest.NewRecorder()
	//
	deleteCustomer(writer, request)
	//
	result := writer.Result()
	defer result.Body.Close()
	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", result.StatusCode, http.StatusOK)
	}
	if data, err := io.ReadAll(result.Body); err != nil {
		t.Errorf("unexpected read error: %v", err)
	} else {
		customer := &crm.Customer{}
		if err = json.Unmarshal(data, customer); err != nil {
			t.Errorf("unexpected json error: %v", err)
		} else {
			fixture := crm.Customer{5, "Bianca Bruxner", "student", "bbruxner@dayrep.com", "(07) 4938 5904", false}
			testCustomerValues(t, customer, fixture)

			request = httptest.NewRequest(http.MethodGet, "/customers/{id}", nil)
			request = mux.SetURLVars(request, urlVars)
			writer = httptest.NewRecorder()
			//
			getCustomer(writer, request)
			//
			result = writer.Result()
			defer result.Body.Close()
			if result.StatusCode != http.StatusNotFound {
				t.Errorf("expected status code %d, got %d", result.StatusCode, http.StatusNotFound)
			}
		}
	}
}
