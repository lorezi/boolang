package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	. "github.com/gobeam/mongo-go-pagination"
	"github.com/gorilla/mux"

	"github.com/lorezi/boolang/inits"
	"github.com/lorezi/boolang/models"
	"go.mongodb.org/mongo-driver/bson"
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
	ps := []models.PermissionGroup{}
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64) // Note: error not handled
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)   // Note: error not handled

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	c := m.Database("boolang").Collection("permissions")

	filter := bson.D{{}}

	res, err := New(c).Context(ctx).Limit(limit).Page(page).Filter(filter).Decode(&ps).Find()
	if err != nil {
		r := models.Result{
			Status:  "error",
			Message: "Internal server error",
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(r)
	}

	for _, v := range res.Data {
		var p *models.PermissionGroup

		if err := bson.Unmarshal(v, &p); err == nil {
			ps = append(ps, *p)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ps)
}

func GetPermission(path string) models.PermissionGroup {

	p := models.PermissionGroup{}

	filter := bson.M{"group_id": path}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	c := m.Database("boolang").Collection("permissions")

	err := c.FindOne(ctx, filter).Decode(&p)
	if err != nil {
		log.Panic(err)

	}

	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(p)
	return p
}

func (pc PermissionController) GetPermission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	p := models.PermissionGroup{}

	pv := mux.Vars(r)
	filter := bson.M{"group_id": pv["id"]}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	c := m.Database("boolang").Collection("permissions")

	err := c.FindOne(ctx, filter).Decode(&p)
	if err != nil {
		r := models.Result{
			Status:  "error",
			Message: "permission id does not exist",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(r)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
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
