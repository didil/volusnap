package api

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type authSvcer interface {
	Signup(email string, password string) (uint, error)
	Login(email string, password string) (string, error)
}

func newAuthService(db *gorm.DB) *authService {
	return &authService{db}
}

type authService struct {
	db *gorm.DB
}

// Signup user
func (svc *authService) Signup(email string, password string) (uint, error) {
	var users []User

	err := svc.db.Where("email = ?", email).Find(&users).Error
	if err != nil {
		return 0, err
	}
	if len(users) > 0 {
		return 0, fmt.Errorf("a user already exists for the email: %v", email)
	}

	hashedP := hashAndSalt([]byte(password))
	user := User{Email: email, Password: hashedP}
	err = svc.db.Create(&user).Error
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

func comparePasswords(hashedPwd []byte, plainPwd []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPwd, plainPwd)
	return err
}

// Login user
func (svc *authService) Login(email string, password string) (string, error) {
	var users []User

	err := svc.db.Where("email = ?", email).Find(&users).Error
	if err != nil {
		return "", err
	}
	if len(users) == 0 {
		return "", fmt.Errorf("user not found for the email: %v", email)
	}

	user := users[0]

	err = comparePasswords([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("password invalid")
	}

	token, err := signJWT(user.ID)
	return token, err
}

type customClaims struct {
	UserID uint `json:"user"`
	jwt.StandardClaims
}

func getJWTSecret() ([]byte, error) {
	secret := []byte(viper.GetString("jwt.secret"))
	if len(secret) == 0 {
		return nil, fmt.Errorf("JWT signing Secret is empty")
	}

	return secret, nil
}

func signJWT(userID uint) (string, error) {
	secret, err := getJWTSecret()

	claims := customClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(365 * 24 * time.Hour).Unix(),
			Issuer:    "app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(secret)
	return tokenStr, err
}

func parseJWT(tokenStr string) (uint, error) {
	secret, err := getJWTSecret()

	token, err := jwt.ParseWithClaims(tokenStr, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("jwt token invalid")
	}

	return claims.UserID, nil
}
