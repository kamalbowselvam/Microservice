package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kamalselvam/Microservice/data"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "unable to marshal the json", http.StatusInternalServerError)
		return
	}

}

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {

		http.Error(rw, "unable to convert id", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT Products")

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	p.l.Printf("%#v", prod)
	err = data.UpdateProduct(id, &prod)
	if err == data.ErrorProductNotFound {
		http.Error(rw, "Product Not Found ", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product Not Found ", http.StatusInternalServerError)
		return
	}
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
		err =prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating the product ", err)
			http.Error(rw,fmt.Sprintf("Error Validating product: %s",err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})

}
