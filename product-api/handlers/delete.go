package handlers

import (
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/kamalbowselvam/Microservice/product-api/data"
)
// swagger:route DELETE /products/{id} products deleteProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  501: errorResponse

// Delete handles DELETE requests and removes items from the database
func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	p.l.Println("Handle Delete Product", id)

	err := data.DeleteProduct(id)

	if err == data.ErrorProductNotFound {

		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}

}
