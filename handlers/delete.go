package handlers

import (
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/kamalselvam/Microservice/data"
)
// swagger:route DELETE /products/{id} products deleteProduct
// Reurns a list of products
// responses:
// 	201: noContent

//DeleteProduct deletes a product from the data store
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