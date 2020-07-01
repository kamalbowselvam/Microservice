package handlers

import (
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/kamalselvam/Microservice/data"
)

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// Update handles PUT requests to update products
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