package api

import (
	"database/sql"

	"github.com/didil/volusnap/pkg/models"
)

type accountSvcer interface {
	List(userID int) (models.AccountSlice, error)
}

func newAccountService(db *sql.DB) *accountService {
	return &accountService{db}
}

type accountService struct {
	db *sql.DB
}

func (svc *accountService) List(userID int) (models.AccountSlice, error) {
	accounts, err := models.Accounts(models.AccountWhere.UserID.EQ(userID)).All(svc.db)
	return accounts, err
}
