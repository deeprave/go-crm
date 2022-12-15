/*
 * Very simplistic in-memory "database" API
 * for customer management
 */
package crm

import (
	"strings"
	"testing"
)

var DATAFILE = "data/customers.json"

func ReadCustomers(t *testing.T) *CustomerTable {
	var customerTable = &CustomerTable{}
	if err := customerTable.ReadCustomerData(DATAFILE); err != nil {
		t.Fatalf("%s: %v", DATAFILE, err)
	}
	return customerTable
}

func TestCustomerCount(t *testing.T) {
	customerTable := ReadCustomers(t)
	count := customerTable.Count()
	if count != 14 {
		t.Errorf("customer count is %d, expected 16", count)
	}
}

func TestNextCustomerId(t *testing.T) {
	customerTable := ReadCustomers(t)
	nextId := customerTable.NextId()
	if nextId != 20 {
		t.Errorf("next customer Id is %d, expected 20", nextId)
	}
}

func TestNewCustomer(t *testing.T) {
	var name, role, email, phone = "Peter Rabbit", "teacher", "pr@bunbun.com.au", "(06) 9345.1126"
	customerTable := ReadCustomers(t)
	customer := customerTable.NewCustomer(name, role, email, phone)
	if customer.Id != 20 {
		t.Errorf("new customer Id is %d, expected 20", customer.Id)
	}
	count := customerTable.Count()
	if count != 15 {
		t.Errorf("customer count is %d, expected 15", count)
	} else if name != customer.Name {
		t.Errorf("customer name is %s, expected %s", customer.Name, name)
	} else if role != customer.Role {
		t.Errorf("customer role is %s, expected %s", customer.Role, role)
	} else if email != customer.Email {
		t.Errorf("customer email is %s, expected %s", customer.Email, email)
	} else if phone != customer.Phone {
		t.Errorf("customer phone is %s, expected %s", customer.Phone, phone)
	}
}

func TestGetAllCustomers(t *testing.T) {
	customerTable := ReadCustomers(t)

	customers := customerTable.GetAllCustomers()
	actual, expected := len(*customers), customerTable.Count()
	if actual != expected {
		t.Errorf("length of all customers incorrect: expected: %d, actual: %d", expected, actual)
	}
}

func TestGetCustomerById(t *testing.T) {
	customerTable := ReadCustomers(t)

	testCustomer := func(id int64, name, role, email, phone string) {
		customer := customerTable.GetCustomerById(id)
		if customer.Id != id {
			t.Errorf("customer id is %d, expected %d", customer.Id, id)
		} else if name != customer.Name {
			t.Errorf("customer name is %s, expected %s", customer.Name, name)
		} else if role != customer.Role {
			t.Errorf("customer role is %s, expected %s", customer.Role, role)
		} else if email != customer.Email {
			t.Errorf("customer email is %s, expected %s", customer.Email, email)
		} else if phone != customer.Phone {
			t.Errorf("customer phone is %s, expected %s", customer.Phone, phone)
		}
	}

	testCustomer(1, "Tyson Danks", "student", "TysonDanks@teleworm.us", "(07) 5398 6183")
	testCustomer(12, "Luca Adcock", "student", "luca.adcock@armyspy.com", "(08) 9063 6440")
	testCustomer(19, "Jett Roth", "student", "jroth@armyspy.com.au", "(08) 9468 2742")
}

func TestDeleteCustomerById(t *testing.T) {
	customerTable := ReadCustomers(t)

	checkCount := func(prefix string, expected int) {
		count := customerTable.Count()
		if count != expected {
			t.Errorf("customer count (%s) is %d, expected %d", prefix, count, expected)
		}
	}

	checkDelete := func(id int64, count int, name, role, email, phone string) *Customer {
		checkCount("before", count)
		customer, err := customerTable.DeleteCustomerById(id)
		if err != nil {
			t.Errorf("DeleteCustomerById: %v", err)
		} else if customer == nil {
			t.Errorf("DeleteCustomerById did not return a customer")
		} else if customer.Id != id {
			t.Errorf("DeleteCustomerById returned customer with incorrect id: expected %d, got %d", id, customer.Id)
		} else {
			count--
			checkCount("after", count)
			if customer.Id != id {
				t.Errorf("customer id is %d, expected %d", customer.Id, id)
			} else if name != customer.Name {
				t.Errorf("customer name is %s, expected %s", customer.Name, name)
			} else if role != customer.Role {
				t.Errorf("customer role is %s, expected %s", customer.Role, role)
			} else if email != customer.Email {
				t.Errorf("customer email is %s, expected %s", customer.Email, email)
			} else if phone != customer.Phone {
				t.Errorf("customer phone is %s, expected %s", customer.Phone, phone)
			}
		}
		return customer
	}

	// delete first
	checkDelete(1, 14, "Tyson Danks", "student", "TysonDanks@teleworm.us", "(07) 5398 6183")
	// delete last
	checkDelete(12, 13, "Luca Adcock", "student", "luca.adcock@armyspy.com", "(08) 9063 6440")
	// delete last
	checkDelete(19, 12, "Jett Roth", "student", "jroth@armyspy.com.au", "(08) 9468 2742")
}

