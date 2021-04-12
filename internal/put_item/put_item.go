package put_item

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
	"github.com/gofrs/uuid"
	"time"
)

func CreateNewOrder(o model.Order) error {
	client := connection.Connect()
	addDynamoKeysToOrders(&o)

	_, err := client.PutItem(context.Background(), &dynamodb.PutItemInput{
		Item:      marshalItem(o),
		TableName: aws.String(model.TableName),
	})

	return err
}

func CreateNewCustomer(id string) error {
	client := connection.Connect()

	item := buildSampleCustomer(id)
	_, err := client.PutItem(context.Background(), &dynamodb.PutItemInput{
		Item:      marshalItem(item),
		TableName: aws.String(model.TableName),
	})

	return err
}

// LockItem will update item with a lock TTL and token.
// This function exercises the ConditionExists
// This is a simple lock that will NOT validate if item is previously locked.
func LockItem(c clock.Clock, id string) (string, error) {
	client := connection.Connect()

	expr := expression.AttributeExists(expression.Name("PK")).
		And(expression.AttributeExists(expression.Name("SK")))

	cond, err := expression.NewBuilder().WithCondition(expr).Build()
	if err != nil {
		panic(err)
	}

	model.PrintCondition(cond)

	item := buildSampleCustomer(id)
	lockToken, _ := uuid.NewV4()
	item.LockToken = lockToken.String()
	lockTime := c.Now().Add(time.Minute * 5)
	item.LockUntil = &lockTime

	_, err = client.PutItem(context.Background(), &dynamodb.PutItemInput{
		Item:      marshalItem(item),
		TableName: aws.String(model.TableName),
		//builds the Dynamo condition string.
		//it will look something like (attribute_exists (#0)) AND (attribute_exists (#1))
		ConditionExpression: cond.Condition(),
		//tells dynamo the attributes names for this query. #0 = PK #1 = SK
		ExpressionAttributeNames: cond.Names(),
		//tells dynamo the attribute values (if any). In this case will be blank, as we're only validating a exists conditions (example: update)
		ExpressionAttributeValues: cond.Values(),
	})

	var conditionalErr *types.ConditionalCheckFailedException
	if err != nil {
		if errors.As(err, &conditionalErr) {
			fmt.Printf("conditional error: %s \n", *conditionalErr.Message)
		} else {
			panic(err)
		}
	}

	return lockToken.String(), err
}

//Update only if token is valid and lock_until > now
// Not generating versions
func Update(c clock.Clock, id, lockToken string) error {
	client := connection.Connect()

	item := buildSampleCustomer(id)
	item.Name = "4d72416e646572736f6e"
	item.LastUpdateAt = c.Now()

	//Ensure that keys exist
	expr := expression.AttributeExists(expression.Name("PK")).
		And(expression.AttributeExists(expression.Name("SK"))).
		//ensure that provided token is the on on database
		And(expression.Name("lock_until").GreaterThan(expression.Value(c.Now())).
			And(expression.Name("token").Equal(expression.Value(lockToken))))
	//if token was not provided, register must not be locked

	cond, err := expression.NewBuilder().WithCondition(expr).Build()
	if err != nil {
		panic(err)
	}

	model.PrintCondition(cond)

	_, err = client.PutItem(context.Background(), &dynamodb.PutItemInput{
		Item:      marshalItem(item),
		TableName: aws.String(model.TableName),
		//builds the Dynamo condition string.
		//it will look something like (attribute_exists (#0)) AND (attribute_exists (#1))
		ConditionExpression: cond.Condition(),
		//tells dynamo the attributes names for this query. #0 = PK #1 = SK
		ExpressionAttributeNames: cond.Names(),
		//tells dynamo the attribute values (if any). In this case will be blank, as we're only validating a exists conditions (example: update)
		ExpressionAttributeValues: cond.Values(),
	})

	var conditionalErr *types.ConditionalCheckFailedException
	if err != nil {
		if errors.As(err, &conditionalErr) {
			fmt.Printf("conditional error: %s \n", *conditionalErr.Message)
		} else {
			panic(err)
		}
	}

	return err

}

func buildSampleCustomer(id string) model.Customer {
	customerItem := model.Customer{
		ID: id,
	}
	customerItem.PK = "CUSTOMER"
	customerItem.SK = "CUSTOMER#" + customerItem.ID
	customerItem.GSI1SK = "LATEST"

	return customerItem
}

func addDynamoKeysToOrders(order *model.Order) {
	order.PK = "ORDER#" + order.ID
	order.SK = "CUSTOMER#" + order.Customer.ID
	order.GSI1SK = "ORDERDATE#" + order.Date.Format(time.RFC3339)
}

func marshalItem(item interface{}) map[string]types.AttributeValue {
	marshaled, err := attributevalue.MarshalMap(item)
	if err != nil {
		panic(err)
	}

	return marshaled
}
