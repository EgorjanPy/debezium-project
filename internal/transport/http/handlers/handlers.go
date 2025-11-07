package handlers

import (
	"context"
	"debez/internal/models"
	"debez/internal/transport/http/modelsDTO"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type UserService interface {
	GetUsers(ctx context.Context, offset, limit int) ([]models.User, error)
	GetUserByID(ctx context.Context, id int64) (models.User, error)
	SaveUser(ctx context.Context, user models.User) error
	UpdateUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, id int64) error
}

type HandlerFacade struct {
	ctx     context.Context
	service UserService
}

func NewHandlerFacade(ctx context.Context, service UserService) *HandlerFacade {
	return &HandlerFacade{
		ctx:     ctx,
		service: service,
	}
}

func (h *HandlerFacade) GetUsers(w http.ResponseWriter, r *http.Request) {
	offset := 0
	limit := 10
	users, err := h.service.GetUsers(h.ctx, offset, limit)
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}
}
func (h *HandlerFacade) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Path[len("/api/v1/user/")+1:])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	user, err := h.service.GetUserByID(h.ctx, int64(userID))

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Failed to get user", http.StatusNotFound)
		return
	}
	userDTO := modelsDTO.UserDTO{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		LastName: user.LastName,
		Role:     user.Role,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userDTO); err != nil {
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}
}
func (h *HandlerFacade) SaveUser(w http.ResponseWriter, r *http.Request) {
	var user modelsDTO.CreateUserDTO
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
		return
	}
	newUser := models.User{
		Name:     user.Name,
		Email:    user.Email,
		LastName: user.LastName,
		Role:     user.Role,
	}
	err = h.service.SaveUser(h.ctx, newUser)
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
func (h *HandlerFacade) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user modelsDTO.UpdateUserDTO
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
		return
	}
	updatedUser := models.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		LastName: user.LastName,
		Role:     user.Role,
	}
	err = h.service.UpdateUser(h.ctx, updatedUser)

	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (h *HandlerFacade) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Path[len("/api/v1/delete_user/"):])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteUser(h.ctx, int64(userID))
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
