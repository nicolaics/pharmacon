package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/nicolaics/pharmacon/service/auth"
	"github.com/nicolaics/pharmacon/types"
	"github.com/nicolaics/pharmacon/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/user/register", h.handleRegister).Methods(http.MethodPost)
	router.HandleFunc("/user/register", func(w http.ResponseWriter, r *http.Request) { utils.WriteJSONForOptions(w, http.StatusOK, nil) }).Methods(http.MethodOptions)

	router.HandleFunc("/user/{params}/{val}", h.handleGetAll).Methods(http.MethodGet)
	router.HandleFunc("/user/{params}/{val}", func(w http.ResponseWriter, r *http.Request) { utils.WriteJSONForOptions(w, http.StatusOK, nil) }).Methods(http.MethodOptions)

	router.HandleFunc("/user/current", h.handleGetCurrentUser).Methods(http.MethodGet)
	router.HandleFunc("/user/current", func(w http.ResponseWriter, r *http.Request) { utils.WriteJSONForOptions(w, http.StatusOK, nil) }).Methods(http.MethodOptions)

	router.HandleFunc("/user/detail", h.handleGetOneUser).Methods(http.MethodPost)
	router.HandleFunc("/user/detail", func(w http.ResponseWriter, r *http.Request) { utils.WriteJSONForOptions(w, http.StatusOK, nil) }).Methods(http.MethodOptions)

	router.HandleFunc("/user", h.handleDelete).Methods(http.MethodDelete)
	router.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) { utils.WriteJSONForOptions(w, http.StatusOK, nil) }).Methods(http.MethodOptions)

	router.HandleFunc("/user/modify", h.handleModify).Methods(http.MethodPatch)
	router.HandleFunc("/user/modify", func(w http.ResponseWriter, r *http.Request) { utils.WriteJSONForOptions(w, http.StatusOK, nil) }).Methods(http.MethodOptions)

	router.HandleFunc("/user/logout", h.handleLogout).Methods(http.MethodGet)
	router.HandleFunc("/user/logout", func(w http.ResponseWriter, r *http.Request) { utils.WriteJSONForOptions(w, http.StatusOK, nil) }).Methods(http.MethodOptions)

	router.HandleFunc("/user/admin", h.handleChangeAdminStatus).Methods(http.MethodPatch)
	router.HandleFunc("/user/admin", func(w http.ResponseWriter, r *http.Request) { utils.WriteJSONForOptions(w, http.StatusOK, nil) }).Methods(http.MethodOptions)
}

func (h *Handler) RegisterUnprotectedRoutes(router *mux.Router) {
	router.HandleFunc("/user/login", h.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/user/login", func(w http.ResponseWriter, r *http.Request) { utils.WriteJSONForOptions(w, http.StatusOK, nil) }).Methods(http.MethodOptions)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// get JSON Payload
	var payload types.LoginUserPayload

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

	user, err := h.store.GetUserByName(payload.Name)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid name: %v", err), nil)
		return
	}

	// check password match
	if !(auth.ComparePassword(user.Password, []byte(payload.Password))) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("not found, invalid password"), nil)
		return
	}

	tokenDetails, err := auth.CreateJWT(user.ID, user.Admin)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	err = h.store.SaveToken(user.ID, tokenDetails)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	err = h.store.UpdateLastLoggedIn(user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	tokens := map[string]string{
		"token": tokenDetails.Token,
	}

	utils.WriteSuccess(w, http.StatusOK, tokens, nil)
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get JSON Payload
	var payload types.RegisterUserPayload

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
	admin, err := h.store.ValidateUserToken(w, r, true)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid admin token or not admin: %v", err), nil)
		return
	}

	// validate admin password
	if !(auth.ComparePassword(admin.Password, []byte(payload.AdminPassword))) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("admin password wrong"), nil)
		return
	}

	// check if the newly created user exists
	_, err = h.store.GetUserByName(payload.Name)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest,
			fmt.Errorf("user with name %s already exists", payload.Name), nil)
		return
	}

	// if it doesn't, we create new user
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	err = h.store.CreateUser(types.User{
		Name:        payload.Name,
		Password:    hashedPassword,
		PhoneNumber: payload.PhoneNumber,
		Admin:       payload.Admin,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
	}

	utils.WriteSuccess(w, http.StatusCreated, fmt.Sprintf("user %s successfully created", payload.Name), nil)
}

