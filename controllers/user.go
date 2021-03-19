package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/lorezi/boolang/helpers"
	"github.com/lorezi/boolang/inits"
	"github.com/lorezi/boolang/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// var mu *mongo.Client = inits.NewDB().MongoConn()
// var validate = validator.New()

// UserController Struct
type UserController struct {
}

// NewUserController instance
func NewUserController() *UserController {
	return &UserController{}
}

// HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	inits.LogFatal(err)
	return string(bs)
}

// VerifyPassword checks the input password while verifying it the password in the DB.
func VerifyPassword(up string, pp string) (bool, string) {
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

}

// CreateUser is the api used to create a new user

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	u := models.User{}

	// map json request to u variable
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		r := models.Result{
			Status:  "error",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(r)
		return
	}

	// validate the struct
	// err = validate.Struct(u)
	// if err != nil {
	// 	r := models.Result{
	// 		Status:  "error",
	// 		Message: err.Error(),
	// 	}
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(r)
	// 	return
	// }

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.Database("boolang").Collection("users")

	pw := HashPassword(u.Password)
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
			Message: "User item was not created 😰😰😰",
		}
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(r)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)

}
