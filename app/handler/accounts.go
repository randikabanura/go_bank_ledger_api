package handler

import (
	"bank_ledger_api/app/handler/common"
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
	for i, _ := range accounts {
		db.Model(accounts[i]).Related(&accounts[i].Customer)
	}
	common.RespondJSON(w, http.StatusOK, accounts)
}

func CreateAccount(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("CreateAccount")
	account := model.Account{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&account); err != nil {
		common.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	account.CustomerID, _ = model.ExtractTokenID(r)

	if err := db.Save(&account).Error; err != nil {
		common.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.RespondJSON(w, http.StatusCreated, account)
}

func GetAccount(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("GetAccount")
	vars := mux.Vars(r)
	id := vars["id"]
	account := getAccountOr404(db, id, w, r)
	if account == nil {
		return
	}

	common.RespondJSON(w, http.StatusOK, account)
}

func GetTransactionsByAccount(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("GetTransactionsByAccount")
	vars := mux.Vars(r)
	id := vars["id"]
	account := getAccountOr404(db, id, w, r)
	if account == nil {
		return
	}

	var transactions []model.Transaction
	var toTransactions []model.Transaction
	db.Model(account).Association("Transactions").Find(&transactions)
	db.Model(account).Association("ToTransactions").Find(&toTransactions)

	for i, _ := range transactions {
		db.Model(transactions[i]).Related(&transactions[i].Account)
		db.Model(transactions[i]).Related(&transactions[i].ToAccount)
	}

	for j, _ := range toTransactions {
		db.Model(toTransactions[j]).Related(&toTransactions[j].ToAccount)
		db.Model(toTransactions[j]).Related(&toTransactions[j].Account)
	}

	common.RespondJSON(w, http.StatusOK, append(transactions, toTransactions...))
}

func GetCustomerByAccount(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("GetCustomerByAccount")
	vars := mux.Vars(r)
	id := vars["id"]
	var customer model.Customer
	account := getAccountOr404(db, id, w, r)
	if account == nil {
		return
	}

	db.Model(&account).Association("Customer").Find(&customer)

	common.RespondJSON(w, http.StatusOK, customer)
}

func UpdateAccount(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("UpdateAccount")
	vars := mux.Vars(r)
	id := vars["id"]
	account := getAccountOr404(db, id, w, r)
	if account == nil {
		return
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&account); err != nil {
		common.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	if err := db.Save(&account).Error; err != nil {
		common.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.RespondJSON(w, http.StatusOK, account)
}

func DeleteAccount(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("DeleteAccount")
	vars := mux.Vars(r)
	id := vars["id"]
	account := getAccountOr404(db, id, w, r)
	if account == nil {
		return
	}
	if err := db.Delete(&account).Error; err != nil {
		common.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.RespondJSON(w, http.StatusNoContent, nil)
}

//gets account instance if exists or responds with 404 otherwise
func getAccountOr404(db *gorm.DB, id string, w http.ResponseWriter, r *http.Request) *model.Account {
	account := model.Account{}
	uid, _ := uuid.FromString(id)
	if err := db.First(&account, model.Account{ID: uid}).Error; err != nil {
		common.RespondError(w, http.StatusNotFound, err.Error())
		return nil
	}

	db.Model(account).Related(&account.Customer)

	return &account
}
