// Helper package for miscellaneous functions for token
package helpers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lorezi/boolang/inits"
	"github.com/lorezi/boolang/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mt *mongo.Client

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func init() {
	mt = inits.NewDB().MongoConn()
}

func getPermission(path string) models.PermissionGroup {

	p := models.PermissionGroup{}

	filter := bson.M{"group_id": path}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	c := mt.Database("boolang").Collection("permissions")

	err := c.FindOne(ctx, filter).Decode(&p)
	if err != nil {

		log.Panic(err)

	}

	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(p)
	return p
}

//GenerateAllTokens generate both the detailed token and refresh token
func GenerateAllTokens(email string, firstName string, lastName string, uid string, permission string) (string, string, error) {

	res := getPermission(permission)

	claims := &models.SignedDetails{
		Email:       email,
		FirstName:   firstName,
		LastName:    lastName,
		UID:         uid,
		Permissions: res,
		StandardClaims: jwt.StandardClaims{
			// duration 1day
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &models.SignedDetails{
		Permissions: res,
		StandardClaims: jwt.StandardClaims{
			// duration 7days
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return "", "", err
	}

	return token, refreshToken, err
}

// ValidateToken validates the jwt token
func ValidateToken(signedToken string) (claims *models.SignedDetails, msg string) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&models.SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*models.SignedDetails)
	if !ok {
		msg = "the token is invalid"
		msg = err.Error()
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		msg = err.Error()
		return
	}

	return claims, msg
}

// UpdateAllTokens renew the user tokens when they login
func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {

	collection := mt.Database("boolang").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})
	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: updatedAt})

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := collection.UpdateOne(
		ctx, filter, bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)

	if err != nil {
		log.Panic(err)
	}

}
