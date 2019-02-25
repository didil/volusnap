package api

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/didil/volusnap/pkg/models"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type authSvcer interface {
	Signup(email string, password string) (int, error)
	Login(email string, password string) (string, error)
}

func newAuthService(db *sql.DB) *authService {
	return &authService{db}
}

type authService struct {
	db *sql.DB
}

// Signup user
func (svc *authService) Signup(email string, password string) (int, error) {
	count, err := models.Users(qm.Where("email = ?", email)).Count(svc.db)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, fmt.Errorf("a user already exists for the email: %v", email)
	}

	hashedP := hashAndSalt([]byte(password))
	user := models.User{Email: email, Password: hashedP}
	err = user.Insert(svc.db, boil.Infer())
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
	users, err := models.Users(models.UserWhere.Email.EQ(email)).All(svc.db)
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
	UserID int `json:"user"`
	jwt.StandardClaims
}

func getJWTSecret() ([]byte, error) {
	secret := []byte(viper.GetString("jwt.secret"))
	if len(secret) == 0 {
		return nil, fmt.Errorf("JWT signing Secret is empty")
	}

	return secret, nil
}

func signJWT(userID int) (string, error) {
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

func parseJWT(tokenStr string) (int, error) {
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
