# Customers API
Crud API for a CRM.

## URL
This project can be viewed and downloaded from the following url: https://github.com/deeprave/go-crm.

## Description
This project provides a basic API written in Go as a backend for a CRM.

## Pre-requisites
- go version 1.18 or later

## Server startup
After checking into the root folder for this project, the following command starts the server:
```bash
go run main.go
```
The server is bound to localhost, port 4000 as hard-coded values.
No other command line options are currently available.

## CRM Domain
This api provides access only to the customer list, the core table in a CRM.
It provides the ability to:
- create new customers   `POST /customers`
- display all customers  `GET /customers`
- display a specific customer `GET /customers/{id}`
- update a specific customer `PUT /customers/{id}`
- delete a specific customer `DELETE /customers/{id}`

Each customer record consists of an Id (assigned on creation), a name, role, email phone number
and a "sticky" (stays true once set) contacted field that indicates whether that customer has
been contacted.

## Go libraries
This project uses:
- gorilla/mux

as the only third party library. All other functionality is provided by Go's standard library.

## Modules & code
This project is structured as follows:
- `github.com/deeprave/go-crm` is the project root. It contains the following two submodules:
  - `crm` contains the customer "database"
  - `api` contains the api including handlers

All files have high test coverage in the provided *_test.go files and may be run using:
```bash
go test ./...
```
The file `httptests.http` contains a sample list of requests used to develop the api.
This file may be used directly by any Jetbrains IDE to run the requests, which requires
that the server itself is running.

`api/crm_test.go` contains httptest tests that interact with and test the handlers.

### Udacity http test

`api/main_test.go` is the Udacity unit test, placed in the module in which the handlers are defined.

> **An important change was made to this test unit. Phone numbers in Australia
> start with the digit "0" which is lost if phone number is a numeric type.
> This API therefore uses a string for the Phone field in the Customer struct so the phone
> number provided in the addCustomer request needed to be quoted.**

### Additional paths

The `main.go` module defines two additional url paths:
- `GET /` displays `public/index.html`, which contains the present content
- `POST /load-test-data` allows loading some suitable test data in json format into the server.
In the body of the request, place the path to the file to load, e..g:
```json
{
  "path": "crm/data/customers.json"
}
```
