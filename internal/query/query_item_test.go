package query_test

import (
	"dynamo-hello-world/internal/clock/clocktest"
	"dynamo-hello-world/internal/model"
	"dynamo-hello-world/internal/put_item"
	"dynamo-hello-world/internal/query"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func buildOrder(customerId, orderId string, date time.Time) model.Order {
	return model.Order{
		ID:     orderId,
		Date:   &date,
		GSI1SK: "ORDERDATE#" + date.Format(time.RFC3339),
		Customer: model.Customer{
			ID: customerId,
		},
	}
}

func TestGetOrdersByCustomer(t *testing.T) {

	require.Nil(t, put_item.CreateNewCustomer("celMarluslaceWal"))

	orderDate := clocktest.NewMockWithTime(time.Date(2006, 01, 02, 0, 0, 0, 0, time.UTC)).Now()
	o := buildOrder("celMarluslaceWal", "order1", orderDate)
	require.Nil(t, put_item.CreateNewOrder(o))

	o = buildOrder("celMarluslaceWal", "order2", orderDate)
	require.Nil(t, put_item.CreateNewOrder(o))

	orders, err := query.GetOrdersByCustomer("celMarluslaceWal")
	require.Nil(t, err)
	require.Equal(t, 2, len(orders))

}

func TestGetOrdersByCustomerAndDate(t *testing.T) {

	require.Nil(t, put_item.CreateNewCustomer("jules"))

	orderDate := clocktest.NewMockWithTime(time.Date(1995, 02, 18, 0, 0, 0, 0, time.UTC)).Now()
	o := buildOrder("jules", "royale-with-cheese", orderDate)
	require.Nil(t, put_item.CreateNewOrder(o))

	orderDate = clocktest.NewMockWithTime(time.Date(1995, 02, 19, 0, 0, 0, 0, time.UTC)).Now()
	o = buildOrder("jules", "five-dollar-shake", orderDate)
	require.Nil(t, put_item.CreateNewOrder(o))

	orders, err := query.GetOrdersByCustomerAndDate("jules", clocktest.NewMockWithTime(orderDate))
	require.Nil(t, err)

	require.Equal(t, 1, len(orders))

}
