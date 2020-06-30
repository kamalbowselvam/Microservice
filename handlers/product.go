package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/kamalselvam/Microservice/data"
)
	
type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		p.addProduct(rw,r)
		return
	}

	if r.Method == http.MethodPut {
		
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)
		
		p.l.Println(g)
		if len(g) != 1 {
			p.l.Println("Invalid URI more than one id")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Println("Invalid URI more than one capture group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)

		if err != nil {
			p.l.Println("Invalid URI unable to convert to number", idString)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		p.updateProduct(id,rw,r)
		
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)
}


func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle POST Products")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "unable to unmarshal the json", http.StatusBadRequest)
		return
	}
	p.l.Printf("%#v",prod)
	data.AddProduct(prod)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request){
	p.l.Println("Handle GET Products")

	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "unable to marshal the json", http.StatusInternalServerError)
		return
	}

}


func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request){

	p.l.Println("Handle PUT Products")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "unable to unmarshal the json", http.StatusBadRequest)
		return
	}
	p.l.Printf("%#v",prod)
	err = data.UpdateProduct(id, prod)
	if err == data.ErrorProductNotFound {
		http.Error(rw, "Product Not Found ", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product Not Found ", http.StatusInternalServerError)
		return
	}
}