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

func CreateTransaction(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("CreateTransaction")
	transaction := model.Transaction{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&transaction); err != nil {
		common.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&transaction).Error; err != nil {
		common.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.RespondJSON(w, http.StatusCreated, transaction)
}

func GetTransaction(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Print("GetTransaction")
	vars := mux.Vars(r)
	id := vars["id"]
	transaction := getTransactionOr404(db, id, w, r)
	if transaction == nil {
		return
	}

	common.RespondJSON(w, http.StatusOK, transaction)
}

//gets transaction instance if exists or responds with 404 otherwise
func getTransactionOr404(db *gorm.DB, id string, w http.ResponseWriter, r *http.Request) *model.Transaction {
	transaction := model.Transaction{}
	uid, _ := uuid.FromString(id)
	if err := db.First(&transaction, model.Transaction{ID: uid}).Error; err != nil {
		common.RespondError(w, http.StatusNotFound, err.Error())
		return nil
	}

	return &transaction
}
