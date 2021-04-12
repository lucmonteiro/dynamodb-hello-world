package create_table_test

import (
	"dynamo-hello-world/internal/create_table"
	"dynamo-hello-world/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateTable(t *testing.T) {
	err := create_table.CreateTable(model.TableName)
	require.Nil(t, err)
}
