package utils

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

func CountDBConnections(ctx context.Context, region, dbInstanceID string) ([]string, error) {

	// Create a new AWS session with the desired region
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		panic(err)
	}

	// Create a CloudWatch client using the session
	cwClient := cloudwatch.NewFromConfig(cfg)

	// Set up the input parameters for the GetMetricData API call
	input := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: []types.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &types.MetricStat{
					Metric: &types.Metric{
						Namespace:  aws.String("AWS/RDS"),
						MetricName: aws.String("DatabaseConnections"),
						Dimensions: []types.Dimension{
							{
								Name:  aws.String("DBInstanceIdentifier"),
								Value: aws.String(dbInstanceID),
							},
						},
					},
					Period: aws.Int32(60), // 1 minute
					Stat:   aws.String("Sum"),
				},
				ReturnData: aws.Bool(true),
			},
		},
		StartTime:     aws.Time(time.Now().Add(-1 * time.Hour)), // 1 hour ago
		EndTime:       aws.Time(time.Now()),
		ScanBy:        types.ScanByTimestampDescending,
		MaxDatapoints: aws.Int32(1),
	}

	// Call the GetMetricData API to get the latest database connection count
	output, err := cwClient.GetMetricData(context.Background(), input)
	if err != nil {
		panic(err)
	}
	var unusedRDS []string

	// Extract the value from the output
	if len(output.MetricDataResults) > 0 && len(output.MetricDataResults[0].Values) > 0 {
		value := output.MetricDataResults[0].Values[0]
		if int(value) == 0 {
			unusedRDS = append(unusedRDS, dbInstanceID)
		}
	}
	return unusedRDS, nil

}
