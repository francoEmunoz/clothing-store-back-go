package handlers

import (
	"cs-go/db"
	"cs-go/models"
	"encoding/json"

	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetProducts(rw http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	name := queryParams.Get("name")
	category := queryParams.Get("category")
	orderBy := queryParams.Get("order_by")

	var products []models.Product
	db := db.Database.Where(&models.Product{})

	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if category != "" {
		db = db.Where("category = ?", category)
	}
	if orderBy == "price_asc" {
		db = db.Order("price ASC")
	} else if orderBy == "price_desc" {
		db = db.Order("price DESC")
	}

	db = db.Order("id DESC")

	db.Find(&products)
	sendData(rw, products, http.StatusOK)
}

func getProductById(r *http.Request) (models.Product, *gorm.DB) {

	vars := mux.Vars(r)
	productId, _ := strconv.Atoi(vars["id"])
	product := models.Product{}

	if err := db.Database.First(&product, productId); err.Error != nil {
		return product, err
	} else {
		return product, nil
	}
}

func GetProduct(rw http.ResponseWriter, r *http.Request) {

	if product, err := getProductById(r); err != nil {
		sendError(rw, http.StatusNotFound, "The product was not found")
	} else {
		sendData(rw, product, http.StatusOK)
	}

}

func CreateProduct(rw http.ResponseWriter, r *http.Request) {

	var validate = validator.New()

	product := models.Product{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&product); err != nil {
		sendError(rw, http.StatusUnprocessableEntity, "Decoding error")
		return
	}

	if err := validate.Struct(product); err != nil {
		sendValidationError(rw, err)
		return
	}

	db.Database.Save(&product)

	sendData(rw, product, http.StatusCreated)

}

func UpdateProduct(rw http.ResponseWriter, r *http.Request) {

	product, err := getProductById(r)
	if err != nil {
		sendError(rw, http.StatusNotFound, "The product was not found")
		return
	}

	var updatedProduct models.Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updatedProduct); err != nil {
		sendError(rw, http.StatusUnprocessableEntity, "Decoding error")
		return
	}

	db.Database.Model(&product).Omit("created_at").Updates(updatedProduct)

	sendData(rw, product, http.StatusOK)

}

func DeleteProduct(rw http.ResponseWriter, r *http.Request) {

	if product, err := getProductById(r); err != nil {
		sendError(rw, http.StatusNotFound, "The product was not found")
	} else {
		db.Database.Delete(&product)
		sendData(rw, product, http.StatusOK)
	}

}
