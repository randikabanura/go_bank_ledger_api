package handler

import (
	"bank_ledger_api/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
)

func GetAllCustomers(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("GetAllCustomers")
	var customers []model.Customer
	db.Find(&customers)

	for i, _ := range customers {
		db.Model(customers[i]).Related(&customers[i].Accounts)
	}

	respondJSON(w, http.StatusOK, customers)
}

func CreateCustomer(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("CreateAccount")
	customer := model.Customer{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&customer); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&customer).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, customer)
}

func GetCustomer(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("GetCustomer")
	vars := mux.Vars(r)
	id := vars["id"]
	customer := getCustomerOr404(db, id, w, r)
	if customer == nil {
		return
	}

	respondJSON(w, http.StatusOK, customer)
}

func GetAccountsByCustomer(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("GetAccountsByCustomer")
	vars := mux.Vars(r)
	id := vars["id"]
	var accounts []model.Account
	customer := getCustomerOr404(db, id, w, r)
	if customer == nil {
		return
	}
	db.Model(&customer).Association("Accounts").Find(&accounts)
	respondJSON(w, http.StatusOK, accounts)
}

func UpdateCustomer(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("UpdateCustomer")
	vars := mux.Vars(r)
	id := vars["id"]
	customer := getCustomerOr404(db, id, w, r)
	if customer == nil {
		return
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&customer); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	if err := db.Save(&customer).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, customer)
}

func DeleteCustomer(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("DeleteCustomer")
	vars := mux.Vars(r)
	id := vars["id"]
	customer := getCustomerOr404(db, id, w, r)
	if customer == nil {
		return
	}
	if err := db.Delete(&customer).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

//gets employee instance if exists or responds with 404 otherwise
func getCustomerOr404(db *gorm.DB, id string, w http.ResponseWriter, r *http.Request) *model.Customer {
	customer := model.Customer{}
	//var accounts []model.Account
	uid, _ := uuid.FromString(id)
	if err := db.First(&customer, model.Customer{ID: uid}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}

	db.Model(customer).Related(&customer.Accounts)

	return &customer
}
