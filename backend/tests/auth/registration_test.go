package auth_tests

import (
	"CoggersProject/internal/app/lib/jwt"
	pbAuth "CoggersProject/pkg/protos/auth"
	"CoggersProject/tests/suite"
	"log"
	"testing"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	ctx, st := suite.New(t)
	ClearTable("users")

	RegReq := &pbAuth.RegistrationRequest{
		Name:      gofakeit.Name(),
		Birthdate: gofakeit.Date().String(),
		Photourl:  "testing",
		Push:      gofakeit.Bool(),
		Email:     gofakeit.Email(),
		Password:  gofakeit.Password(true, true, true, true, false, 10),
		RoleId:    1,
	}

	exceptedKey, _ := jwt.Keygen(RegReq.Email, RegReq.Password, st.Env.JwtSecret)

	t.Run("NewUserRegistration Test", func(t *testing.T) {
		response, err := st.AuthClient.Registration(ctx, RegReq)
		if err != nil {
			log.Fatalf("Error when calling Registration: %v", err)
		}

		assert.Equal(t, exceptedKey, response.Key)
	})
}
