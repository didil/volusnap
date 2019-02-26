package api

import (
	"database/sql"
	"fmt"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/volatiletech/sqlboiler/boil"

	"github.com/didil/volusnap/pkg/models"
)

type accountSvcer interface {
	List(userID int) (models.AccountSlice, error)
	Create(userID int, provider string, name string, token string) (int, error)
	GetForUser(userID, accountID int) (*models.Account, error)
	Get(accountID int) (*models.Account, error)
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

func (svc *accountService) Create(userID int, provider string, name string, token string) (int, error) {
	if !pRegistry.isValidProvider(provider) {
		return 0, fmt.Errorf("invalid provider: %v", provider)
	}
	if token == "" {
		return 0, fmt.Errorf("empty token")
	}

	account := models.Account{UserID: userID, Provider: provider, Name: name, Token: token}
	err := account.Insert(svc.db, boil.Infer())

	if err != nil {
		return 0, err
	}

	return account.ID, nil
}

func (svc *accountService) GetForUser(userID, accountID int) (*models.Account, error) {
	account, err := models.Accounts(
		qm.Where("user_id = ?", userID),
		qm.Where("id = ?", accountID),
	).One(svc.db)
	return account, err
}

func (svc *accountService) Get(accountID int) (*models.Account, error) {
	account, err := models.Accounts(
		qm.Where("id = ?", accountID),
	).One(svc.db)
	return account, err
}
