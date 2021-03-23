package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/lorezi/boolang/inits"
	"github.com/lorezi/boolang/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PermissionController struct
type PermissionController struct {
}

func init() {
	m = inits.NewDB().MongoConn()
}

// NewPermissionController instance
func NewPermissionController() *PermissionController {
	return &PermissionController{}
}

func (pc PermissionController) GetPermissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// bussiness logic

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("empty")
}

func (pc PermissionController) CreatePermission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// bussiness logic
	p := models.PermissionGroup{}

	json.NewDecoder(r.Body).Decode(&p)

	c := m.Database("boolang").Collection("permissions")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	p.GroupID = primitive.NewObjectID().Hex()

	_, err := c.InsertOne(ctx, p)
	if err != nil {
		r := models.Result{
			Status:  "fail",
			Message: "Server error...",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(r)
		return
	}

	msg := models.Result{
		Status:  "success",
		Message: "permission created successfully",
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)

}
