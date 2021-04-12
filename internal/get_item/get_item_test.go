package get_item

import (
	"dynamo-hello-world/internal/put_item"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetItem(t *testing.T) {
	err := put_item.CreateNewCustomer("bale")
	require.Nil(t, err)

	item, err := GetItem("bale")
	require.Nil(t, err)
	assert.Equal(t, "bale", item.ID)
}

func TestNotFound(t *testing.T) {
	item, err := GetItem("nolan")
	require.NotNil(t, err)
	assert.Empty(t, item)

	item, err = GetItem("jackman")
	require.NotNil(t, err)
	assert.Empty(t, item)
}
