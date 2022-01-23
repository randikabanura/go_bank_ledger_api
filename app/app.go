package app

import (
	"bank_ledger_api/app/handler"
	"bank_ledger_api/app/handler/common"
	"bank_ledger_api/config"
	"bank_ledger_api/model"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialize(config *config.Config) {
	dbURI := fmt.Sprintf("%s:%s@/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Name,
		config.DB.Charset)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	if err != nil {
		log.Fatal("Could not connect database", err)
	}
	a.DB = model.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.setRoutes()
}

//set required routes
func (a *App) setRoutes() {
	a.Post("/api/v1/auth/sign_in", a.LoginCustomer)
	a.Get("/api/v1/customers", a.GetAllCustomers)
	a.Post("/api/v1/auth", a.CreateCustomer)
	a.Get("/api/v1/customers/{id}", a.GetCustomer)
	a.Get("/api/v1/customers/{id}/accounts", a.GetAccountsByCustomer)
	a.Put("/api/v1/customers/{id}", a.UpdateCustomer)
	a.Delete("/api/v1/customers/{id}", a.DeleteCustomer)
	a.Get("/api/v1/accounts", a.GetAllAccounts)
	a.Post("/api/v1/accounts", a.CreateAccount)
	a.Get("/api/v1/accounts/{id}", a.GetAccount)
	a.Get("/api/v1/accounts/{id}/customer", a.GetCustomerByAccount)
	a.Put("/api/v1/accounts/{id}", a.UpdateAccount)
	a.Delete("/api/v1/accounts/{id}", a.DeleteAccount)
	a.Post("/api/v1/transactions", a.CreateTransaction)
	a.Get("/api/v1/transactions/{id}", a.GetTransaction)
	a.Get("/api/v1/transactions", a.GetTransactionsByCustomer)
	a.Get("/api/v1/accounts/{id}/transactions", a.GetTransactionsByAccount)
}

func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

func (a *App) LoginCustomer(w http.ResponseWriter, r *http.Request) {
	handler.LoginCustomer(a.DB, w, r)
}

func (a *App) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	err := SetMiddlewareAuthentication(r)
	if err == nil {
		handler.GetAllCustomers(a.DB, w, r)
	} else {
		common.RespondError(w, http.StatusUnauthorized, "Unauthorized")
	}
}

func (a *App) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	handler.CreateCustomer(a.DB, w, r)
}

func (a *App) GetCustomer(w http.ResponseWriter, r *http.Request) {
	err := SetMiddlewareAuthentication(r)
	if err == nil {
		handler.GetCustomer(a.DB, w, r)
	} else {
		common.RespondError(w, http.StatusUnauthorized, "Unauthorized")
	}
}

func (a *App) GetCustomerByAccount(w http.ResponseWriter, r *http.Request) {
	err := SetMiddlewareAuthentication(r)
	if err == nil {
		handler.GetCustomerByAccount(a.DB, w, r)
	} else {
		common.RespondError(w, http.StatusUnauthorized, "Unauthorized")
	}
}

func (a *App) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	err := SetMiddlewareAuthentication(r)
	if err == nil {
		handler.UpdateCustomer(a.DB, w, r)
	} else {
		common.RespondError(w, http.StatusUnauthorized, "Unauthorized")
	}
}

func (a *App) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	err := SetMiddlewareAuthentication(r)
	if err == nil {
		handler.DeleteCustomer(a.DB, w, r)
	} else {
		common.RespondError(w, http.StatusUnauthorized, "Unauthorized")
	}
}

func (a *App) GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	err := SetMiddlewareAuthentication(r)
	if err == nil {
		handler.GetAllAccounts(a.DB, w, r)
	} else {
		common.RespondError(w, http.StatusUnauthorized, "Unauthorized")
	}
}

func (a *App) CreateAccount(w http.ResponseWriter, r *http.Request) {
	err := SetMiddlewareAuthentication(r)
	if err == nil {
		handler.CreateAccount(a.DB, w, r)
	} else {
		common.RespondError(w, http.StatusUnauthorized, "Unauthorized")
	}
}

func (a *App) GetAccount(w http.ResponseWriter, r *http.Request) {
	err := SetMiddlewareAuthentication(r)
	if err == nil {
		handler.GetAccount(a.DB, w, r)
	} else {
		common.RespondError(w, http.StatusUnauthorized, "Unauthorized")
	}
}

func (a *App) GetAccountsByCustomer(w http.ResponseWriter, r *http.Request) {
	err := SetMiddlewareAuthentication(r)
	if err == nil {
		handler.GetAccountsByCustomer(a.DB, w, r)
	} else {
		common.RespondError(w, http.StatusUnauthorized, "Unauthorized")
	}
}

func (a *App) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	err := SetMiddlewareAuthentication(r)
	if err == nil {
		handler.UpdateAccount(a.DB, w, r)
	} else {
		common.RespondError(w, http.StatusUnauthorized, "Unauthorized")
	}
}

func (a *App) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	err := SetMiddlewareAuthentication(r)
	if err == nil {
		handler.DeleteAccount(a.DB, w, r)
	} else {
		common.RespondError(w, http.StatusUnauthorized, "Unauthorized")
	}
}

func (a *App) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	err := SetMiddlewareAuthentication(r)
	if err == nil {
		handler.CreateTransaction(a.DB, w, r)
	} else {
		common.RespondError(w, http.StatusUnauthorized, "Unauthorized")
	}
}

func (a *App) GetTransaction(w http.ResponseWriter, r *http.Request) {
	err := SetMiddlewareAuthentication(r)
	if err == nil {
		handler.GetTransaction(a.DB, w, r)
	} else {
		common.RespondError(w, http.StatusUnauthorized, "Unauthorized")
	}
}

func (a *App) GetTransactionsByCustomer(w http.ResponseWriter, r *http.Request) {
	err := SetMiddlewareAuthentication(r)
	if err == nil {
		handler.GetTransactionsByCustomer(a.DB, w, r)
	} else {
		common.RespondError(w, http.StatusUnauthorized, "Unauthorized")
	}
}

func (a *App) GetTransactionsByAccount(w http.ResponseWriter, r *http.Request) {
	err := SetMiddlewareAuthentication(r)
	if err == nil {
		handler.GetTransactionsByAccount(a.DB, w, r)
	} else {
		common.RespondError(w, http.StatusUnauthorized, "Unauthorized")
	}
}

func SetMiddlewareAuthentication(r *http.Request) error {
	err := model.TokenValid(r)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
