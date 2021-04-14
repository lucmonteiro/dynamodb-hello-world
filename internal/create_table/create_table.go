package create_table

import (
	"context"
	"dynamo-hello-world/internal/connection"
	"dynamo-hello-world/internal/describe_table"
	"dynamo-hello-world/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func CreateTable(name string) error {
	client := connection.Connect()

	_, err := client.CreateTable(context.Background(), createTableInput(name))
	if err != nil {
		panic(err)
	}

	describe_table.DescribeTable(name)
	return err
}

func createTableInput(name string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions:   attributeDefinitions(),
		KeySchema:              keySchema(),
		TableName:              aws.String(name),
		GlobalSecondaryIndexes: globalSecondaryIndexes(),
		//Provisioned = You define maximum read/write capacity.
		//PayPerRequest = Dynamodb will auto scale according to use.
		BillingMode: types.BillingModeProvisioned,

		// ProvisionedThroughput of the table. Must be provided if BillingMode = Provisioned
		// If billingMode = PayPerRequest, is not necessary, as table will scale according to use.
		// It's recommended for free users to always use 1 here to avoid being charged.
		ProvisionedThroughput: provisionedThroughput(),
	}

}

//Defining Key for table
//At least one Hash key must be provided.
//If just Hash attribute is provided, unique items will be defined by Hash
//If Range key is provided, unique items will be defined by Hash+Range
func keySchema() []types.KeySchemaElement {
	return []types.KeySchemaElement{
		{
			AttributeName: aws.String("PK"),
			KeyType:       types.KeyTypeHash, //HASH must always be = on queries. It's the key type for "Partition Key"
		},
	}
}

//Defining indexes for table. Table might have up to 20 GSI's. It is advised to overload indexes
//in the case more is needed.
func globalSecondaryIndexes() []types.GlobalSecondaryIndex {
	return []types.GlobalSecondaryIndex{
		{
			IndexName: aws.String(model.GSI1IndexName),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("gsi1pk"),
					KeyType:       types.KeyTypeHash, //HASH must always be = on queries. It's the key type for "Partition Key"
				},
				{
					AttributeName: aws.String("gsi1sk"),
					KeyType:       types.KeyTypeRange, //RANGE is more flexibe. Can be >, <, begins_with on query, It's the type for "Sort Key"
				},
			},
			Projection: &types.Projection{
				//KEYS_ONLY = Only keys are projected
				//ALL = All attributes are projected. Use with caution, as this causes the table to be duplicated.
				//      One recommended use case by Amazon is to create an index with the same Key Schema of the table
				//      and projection = ALL to generate an eventually consistent read-only copy of the table.

				//INCLUDE = Select what attributes will be included on projection.
				ProjectionType: types.ProjectionTypeKeysOnly,

				//if projection = INCLUDE, fields that will be included must be provided here
				NonKeyAttributes: nil,
			},

			//Indexes have their own provisioned throughput
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(1),
				WriteCapacityUnits: aws.Int64(1),
			},
		},
	}
}

func provisionedThroughput() *types.ProvisionedThroughput {
	return &types.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(1),
		WriteCapacityUnits: aws.Int64(1),
	}
}

//Defining table attribute names. Not all attributes need to be defined here, only those that will be used
//on the Keys and Indexes
func attributeDefinitions() []types.AttributeDefinition {
	return []types.AttributeDefinition{
		{
			AttributeName: aws.String("PK"),
			AttributeType: types.ScalarAttributeTypeS, //AttributeTypeS = string
		},
		{
			AttributeName: aws.String("gsi1pk"),
			AttributeType: types.ScalarAttributeTypeS, //AttributeTypeS = string
		},
		{
			AttributeName: aws.String("gsi1sk"),
			AttributeType: types.ScalarAttributeTypeS, //AttributeTypeS = string
		},
	}
}
