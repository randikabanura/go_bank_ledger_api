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

func GetAllAccounts(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("GetAllAccounts")
	var accounts []model.Account
	db.Find(&accounts)
	respondJSON(w, http.StatusOK, accounts)
}

func CreateAccount(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("CreateAccount")
	account := model.Account{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&account); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&account).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, account)
}

func GetAccount(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	account := getAccountOr404(db, id, w, r)
	if account == nil {
		return
	}

	log.Println("Count: ", db.Model(&account).Association("Customer").Count())
	respondJSON(w, http.StatusOK, account)
}

func GetCustomerByAccount(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var customer model.Customer
	account := getAccountOr404(db, id, w, r)
	if account == nil {
		return
	}

	db.Model(&account).Association("Customer").Find(&customer)

	respondJSON(w, http.StatusOK, customer)
}

func UpdateAccount(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	account := getAccountOr404(db, id, w, r)
	if account == nil {
		return
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&account); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	if err := db.Save(&account).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, account)
}

func DeleteAccount(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	account := getAccountOr404(db, id, w, r)
	if account == nil {
		return
	}
	if err := db.Delete(&account).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

//gets employee instance if exists or responds with 404 otherwise
func getAccountOr404(db *gorm.DB, id string, w http.ResponseWriter, r *http.Request) *model.Account {
	log.Println("Test")
	account := model.Account{}
	uid, _ := uuid.FromString(id)
	if err := db.First(&account, model.Account{ID: uid}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}

	return &account
}
