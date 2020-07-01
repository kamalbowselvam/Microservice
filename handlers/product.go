// Package classification Petstore API.
//
// Documentation of API
//
//     Schemes: http
//     Host: localhost
//     BasePath: /
//     Version: 1.0.0
//     Contact: Kamal SELVAM<kselvam.phd@gmail.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"github.com/kamalselvam/Microservice/data"
	"log"
	"net/http"
)

type Products struct {
	l *log.Logger
}

// A list of products returns in the response
// swagger:response productsResponse
type proudctsResponseWrapper struct {
	// All products in the system
	// in: body
	Body []data.Product
}

// swagger:parameters deleteProduct
type productIDParameterWrapper struct {
	// The id of the product to delete from the database
	// in: path
	// required: true
	ID int `json:"id"`
}
// swagger:response noContent
type productsNoContentWrapper struct {

}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

type KeyProduct struct {
}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		prod := data.Product{}
		err := prod.FromJSON(r.Body)

		if err != nil {
			p.l.Println("[ERROR] deserializing product error ", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// validate the product
		err = prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating the product ", err)
			http.Error(rw, fmt.Sprintf("Error Validating product: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})

}