func (h *Handler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	// validate token
	_, err := h.store.ValidateUserToken(w, r, true)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid admin token or not admin: %v", err), nil)
		return
	}

	vars := mux.Vars(r)
	params := vars["params"]
	val := vars["val"]

	var users []types.User

	if val == "all" {
		users, err = h.store.GetAllUsers()
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err, nil)
			return
		}
	} else if params == "name" {
		users, err = h.store.GetUserBySearchName(val)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user %s not found", val), nil)
			return
		}
	} else if params == "id" {
		id, err := strconv.Atoi(val)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err, nil)
			return
		}

		user, err := h.store.GetUserByID(id)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user id %s not found", val), nil)
			return
		}

		users = append(users, *user)
	} else if params == "phone-number" {
		users, err = h.store.GetUserBySearchPhoneNumber(val)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user %s not found", val), nil)
			return
		}
	} else {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("unknown query"), nil)
		return
	}

	utils.WriteSuccess(w, http.StatusOK, users, nil)
}

func (h *Handler) handleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// validate token
	user, err := h.store.ValidateUserToken(w, r, true)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid admin token or not admin: %v", err), nil)
		return
	}

	utils.WriteSuccess(w, http.StatusOK, user, nil)
}

// get one user other than the current user
func (h *Handler) handleGetOneUser(w http.ResponseWriter, r *http.Request) {
	// get JSON Payload
	var payload types.GetOneUserPayload

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
	_, err := h.store.ValidateUserToken(w, r, true)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid admin token or not admin: %v", err), nil)
		return
	}

	// check if the newly created user exists
	user, err := h.store.GetUserByID(payload.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest,
			fmt.Errorf("user id %d doesn't exist", payload.ID), nil)
		return
	}

	utils.WriteSuccess(w, http.StatusOK, user, nil)
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	var payload types.RemoveUserPayload

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
	admin, err := h.store.ValidateUserToken(w, r, true)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid admin token or not admin: %v", err), nil)
		return
	}

	// validate admin password
	if !(auth.ComparePassword(admin.Password, []byte(payload.AdminPassword))) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("admin password wrong"), nil)
		return
	}

	users, err := h.store.GetAllUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	if len(users) == 1 || payload.ID == 1 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("cannot delete initial admin"), nil)
		return
	}

	user, err := h.store.GetUserByID(payload.ID)
	if user == nil || err != nil {
		utils.WriteError(w, http.StatusBadRequest, err, nil)
		return
	}

	err = h.store.DeleteUser(user, admin)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	utils.WriteSuccess(w, http.StatusOK, fmt.Sprintf("%s successfully deleted", user.Name), nil)
}

func (h *Handler) handleModify(w http.ResponseWriter, r *http.Request) {
	var payload types.ModifyUserPayload

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
	admin, err := h.store.ValidateUserToken(w, r, true)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid admin token or not admin: %v", err), nil)
		return
	}

	// validate admin password
	if !(auth.ComparePassword(admin.Password, []byte(payload.NewData.AdminPassword))) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("admin password wrong"), nil)
		return
	}

	user, err := h.store.GetUserByID(payload.ID)
	if user == nil {
		utils.WriteError(w, http.StatusBadRequest, err, nil)
		return
	}

	if user.Name != payload.NewData.Name {
		_, err = h.store.GetUserByName(payload.NewData.Name)
		if err == nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with name %s already exists", payload.NewData.Name), nil)
		}
	}

	err = h.store.ModifyUser(user.ID, types.User{
		Name:        payload.NewData.Name,
		Password:    payload.NewData.Password,
		Admin:       payload.NewData.Admin,
		PhoneNumber: payload.NewData.PhoneNumber,
	}, admin)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	utils.WriteSuccess(w, http.StatusOK, fmt.Sprintf("%s updated into", payload.NewData.Name), nil)
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	accessDetails, err := auth.ExtractTokenFromClient(r)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token"), nil)
		return
	}

	// check user exists or not
	_, err = h.store.GetUserByID(accessDetails.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user id %d doesn't exists", accessDetails.UserID), nil)
		return
	}

	err = h.store.UpdateLastLoggedIn(accessDetails.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	err = h.store.DeleteToken(accessDetails.UserID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	utils.WriteSuccess(w, http.StatusOK, "successfully logged out", nil)
}

func (h *Handler) handleChangeAdminStatus(w http.ResponseWriter, r *http.Request) {
	var payload types.ChangeAdminStatusPayload

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
	admin, err := h.store.ValidateUserToken(w, r, true)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid admin token or not admin: %v", err), nil)
		return
	}

	// validate admin password
	if !(auth.ComparePassword(admin.Password, []byte(payload.AdminPassword))) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("admin password wrong"), nil)
		return
	}

	// check whether user exists or not
	user, err := h.store.GetUserByID(payload.ID)
	if user == nil {
		utils.WriteError(w, http.StatusBadRequest, err, nil)
		return
	}

	err = h.store.ModifyUser(user.ID, types.User{
		Name:        user.Name,
		Password:    user.Password,
		Admin:       payload.Admin,
		PhoneNumber: user.PhoneNumber,
	}, admin)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err, nil)
		return
	}

	utils.WriteSuccess(w, http.StatusOK, fmt.Sprintf("%s updated into admin: %t", user.Name, payload.Admin), nil)
}
