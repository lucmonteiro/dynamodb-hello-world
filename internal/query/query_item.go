package query

import (
	"context"
	"dynamo-hello-world/internal/clock"
	"dynamo-hello-world/internal/connection"
	"dynamo-hello-world/internal/model"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetOrdersByCustomerAndDate_OnlyQueryingIndex(customerID string, c clock.Clock) ([]model.Order, error) {
	client := connection.Connect()

	keyExp := expression.Key("gsi1pk").
		Equal(expression.Value("CUSTOMER#" + customerID)).
		And(expression.Key("gsi1sk").BeginsWith("ORDERDATE#" + c.Now().Format("2006-01-02")))

	//filter registers that are not orders
	filterCondition := expression.Name("PK").BeginsWith("ORDER")

	cond, err := expression.NewBuilder().
		WithKeyCondition(keyExp).           //key condition for this query. must ONLY have keys defined on table OR gsi being used
		WithFilter(filterCondition).Build() //filter condition. can contain any attribute, as long it's projected on index
	if err != nil {
		return nil, err
	}

	model.PrintCondition(cond)

	result, err := client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:                 aws.String(model.TableName),
		ExpressionAttributeNames:  cond.Names(),
		ExpressionAttributeValues: cond.Values(),
		FilterExpression:          cond.Filter(),
		IndexName:                 aws.String(model.GSI1IndexName),
		KeyConditionExpression:    cond.KeyCondition(),
	})

	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}

	if result.Count == 0 {
		return nil, errors.New("no orders found")
	}

	out := []model.Order{}
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &out); err != nil {
		return nil, err
	}

	return out, nil
}

func GetOrdersByCustomerAndDate(customerID string, c clock.Clock) ([]model.Order, error) {
	client := connection.Connect()

	keyExp := expression.Key("gsi1pk").
		Equal(expression.Value("CUSTOMER#" + customerID)).
		And(expression.Key("gsi1sk").BeginsWith("ORDERDATE#" + c.Now().Format("2006-01-02")))

	//filter registers that are not orders
	filterCondition := expression.Name("PK").BeginsWith("ORDER")

	cond, err := expression.NewBuilder().
		WithKeyCondition(keyExp).           //key condition for this query. must ONLY have keys defined on table OR gsi being used
		WithFilter(filterCondition).Build() //filter condition. can contain any attribute, as long it's projected on index
	if err != nil {
		return nil, err
	}

	model.PrintCondition(cond)

	queryResult, err := client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:                 aws.String(model.TableName),
		ExpressionAttributeNames:  cond.Names(),
		ExpressionAttributeValues: cond.Values(),
		FilterExpression:          cond.Filter(),
		IndexName:                 aws.String(model.GSI1IndexName),
		KeyConditionExpression:    cond.KeyCondition(),
	})

	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}

	if queryResult.Count == 0 {
		return nil, errors.New("no orders found")
	}

	for _, r := range queryResult.Items {
		delete(r, "gsi1pk")
		delete(r, "gsi1sk")
	}

	result, err := client.BatchGetItem(context.Background(), &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			model.TableName: types.KeysAndAttributes{
				Keys:           queryResult.Items, //fetch me all those keys
				ConsistentRead: aws.Bool(true),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	out := []model.Order{}
	if err := attributevalue.UnmarshalListOfMaps(result.Responses[model.TableName], &out); err != nil {
		return nil, err
	}

	return out, nil
}

func GetOrdersByCustomer(customerID string) ([]model.Order, error) {
	client := connection.Connect()

	keyExp := expression.Key("gsi1pk").Equal(expression.Value("CUSTOMER#" + customerID))
	orderCond := expression.Name("PK").BeginsWith("ORDER")

	cond, err := expression.NewBuilder().WithKeyCondition(keyExp).WithFilter(orderCond).Build()
	if err != nil {
		return nil, err
	}

	model.PrintCondition(cond)

	result, err := client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:                 aws.String(model.TableName),
		ExpressionAttributeNames:  cond.Names(),
		ExpressionAttributeValues: cond.Values(),
		FilterExpression:          cond.Filter(),
		IndexName:                 aws.String(model.GSI1IndexName),
		KeyConditionExpression:    cond.KeyCondition(),
	})

	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}

	if result.Count == 0 {
		return nil, errors.New("no orders found")
	}

	out := []model.Order{}
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &out); err != nil {
		return nil, err
	}

	return out, nil
}
