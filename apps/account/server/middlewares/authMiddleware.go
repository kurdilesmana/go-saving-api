package middlewares

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/models/accountModel"
	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/ports/accountPort"
	"github.com/kurdilesmana/go-saving-api/pkg/logging"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthMiddleware struct {
	AccountRepo accountPort.IAccountRepository
	log         *logging.Logger
}

func NewAuthMiddleware(
	accountRepo accountPort.IAccountRepository,
	log *logging.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		AccountRepo: accountRepo,
		log:         log,
	}
}

const (
	basic = "basic"
)

func (m *AuthMiddleware) ValidateAccount(account_no, pin string, c echo.Context) (valid bool, err error) {
	// validate account_no
	m.log.Info(logrus.Fields{"account_no": account_no}, nil, "validate account_no")
	account, err := m.AccountRepo.GetByAccountNo(c.Request().Context(), account_no)
	if err != nil {
		m.log.Warn(logrus.Fields{}, nil, err.Error())
		err = fmt.Errorf("internal error, failed validated account no")
		valid = false
		return
	}

	// user info
	m.log.Info(logrus.Fields{"account_no": account_no}, account, "account data")
	if account.AccountNo == "" {
		m.log.Warn(logrus.Fields{}, nil, "invalid account")
		err = fmt.Errorf("invalid account")
		valid = false
		return
	}

	// validate pin / password
	byteHash := []byte(account.PIN)
	err = bcrypt.CompareHashAndPassword(byteHash, []byte(pin))
	if err != nil {
		m.log.Warn(logrus.Fields{}, nil, err.Error())
		err = fmt.Errorf("invalid pin")
		valid = false
		return
	}

	// set info user session
	accountPayload := accountModel.AccountPayload{
		AccountNo:   account.AccountNo,
		Name:        account.Name,
		PhoneNumber: account.PhoneNumber,
		Balance:     account.Balance,
	}
	c.Set("session", &accountPayload)

	return true, nil
}

func (m *AuthMiddleware) AuthorizeToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get(echo.HeaderAuthorization)
			l := len(basic)

			if len(auth) > l+1 && strings.EqualFold(auth[:l], basic) {
				// Invalid base64 shouldn't be treated as error
				// instead should be treated as invalid client input
				b, err := base64.StdEncoding.DecodeString(auth[l+1:])
				if err != nil {
					return c.JSON(http.StatusBadRequest, map[string]interface{}{
						"message": "bad request",
						"data":    nil,
						"error":   err.Error(),
					})
				}

				cred := string(b)
				for i := 0; i < len(cred); i++ {
					if cred[i] == ':' {
						// Verify credentials
						valid, err := m.ValidateAccount(cred[:i], cred[i+1:], c)
						if err != nil {
							return c.JSON(http.StatusUnauthorized, map[string]interface{}{
								"message": "unauthorized",
								"data":    nil,
								"error":   err.Error(),
							})
						} else if valid {
							return next(c)
						}
						break
					}
				}
			}

			return nil
		}
	}
}
