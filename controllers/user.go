package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	. "github.com/gobeam/mongo-go-pagination"
	"github.com/gorilla/mux"

	"github.com/go-playground/validator"
	"github.com/lorezi/boolang/helpers"
	"github.com/lorezi/boolang/inits"
	"github.com/lorezi/boolang/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// var mu *mongo.Client = inits.NewDB().MongoConn()
var validate = validator.New()

// UserController Struct
type UserController struct {
}

// NewUserController instance
func NewUserController() *UserController {
	return &UserController{}
}

// hashPassword is used to encrypt the password before it is stored in the DB
func hashPassword(password string) string {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	inits.LogFatal(err)
	return string(bs)
}

// verifyPassword checks the input password while verifying it the password in the DB.
func verifyPassword(up string, pp string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(pp), []byte(up))
	check := true
	msg := ""
	if err != nil {
		msg = "login or password is incorrect"
		check = false
	}

	return check, msg
}

// GetUsers returns all users
func (uc UserController) GetUsers(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	su := []models.User{}

	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		http.Error(w, "limit query params missing...😵😵😵", http.StatusNotFound)
	}

	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		http.Error(w, "page query params missing...😵😵😵", http.StatusNotFound)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.Database("boolang").Collection("users")
	filter := bson.D{{}}

	res, err := New(collection).Context(ctx).Limit(limit).Page(page).Filter(filter).Decode(&su).Find()
	if err != nil {
		http.Error(w, "Server error 😵😵😵", http.StatusInternalServerError)
	}

	for _, v := range res.Data {
		var u *models.User

		if err := bson.Unmarshal(v, &u); err == nil {
			su = append(su, *u)
		}
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(su)

}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := models.User{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c := m.Database("boolang").Collection("users")

	p := mux.Vars(r)
	f := bson.M{"user_id": p["id"]}

	err := c.FindOne(ctx, f).Decode(&u)
	if err != nil {
		r := models.Result{
			Status:  "error",
			Message: "user does not exist",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(r)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)

}

// CreateUser is the api used to create a new user

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	u := models.User{}

	// map json request to u variable
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	err := d.Decode(&u)
	if err != nil {
		msg := helpers.JSONValidator(err)
		r := models.Result{
			Status:  "error",
			Message: msg,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(r)
		return
	}

	// validate the struct
	err = validate.Struct(u)
	if err != nil {
		// var msg string
		// for _, err := range err.(validator.ValidationErrors) {
		// 	msg += err.Field()
		// 	msg += " " + err.Tag()
		// 	msg += " " + err.Type().String() + ", "

		// }
		r := models.Result{
			Status:  "validation error",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.Database("boolang").Collection("users")

	pw := hashPassword(u.Password)
	u.Password = pw
	u.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	u.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	u.ID = primitive.NewObjectID()
	u.UserID = u.ID.Hex()

	tk, rtk, _ := helpers.GenerateAllTokens(u.Email, u.FirstName, u.LastName, u.UserID)

	u.Token = tk
	u.RefreshToken = rtk

	_, err = collection.InsertOne(ctx, u)
	if err != nil {
		r := models.Result{
			Status:  "fail",
			Message: "User account was not created 😰😰😰",
		}
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(r)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)

}

// Update user method
func (uc UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := models.User{}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	err := d.Decode(&u)
	if err != nil {
		msg := helpers.JSONValidator(err)
		r := models.Result{
			Status:  "error",
			Message: msg,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.Database("boolang").Collection("users")

	p := mux.Vars(r)
	filter := bson.M{"user_id": p["id"]}

	var uu primitive.D

	if u.FirstName != "" {
		uu = append(uu, bson.E{Key: "first_name", Value: u.FirstName})
	}

	if u.LastName != "" {
		uu = append(uu, bson.E{Key: "last_name", Value: u.LastName})
	}

	if u.Email != "" {
		uu = append(uu, bson.E{Key: "email", Value: u.Email})
	}

	if u.PhoneNo != "" {
		uu = append(uu, bson.E{Key: "phone_no", Value: u.PhoneNo})
	}

	if u.Address != "" {
		uu = append(uu, bson.E{Key: "address", Value: u.Address})
	}

	u.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	uu = append(uu, bson.E{Key: "updated_at", Value: u.UpdatedAt})

	upsert := true
	opts := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err = collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: uu}}, &opts)
	if err != nil {
		r := models.Result{
			Status:  "error",
			Message: "unable to update user",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(r)
		return
	}

	msg := models.Result{
		Status:  "success",
		Message: "user account updated successfully",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)

}

func (uc UserController) Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	// logged user
	logu := models.Login{}
	// user
	u := models.User{}

	// map json request to u variable
	err := json.NewDecoder(r.Body).Decode(&logu)
	if err != nil {
		r := models.Result{
			Status:  "error",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.Database("boolang").Collection("users")

	filter := bson.M{
		"email": logu.Email,
	}
	err = collection.FindOne(ctx, filter).Decode(&u)
	if err != nil {
		r := models.Result{
			Status:  "error",
			Message: "login or password incorrect",
		}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(r)
		return
	}
	fmt.Println("am here")

	ok, msg := verifyPassword(logu.Password, u.Password)
	if !ok {
		r := models.Result{
			Status:  "error",
			Message: msg,
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(r)
		return
	}

	tk, rtk, _ := helpers.GenerateAllTokens(u.Email, u.FirstName, u.LastName, u.UserID)
	helpers.UpdateAllTokens(tk, rtk, u.UserID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)

}
