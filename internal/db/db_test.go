package db_test

import (
	"goRoadMap/internal/db"
	"goRoadMap/internal/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitiateTables(t *testing.T) {
	logger.Init("info")

	// проверка, выполняется ли функция без ошибок
	err := db.InitiateTables()
	assert.Nil(t, err, "Expected nil error, got: %v", err)
}
