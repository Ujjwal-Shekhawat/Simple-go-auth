package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Ujjwal-Shekhawat/golang-gin-poc/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var (
	collection *mongo.Collection
	client     *mongo.Client
)

// Init - Connect to db and ping it
func Init() {
	fmt.Println("Connecting to database")
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}

	// Creating context
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Connecting to my mongodb database
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database")

	collection = client.Database("users").Collection("users")
}

// Protect - A function that handels auth with JWT
func Protect() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Getting the cookie vals if exsists (Finding a better way)
		tkn, _ := ctx.Cookie("token")
		opt1 := verifyToken(ctx, tkn)
		if tkn == "" {
			opt1 = verifyToken(ctx)
		}
		if !(opt1) {
			ctx.AbortWithStatusJSON(401, gin.H{
				"message": "unauthorized (maybe the token expired or is invalid)",
			})
			return
		}

		// As I cannot use userschema struct here
		type loginForm struct {
			Name     string `json:"name,omitempty"`
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		var user loginForm
		// Unmarshalling the body json to struct and checking for errors
		// ShoudBindJSON shoud be used instead if BindJSON (this returns with Header text/json instead of text/plain)
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Server error",
			})
			return
		}

		// Retreving userdata
		var check loginForm
		filter := bson.D{{"email", user.Email}}
		singleresult := collection.FindOne(context.TODO(), filter)

		err := singleresult.Err()
		if err == mongo.ErrNoDocuments {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "No user found",
			})
		}

		err = singleresult.Decode(&check)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Server error",
			})
			return
		}
		// if check.Email == "" {
		// 	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		// 		"message": "invalid credentials",
		// 	})
		// 	return
		// }

		// Verify user password
		pass := func() bool {
			err := bcrypt.CompareHashAndPassword([]byte(check.Password), []byte(user.Password))
			if err != nil {
				return false
			}
			return true
		}()

		if !pass {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "invalid credentials",
			})
			return
		}

		// Generate and send the jwt back with response
		tok, err := createToken(user.Name)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Server error",
			})
			return
		}

		// Set cookies with token in them (expires in 15 minutes)
		ctx.SetCookie("token", tok, 15*60, "/", "localhost", false, true)

		check.Password = "nil"
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"message": "User found",
			"data":    check,
			"token":   tok,
		})

		// Prining the users name to console
		// fmt.Sprintf("%s Protect ran\n", name.Name)

		// This ctx.Next() func tion is a dead code do something about it (Its not like nodejs okay ujjwal)
		ctx.Next()
	}
}

func extractToken(r /* *http.Request */ *gin.Context) string {
	bearToken := r.Request.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func verifyToken(r /* *http.Request */ *gin.Context, value ...string) bool {
	tokenString := extractToken(r)

	if len(value) > 0 {
		tokenString = value[0]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECERET")), nil
	})
	if err != nil {
		return false
	}

	fmt.Println("Retrived token claim username : ", token.Claims.(jwt.MapClaims)["username"])
	return token.Valid
}

func createToken(username string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["username"] = username
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("JWT_SECERET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

// Regester - Docs incomplete
func Regester() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// testing var name of type model.UserModel struct
		var user model.UserModel
		// Unmarshalling the body json to struct and checking for errors
		// ShoudBindJSON shoud be used instead if BindJSON (this returns with Header text/json instead of text/plain)
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Server error : " + err.Error(),
			})
			return
		}

		// Creating mongo indexes for the users model in users database in mongodb (Can be shifter to Init function [should be])
		indexes := []mongo.IndexModel{
			{
				Keys:    bson.M{"email": 1},
				Options: options.Index().SetUnique(true),
			},
			// {
			// 	Keys:    bson.M{"name": 1},
			// 	Options: options.Index().SetUnique(true),
			// },
		}
		collection.Indexes().CreateMany(ctx, indexes)
		// Checking if the user already exsists (Depricated now I use indexes)
		// var check model.UserModel
		// filter := bson.D{{"email", user.Email}}
		// result := collection.FindOne(context.TODO(), filter).Decode(&check)
		// fmt.Println(result)
		// if check.Email != "" {
		// 	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		// 		"message": "User alread exsists cannot create user",
		// 	})
		// 	return
		// }

		// Encrypt the user password
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		user.Password = string(hash)

		// Save the user to the database
		_, err = collection.InsertOne(context.TODO(), user)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Error occured while savin user information : " + err.Error(),
			})
			return
		}

		// Generate and send the jwt back with response
		tok, err := createToken(user.Name)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Server error",
			})
			return
		}

		// Set cookies with token in them (expires in 15 minutes)
		ctx.SetCookie("token", tok, 15*60, "/", "localhost", false, true)

		// Send the user message
		ctx.AbortWithStatusJSON(200, gin.H{
			"message": "Success",
			"token":   tok,
		})

		ctx.Next()
	}
}
