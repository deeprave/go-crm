package crm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type Customer struct {
	Id        int64  `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Role      string `json:"role,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Contacted bool   `json:"contacted,omitempty"`
}

type Customers []Customer

type CustomerTable struct {
	customers Customers
}

func (t *CustomerTable) ReadCustomerData(filename string) error {
	var (
		err  error
		data []byte
	)
	if data, err = os.ReadFile(filename); err == nil {
		err = json.Unmarshal(data, &t.customers)
	}
	return err
}

func (c *Customer) ToJSON() (string, error) {
	buffer := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buffer).Encode(c); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (c *Customer) FromJSON(buffer []byte) error {
	reader := bytes.NewReader(buffer)
	return json.NewDecoder(reader).Decode(c)
}

func (C *Customers) ToJSON() (string, error) {
	buffer := bytes.NewBuffer(nil)
	err := json.NewEncoder(buffer).Encode(C)
	return buffer.String(), err
}

func (C *Customers) FromJSON(buffer []byte) error {
	reader := bytes.NewReader(buffer)
	return json.NewDecoder(reader).Decode(C)
}

func (t *CustomerTable) InitCustomerTable() {
	t.customers = make(Customers, 0, 16)
}

func (t *CustomerTable) NewCustomer(name, role, email, phone string) (c *Customer) {
	customer := Customer{
		Id:        t.NextId(),
		Name:      name,
		Role:      role,
		Email:     email,
		Phone:     phone,
		Contacted: false,
	}
	index := len(t.customers)
	t.customers = append(t.customers, customer)
	return &t.customers[index]
}

func (t *CustomerTable) Count() int {
	return len(t.customers)
}

func (t *CustomerTable) NextId() int64 {
	var highestId int64 = 0
	for index := 0; index < t.Count(); index++ {
		if t.customers[index].Id > highestId {
			highestId = t.customers[index].Id
		}
	}
	return highestId + 1
}

func (t *CustomerTable) GetAllCustomers() *Customers {
	return &t.customers
}

func (t *CustomerTable) GetCustomerById(id int64) *Customer {
	for index := 0; index < len(t.customers); index++ {
		if t.customers[index].Id == id {
			return &t.customers[index]
		}
	}
	return nil
}

func (t *CustomerTable) DeleteCustomerById(id int64) (*Customer, error) {
	var (
		customer Customer
	)
	for index := 0; index < len(t.customers); index++ {
		if t.customers[index].Id == id {
			customer = t.customers[index]
			length := len(t.customers)
			copy(t.customers[index:], t.customers[index+1:])
			t.customers = t.customers[:length-1]
			return &customer, nil
		}
	}
	return nil, fmt.Errorf("customer id %d not found", id)
}

func (t *CustomerTable) UpdateCustomerById(id int64, v *Customer) (*Customer, error) {
	customer := t.GetCustomerById(id)
	if customer == nil {
		return nil, fmt.Errorf("customer id %d not found", id)
	} else {
		if v.Name != "" {
			customer.Name = v.Name
		}
		if v.Role != "" {
			customer.Role = v.Role
		}
		if v.Email != "" {
			customer.Email = v.Email
		}
		if v.Phone != "" {
			customer.Phone = v.Phone
		}
		if v.Contacted {
			customer.Contacted = v.Contacted
		}
	}
	return customer, nil
}
