package put_item

import (
	"dynamo-hello-world/internal/clock/clocktest"
	"dynamo-hello-world/internal/model"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var c = clocktest.NewMockWithTime(time.Date(1999, 05, 21, 0, 0, 0, 0, time.UTC))

func TestPutItem(t *testing.T) {
	id, _ := uuid.NewV4()
	require.Nil(t, CreateNewCustomer(id.String()))
}

func TestCreateOrder(t *testing.T) {
	orderDate := c.Now()

	orderId, _ := uuid.NewV4()

	require.Nil(t, CreateNewCustomer("Thomas Anderson"))

	err := CreateNewOrder(model.Order{
		ID:     orderId.String(),
		Date:   &orderDate,
		GSI1SK: "ORDERDATE#" + orderDate.Format(time.RFC3339),
		Customer: model.Customer{
			ID: "Thomas Anderson",
		},
	})

	assert.Nil(t, err)
}

func TestLockItem(t *testing.T) {
	idForTest := "4d72416e646572736f6e"

	require.Nil(t, CreateNewCustomer(idForTest))
	lock, err := LockItem(c, idForTest)
	require.Nil(t, err)

	//time is < ttl, so its locked
	//since provided token is no good, will not update
	beforeTTL := time.Date(1977, 05, 25, 0, 0, 0, 0, time.UTC)

	//update with token will work
	err = Update(clocktest.NewMockWithTime(beforeTTL), idForTest, lock)
	assert.Nil(t, err, "update before ttl with token should not error")

}

func TestExpiredLock(t *testing.T) {
	idForTest := "nebuchadnezzar"

	require.Nil(t, CreateNewCustomer(idForTest))
	lock, err := LockItem(c, idForTest)
	require.Nil(t, err)

	//time is > ttl, and token is valid
	expiredTTL := time.Date(2002, 01, 01, 0, 0, 0, 0, time.UTC)
	c2 := clocktest.NewMockWithTime(expiredTTL)
	err = Update(c2, idForTest, lock)
	assert.NotNil(t, err, "lock is expired, should not update")
}

func TestInvalidToken(t *testing.T) {
	idForTest := "deep rabbit hole"

	require.Nil(t, CreateNewCustomer(idForTest))
	_, err := LockItem(c, idForTest)
	require.Nil(t, err)

	//time is < ttl, so its locked
	//since provided token is no good, will not update
	beforeTTL := time.Date(1977, 05, 25, 0, 0, 0, 0, time.UTC)
	err = Update(clocktest.NewMockWithTime(beforeTTL), idForTest, "invalidToken")
	assert.NotNil(t, err, "should update with invalid token")
}

func TestConditionNotExistsPutItem(t *testing.T) {
	idForTest := "architect"

	_, err := LockItem(c, idForTest)
	require.NotNil(t, err)

	require.Nil(t, CreateNewCustomer(idForTest))
	lock, err := LockItem(c, idForTest)

	require.Nil(t, err)
	require.NotEmpty(t, lock)
}