func TestUpdateCustomerById(t *testing.T) {
	customerTable := ReadCustomers(t)

	checkUpdate := func(id int64, count int, v *Customer) *Customer {
		customer, err := customerTable.UpdateCustomerById(id, v)
		if err != nil {
			t.Errorf("UpdateCustomerById: %v", err)
		} else if customer == nil {
			t.Errorf("UpdateCustomerById did not return a customer")
		} else if customer.Id != id {
			t.Errorf("UpdateeCustomerById returned customer with incorrect id: expected %d, got %d", id, customer.Id)
		}
		return customer
	}

	var customer *Customer
	// update first
	customer = checkUpdate(1, 14, &Customer{Name: "Tyson Davies", Role: "teacher"})
	if customer.Name != "Tyson Davies" {
		t.Errorf("UpdateCustomerById did not update name")
	}
	if customer.Role != "teacher" {
		t.Errorf("UpdateCustomerById did not update role")
	}

	// update random (that exists)
	customer = checkUpdate(12, 13, &Customer{Phone: "(08) 9988 6044", Contacted: true})
	if customer.Phone != "(08) 9988 6044" {
		t.Errorf("UpdateCustomerById did not update phone")
	}
	if !customer.Contacted {
		t.Errorf("UpdateCustomerById did not update contacted")
	}

	// update last
	customer = checkUpdate(19, 12, &Customer{Role: "teacher", Email: "jr@armyspy.org.au"})
	if customer.Role != "teacher" {
		t.Errorf("UpdateCustomerById did not update role")
	}
	if customer.Email != "jr@armyspy.org.au" {
		t.Errorf("UpdateCustomerById did not update email")
	}
}

var (
	CustomerJson   = "{\"id\":50,\"name\":\"Bill Gates\",\"role\":\"teacher\",\"email\":\"bill.gates@microsoft.com\",\"phone\":\"(555) 555 5555\"}\n"
	CustomerRecord = Customer{
		Id:    50,
		Name:  "Bill Gates",
		Role:  "teacher",
		Email: "bill.gates@microsoft.com",
		Phone: "(555) 555 5555",
	}
	CustomersJson = "[{\"id\":50,\"name\":\"Bill Gates\",\"role\":\"teacher\",\"email\":\"bill.gates@microsoft.com\",\"phone\":\"(555) 555 5555\"}," +
		"{\"id\":51,\"name\":\"Elon Musk\",\"role\":\"director\",\"email\":\"em@tesla.com\",\"phone\":\"(999) 555 0000\"}]\n"
	CustomersRecord = Customers{
		{
			Id:    50,
			Name:  "Bill Gates",
			Role:  "teacher",
			Email: "bill.gates@microsoft.com",
			Phone: "(555) 555 5555",
		},
		{
			Id:    51,
			Name:  "Elon Musk",
			Role:  "director",
			Email: "em@tesla.com",
			Phone: "(999) 555 0000",
		},
	}
)

func TestCustomerToJSON(t *testing.T) {
	expected := CustomerJson
	customer := CustomerRecord
	result, err := customer.ToJSON()
	if err != nil {
		t.Errorf("Customer.ToJSON: %v", err)
	} else if strings.Compare(result, expected) != 0 {
		t.Errorf("Customer.ToJSON failed\n Expected (%d): %s\n   Actual (%d): %s", len(expected), expected, len(result), result)
	}
}

func TestCustomerFromJSON(t *testing.T) {
	expected := &CustomerRecord
	customer := &Customer{}
	err := customer.FromJSON([]byte(CustomerJson))
	if err != nil {
		t.Errorf("Customer.FromJSON: %v", err)
	} else if *customer != *expected {
		t.Errorf("Customer.ToJSON failed\n Expected: %v\n   Actual: %v", expected, customer)
	}
}

func TestCustomersToJSON(t *testing.T) {
	expected := CustomersJson
	customers := CustomersRecord
	result, err := customers.ToJSON()
	if err != nil {
		t.Errorf("Customers.ToJSON: %v", err)
	} else if strings.Compare(result, expected) != 0 {
		t.Errorf("Customers.ToJSON failed\n Expected (%d): %s\n   Actual (%d): %s", len(expected), expected, len(result), result)
	}
}

func TestCustomersFromJSON(t *testing.T) {
	expected := &CustomersRecord
	customers := &Customers{}
	err := customers.FromJSON([]byte(CustomersJson))
	if err != nil {
		t.Errorf("Customer.FromJSON: %v", err)
	} else if len(*customers) != len(*expected) || (*customers)[0] != (*expected)[0] || (*customers)[1] != (*expected)[1] {
		t.Errorf("Customer.ToJSON failed\n Expected: %v\n   Actual: %v", expected, customers)
	}
}
