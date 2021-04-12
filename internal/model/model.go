package model

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"time"
)

const (
	TableName     = "test-table"
	GSI1IndexName = "gsi1"
)

type DynamoKey struct {
	PK string
	SK string
}

type Customer struct {
	DynamoKey
	ID           string     `dynamodbav:"id"`
	Name         string     `dynamodbav:"name"`
	LockToken    string     `dynamodbav:"token,omitempty"`
	LockUntil    *time.Time `dynamodbav:"lock_until,omitempty"`
	LastUpdateAt time.Time  `dynamodbav:"last_update_at"`
	GSI1SK       string     `dynamodbav:"gsi1sk"`
}

type Order struct {
	DynamoKey
	ID       string     `dynamodbav:"id"`
	Date     *time.Time `dynamodbav:"date"`
	GSI1SK   string     `dynamodbav:"gsi1sk"`
	Customer Customer   `dynamodbav:"-"`
}

func PrintCondition(cond expression.Expression) {
	fmt.Printf("\n\n-------- PRINTING CONDITION -----------\n\n")
	fmt.Printf("names:    \t %s \n", cond.Names())
	fmt.Printf("values:   \t %s \n", cond.Values())
	if cond.Condition() != nil {
		fmt.Printf("condition:\t %s \n", *cond.Condition())
	}
	if cond.KeyCondition() != nil {
		fmt.Printf("key condition:\t %s \n", *cond.KeyCondition())
	}

	fmt.Printf("\n-------- PRINTING CONDITION -----------\n\n")
}
