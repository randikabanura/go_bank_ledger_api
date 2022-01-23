package handler

import (
	"bank_ledger_api/app/handler/common"
	"bank_ledger_api/model"
	"encoding/json"
	"fmt"
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

	common.RespondJSON(w, http.StatusOK, customers)
}

func LoginCustomer(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("LoginCustomer")
	jsonCustomer := model.Customer{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&jsonCustomer); err != nil {
		common.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	email := jsonCustomer.Email

	customer := getCustomerByEmailOr404(db, email, w, r)
	if customer == nil {
		common.RespondError(w, http.StatusBadRequest, "Password not correct")
		return
	}

	isVerified := model.VerifyPassword(customer.Password, jsonCustomer.Password)

	if isVerified == nil {
		token, _ := model.CreateToken(customer.ID)
		common.RespondJSON(w, http.StatusCreated, map[string]string{"token": token})
		return
	} else {
		common.RespondError(w, http.StatusBadRequest, "Password not correct")
		return
	}
}

func CreateCustomer(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("CreateCustomer")
	customer := model.Customer{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&customer); err != nil {
		common.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&customer).Error; err != nil {
		common.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	token, _ := model.CreateToken(customer.ID)

	common.RespondJSON(w, http.StatusCreated, map[string]string{"token": token})
}

func GetCustomer(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("GetCustomer")
	vars := mux.Vars(r)
	id := vars["id"]
	customer := getCustomerOr404(db, id, w, r)
	if customer == nil {
		return
	}

	common.RespondJSON(w, http.StatusOK, customer)
}

func GetTransactionsByCustomer(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("GetTransactionsByCustomer")
	customerId, err := model.ExtractTokenID(r)
	if err != nil {
		common.RespondJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	customer := getCustomerOr404(db, fmt.Sprint(customerId), w, r)
	if customer == nil {
		return
	}

	var transactions []model.Transaction
	db.Model(&customer).Association("Transactions").Find(&transactions)

	for i, _ := range transactions {
		db.Model(transactions[i]).Related(&transactions[i].Account)
		db.Model(transactions[i]).Related(&transactions[i].ToAccount)
	}

	common.RespondJSON(w, http.StatusOK, transactions)
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
	common.RespondJSON(w, http.StatusOK, accounts)
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
		common.RespondJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	if err := db.Save(&customer).Error; err != nil {
		common.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.RespondJSON(w, http.StatusOK, customer)
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
		common.RespondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.RespondJSON(w, http.StatusNoContent, nil)
}

//gets customer instance if exists or responds with 404 otherwise
func getCustomerOr404(db *gorm.DB, id string, w http.ResponseWriter, r *http.Request) *model.Customer {
	customer := model.Customer{}
	uid, _ := uuid.FromString(id)
	if err := db.First(&customer, model.Customer{ID: uid}).Error; err != nil {
		common.RespondError(w, http.StatusNotFound, err.Error())
		return nil
	}

	db.Model(customer).Related(&customer.Accounts)

	return &customer
}

//gets customer instance if exists or responds with 404 otherwise
func getCustomerByEmailOr404(db *gorm.DB, email string, w http.ResponseWriter, r *http.Request) *model.Customer {
	customer := model.Customer{}
	if err := db.First(&customer, model.Customer{Email: email}).Error; err != nil {
		common.RespondError(w, http.StatusNotFound, err.Error())
		return nil
	}

	db.Model(customer).Related(&customer.Accounts)

	return &customer
}
