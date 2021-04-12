package get_item

import (
	"context"
	"dynamo-hello-world/internal/connection"
	"dynamo-hello-world/internal/model"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetItem(id string) (model.Customer, error) {
	client := connection.Connect()

	result, err := client.GetItem(context.Background(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "CUSTOMER"},
			"SK": &types.AttributeValueMemberS{Value: "CUSTOMER#" + id},
		},
		TableName: aws.String(model.TableName),
	})

	if err != nil {
		return model.Customer{}, err
	}

	if result.Item == nil {
		return model.Customer{}, errors.New("not found")
	}

	out := model.Customer{}
	if err := attributevalue.UnmarshalMap(result.Item, &out); err != nil {
		return model.Customer{}, err
	}

	return out, nil
}
