package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func FindTerraformTables(ctx context.Context, svc *dynamodb.Client) ([]string, error) {
	input := &dynamodb.ListTablesInput{}
	output, err := svc.ListTables(ctx, input)
	if err != nil {
		return nil, err
	}

	var terraformTables []string
	for _, tableName := range output.TableNames {
		if strings.Contains(strings.ToLower(tableName), "terraform") {
			table, err := svc.DescribeTable(ctx, &dynamodb.DescribeTableInput{
				TableName: aws.String(tableName),
			})
			if err != nil {
				return nil, err
			}
			if table.Table.TableStatus == types.TableStatusActive {
				if table.Table.BillingModeSummary != nil && table.Table.BillingModeSummary.BillingMode != types.BillingModePayPerRequest {
					terraformTables = append(terraformTables, fmt.Sprintf("%s. Warning: this table may incur additional costs as it is not set to on-demand mode", tableName))
				}
			}
		}
	}
	return terraformTables, nil
}
