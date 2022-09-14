package helpers

import (
	"context"
	"log"
	"time"
	"todolist/configs"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//SignedDetails
type SignedDetails struct { //struct to store user_id and token expiration time
	User_id            string `json:"user_id"`
	jwt.StandardClaims        //struct to store token expiration time
}

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users") //get user collection from database

var SECRET_KEY = configs.EnvMongoURI("SECRET_KEY") //get secret key from environment variable

//GenerateAllToken generates both detailed token and refreshtoken
func GenerateAllTokens(uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{ //create claims for token with user_id and token expiration time
		User_id: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 24).Unix(), //token expires in 24 hours
		},
	}

	refreshCliams := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 24 * 7).Unix(), //refresh token expires in 7 days
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY)) //sign token with secret key and return signed token and error if any error occurs
	if err != nil {
		log.Fatal("Error while signing token:", err)
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshCliams).SignedString([]byte(SECRET_KEY)) //sign refresh token with secret key and return signed refresh token and error if any occurs

	if err != nil {
		log.Fatal("Error while signing refresh token:", err)
		return "", "", err
	}

	return token, refreshToken, nil
}

//ValidateToken validates the jwt token

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims( //parse token with claims and return claims and error if any error occurs
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) { //function to validate token with secret key
			return []byte(SECRET_KEY), nil //return secret key and nil error
		},
	)

	if err != nil {
		log.Fatal("Error while parsing token:", err)
		return nil, err.Error()
	}

	claims, ok := token.Claims.(*SignedDetails) //get claims from token
	if !ok {
		log.Fatal("Token claims are not valid")
		return nil, "Token claims are not valid"
	}

	if claims.ExpiresAt < time.Now().Local().Unix() { //check if token is expired
		log.Fatal("Token is expired")
		return nil, "Token is expired"
	}
	log.Fatal("Token Parse jwt:", token)

	return claims, ""
}

//UpdateToken updates the jwt token renews the token when they login
func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var updatedObj primitive.D //create a primitive.D object to update the token and refresh token in database

	updatedObj = append(updatedObj, primitive.E{Key: "token", Value: signedToken})
	updatedObj = append(updatedObj, primitive.E{Key: "refreshToken", Value: signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updatedObj = append(updatedObj, primitive.E{Key: "updated_at", Value: Updated_at})

	upsert := true //if user does not exist then create new user

	filter := bson.M{"user_id": userId} //filter to find the user with user_id
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(ctx, filter, bson.D{{"$set", updatedObj}}, &opt) //update token and refresh token in database and return error if any error occurs

	defer cancel()

	if err != nil {
		log.Fatal("Error while updating token:", err)
		return
	}
}
