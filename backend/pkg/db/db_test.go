package db_test

import (
	"CoggersProject/backend/config"
	"CoggersProject/backend/internal/endpoints/auth"
	"CoggersProject/backend/internal/mw"
	"CoggersProject/backend/pkg/cache"
	"CoggersProject/backend/pkg/db"
	"CoggersProject/backend/pkg/jwtAuth"
	"CoggersProject/backend/pkg/service/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type App struct {
	e    *auth.Endpoint
	s    *jwtAuth.Service
	echo *echo.Echo
}

func TestInitiateTables(t *testing.T) {
	config.Init()
	config := config.GetConfig()
	logger.Init(config.LoggerMode)
	godotenv.Load(config.EnvPath)
	cache.ClearCache()

	a := &App{}
	a.s = jwtAuth.New()
	a.e = auth.New(a.s)
	a.echo = echo.New()
	a.echo.Use(mw.Recovery)

	expectedUserData := map[string]string{
		"username": "testuser",
		"password": "testpass",
		"roleid":   "1",
	}

	req := httptest.NewRequest(http.MethodGet, "/NewUserRegistration", nil)
	req.Header.Set("username", "testuser")
	req.Header.Set("password", "testpass")
	rec := httptest.NewRecorder()
	ctx := a.echo.NewContext(req, rec)
	a.e.NewUserRegistration(ctx)

	t.Run("InitiateTables Test", func(t *testing.T) {
		err := db.InitiateTables()
		assert.Nil(t, err, "Не ожидаем ошибку, получили: %v", err)
	})

	t.Run("FetchDataFromDB Test", func(t *testing.T) {
		userData, err := db.FetchUserData("testuser")

		assert.Nil(t, err, "Не ожидаем ошибку, получили: %v", err)
		assert.Equalf(t, expectedUserData, userData, "Ожидали %v, получили %v", expectedUserData, userData)
	})
}
