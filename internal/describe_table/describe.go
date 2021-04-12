package describe_table

import (
	"context"
	"dynamo-hello-world/internal/connection"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func DescribeTable(table string) {
	client := connection.Connect()

	t, err := client.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{TableName: aws.String(table)})
	if err != nil {
		panic(err)
	}

	fmt.Printf("name: %s \n", *t.Table.TableName)
	fmt.Print("----- indexes -----\n")

	for _, i := range t.Table.GlobalSecondaryIndexes {
		fmt.Printf("name: %s \n", *i.IndexName)
		fmt.Printf("projection type: %s \n", i.Projection.ProjectionType)

		for k, v := range i.KeySchema {
			fmt.Printf("[key %d] attribute: %s type: %s \n", k, *v.AttributeName, v.KeyType)
		}
	}
}
