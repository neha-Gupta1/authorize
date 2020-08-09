package controller

import (
	"authorize/models"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

)

var jwtKey = []byte("my_secret_key")

// SetupServer ....
func SetupServer() *gin.Engine {
	r := gin.Default()

	// Routes
	r.POST("/signup", CreateUser)
	r.POST("/login", LoginUser)
	r.GET("/validate", Validate)

	return r
}

// CreateUser user creation
// POST /signup
// Create new user
func CreateUser(c *gin.Context) {
	// Validate input
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil { // validate payload
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{UserName: input.UserName, Password: input.Password}
	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// LoginUser create cab booking
// POST /login
// Login User
func LoginUser(c *gin.Context) {
	// Validate input
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil { // validate payload
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{UserName: input.UserName, Password: input.Password}
	if err := models.DB.Where("user_name = ? AND password = ?", input.UserName, input.Password).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Internal server error!"})
		return
	}
	token, err := generateJWT(input.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": token})
}

//Claims ...
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func generateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

//Validate function
// GET /validate
func Validate(c *gin.Context) {
	//tknStr := c.Param("id")
	// Initialize a new instance of `Claims`
	claims := &Claims{}
	tknStr := c.Request.Header["Token"][0]
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invlid token!"})
			return
		}
		log.Println("Err:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invlid token!"})
		return
	}
	if !tkn.Valid {
		log.Println("tkn.Valid:", tkn.Valid)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invlid token!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "authenticated user!"})
}
