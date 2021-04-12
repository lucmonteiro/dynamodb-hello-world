package describe_table

import (
	"dynamo-hello-world/internal/model"
	"testing"
)

func TestDescribeTable(t *testing.T) {
	DescribeTable(model.TableName)
}
