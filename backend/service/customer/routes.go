package customer

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/nicolaics/pharmacon/types"
	"github.com/nicolaics/pharmacon/utils"
)

type Handler struct {
	custStore types.CustomerStore
	userStore types.UserStore
}

func NewHandler(custStore types.CustomerStore, userStore types.UserStore) *Handler {
	return &Handler{custStore: custStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/customer", h.handleRegister).Methods(http.MethodPost)
	router.HandleFunc("/customer/{val}", h.handleGetAll).Methods(http.MethodGet)
	router.HandleFunc("/customer/detail", h.handleGetOne).Methods(http.MethodPost)
	router.HandleFunc("/customer", h.handleDelete).Methods(http.MethodDelete)
	router.HandleFunc("/customer", h.handleModify).Methods(http.MethodPatch)

	router.HandleFunc("/customer", func(w http.ResponseWriter, r *http.Request) { utils.WriteJSONForOptions(w, http.StatusOK, nil) }).Methods(http.MethodOptions)
	router.HandleFunc("/customer/{val}", func(w http.ResponseWriter, r *http.Request) { utils.WriteJSONForOptions(w, http.StatusOK, nil) }).Methods(http.MethodOptions)
	router.HandleFunc("/customer/detail", func(w http.ResponseWriter, r *http.Request) { utils.WriteJSONForOptions(w, http.StatusOK, nil) }).Methods(http.MethodOptions)
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get JSON Payload
	var payload types.RegisterCustomerPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err, nil)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors), nil)
		return
	}

	// validate token
	user, err := h.userStore.ValidateUserToken(w, r, false)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user token invalid: %v", err), nil)
		return
	}

	// check if the customer exists
	_, err = h.custStore.GetCustomerByName(payload.Name)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest,
			fmt.Errorf("customer with name %s already exists", payload.Name), nil)
		return
	}

	err = h.custStore.CreateCustomer(types.Customer{
		Name: payload.Name,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, fmt.Sprintf("customer %s successfully created by %s", payload.Name, user.Name), nil)
}

func (h *Handler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	// validate token
	_, err := h.userStore.ValidateUserToken(w, r, false)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user token invalid: %v", err), nil)
		return
	}

	vars := mux.Vars(r)
	val := vars["val"]

	log.Println("customer val: ", val)

	var customers []types.Customer

	if val == "all" {
		customers, err = h.custStore.GetAllCustomers()
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err, nil)
			return
		}
	} else {
		id, err := strconv.Atoi(val)
		if err != nil {
			customers, err = h.custStore.GetCustomersBySearchName(val)
			if err != nil {
				utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("customer %s not found", val), nil)
				return
			}
		} else {
			customer, err := h.custStore.GetCustomerByID(id)
			if err != nil {
				utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("customer id %d not found", id), nil)
				return
			}

			customers = append(customers, *customer)
		}
	}

	utils.WriteSuccess(w, http.StatusOK, customers, nil)
}

func (h *Handler) handleGetOne(w http.ResponseWriter, r *http.Request) {
	// get JSON Payload
	var payload types.GetOneCustomerPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err, nil)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors), nil)
		return
	}

	// validate token
	_, err := h.userStore.ValidateUserToken(w, r, false)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user token invalid: %v", err), nil)
		return
	}

	// get customer data
	customer, err := h.custStore.GetCustomerByID(payload.ID)
	if customer == nil || err != nil {
		utils.WriteError(w, http.StatusBadRequest,
			fmt.Errorf("customer id %d doesn't exist", payload.ID), nil)
		return
	}

	utils.WriteSuccess(w, http.StatusOK, customer, nil)
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	// get JSON Payload
	var payload types.DeleteCustomerPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err, nil)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors), nil)
		return
	}

	// validate token
	user, err := h.userStore.ValidateUserToken(w, r, false)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user token invalid: %v", err), nil)
		return
	}

	// check if the customer exists
	customer, err := h.custStore.GetCustomerByID(payload.ID)
	if customer == nil || err != nil {
		utils.WriteError(w, http.StatusBadRequest,
			fmt.Errorf("customer %s doesn't exist", payload.Name), nil)
		return
	}

	err = h.custStore.DeleteCustomer(user, customer)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	utils.WriteSuccess(w, http.StatusOK, fmt.Sprintf("customer %s deleted by %s", payload.Name, user.Name), nil)
}

func (h *Handler) handleModify(w http.ResponseWriter, r *http.Request) {
	// get JSON Payload
	var payload types.ModifyCustomerPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err, nil)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors), nil)
		return
	}

	// validate token
	user, err := h.userStore.ValidateUserToken(w, r, false)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user token invalid: %v", err), nil)
		return
	}

	// check if the customer exists
	customer, err := h.custStore.GetCustomerByID(payload.ID)
	if err != nil || customer == nil {
		utils.WriteError(w, http.StatusBadRequest,
			fmt.Errorf("customer with id %d doesn't exists", payload.ID), nil)
		return
	}

	_, err = h.custStore.GetCustomerByName(payload.NewData.Name)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest,
			fmt.Errorf("customer with name %s already exist", payload.NewData.Name), nil)
		return
	}

	err = h.custStore.ModifyCustomer(customer.ID, payload.NewData.Name, user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, fmt.Sprintf("customer modified into %s by %s",
		payload.NewData.Name, user.Name), nil)
}
