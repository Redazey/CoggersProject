package servparser_test

import (
	"CoggersProject/tests/suite"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestGetServersInfo(t *testing.T) {
	ctx, st := suite.New(t)

	serversInfo, err := st.ServParserClient.GetServersInfo(ctx, &emptypb.Empty{})
	assert.Nil(t, err, "неожиданная ошибка: ", zap.Error(err))

	fmt.Printf("%s", serversInfo)
}
