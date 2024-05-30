package auth_test

import (
	"CoggersProject/backend/config"
	"CoggersProject/backend/internal/endpoints/auth"
	"CoggersProject/backend/internal/errorz"
	"CoggersProject/backend/internal/mw"
	"CoggersProject/backend/pkg/jwtAuth"
	"CoggersProject/backend/pkg/service/cacher"
	"CoggersProject/backend/pkg/service/logger"
	"log"
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

func TestAuth(t *testing.T) {
	config.Init()
	config := config.GetConfig()

	err := godotenv.Load(config.EnvPath)

	if err != nil {
		log.Fatal("Ошибка при открытии .env файла: ", err)
	}

	logger.Init(config.LoggerMode)
	cacher.Init(config.Cache.UpdateInterval)

	a := &App{}
	a.s = jwtAuth.New()
	a.e = auth.New(a.s)
	a.echo = echo.New()
	a.echo.Use(mw.Recovery)

	t.Run("NewUserRegistration Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/NewUserRegistration", nil)
		req.Header.Set("username", "testuser")
		req.Header.Set("password", "testpass")
		rec := httptest.NewRecorder()
		ctx := a.echo.NewContext(req, rec)

		err := a.e.NewUserRegistration(ctx)
		if err != errorz.ErrUserExists {
			assert.NoError(t, err, "не ожидали ошибку, получили: ", err)
		}
	})

	t.Run("UserLogin Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/UserLogin", nil)
		req.Header.Set("username", "testuser")
		req.Header.Set("password", "testpass")
		rec := httptest.NewRecorder()
		ctx := a.echo.NewContext(req, rec)

		err := a.e.UserLogin(ctx)
		assert.NoError(t, err, "не ожидали ошибку, получили: ", err)
	})
}
