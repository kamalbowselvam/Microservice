package handlers

import (
	"github.com/kamalselvam/Microservice/data"
	"net/http"
)

// swagger:route GET /products products listProducts
// Reurns a list of products
// responses:
// 	200: productsResponse

//GetProducts returns the products from the data store
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "unable to marshal the json", http.StatusInternalServerError)
		return
	}

}
