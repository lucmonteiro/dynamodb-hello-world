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
	model.AddDynamoKeysToOrders(&o)

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

	expr := expression.AttributeExists(expression.Name("PK"))

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
		//tells dynamo the attributes names for this query. #0 = PK #1 = GSI1PK
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

	expr := buildLockCondition(c, lockToken)
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
		//tells dynamo the attributes names for this query. #0 = PK #1 = GSI1PK
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

func buildLockCondition(c clock.Clock, lockToken string) expression.ConditionBuilder {
	//Ensure that keys exist
	expr := expression.AttributeExists(expression.Name("PK")).
		//ensure that provided token is the on on database
		And(expression.Name("lock_until").GreaterThan(expression.Value(c.Now())).
			And(expression.Name("token").Equal(expression.Value(lockToken))))
	return expr
}

func buildSampleCustomer(id string) model.Customer {
	customerItem := model.Customer{
		ID: id,
	}

	model.AddDynamoKeysToCustomer(&customerItem)

	return customerItem
}

func marshalItem(item interface{}) map[string]types.AttributeValue {
	marshaled, err := attributevalue.MarshalMap(item)
	if err != nil {
		panic(err)
	}

	return marshaled
}
