package auth_tests

import (
	pbAuth "CoggersProject/gen/go/auth"
	"CoggersProject/internal/app/lib/jwt"
	"CoggersProject/tests/suite"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	ctx, st := suite.New(t)
	UserData, err := MockUser(1)
	if err != nil {
		log.Fatalf("Ошибка при добавлении тестового админа в бд: %v", err)
	}

	loginReq := &pbAuth.LoginRequest{
		Email:    UserData["email"].(string),
		Password: UserData["password"].(string),
	}

	exceptedKey, _ := jwt.Keygen(loginReq.Email, loginReq.Password, st.Cfg.JwtSecret)

	t.Run("UserLogin Test", func(t *testing.T) {
		response, err := st.AuthClient.Login(ctx, loginReq)
		if err != nil {
			log.Fatalf("Error when calling Login: %v", err)
		}

		assert.Equal(t, exceptedKey, response.Key)
	})
}
