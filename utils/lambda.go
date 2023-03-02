package utils

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

func FindTotalCountOfLambdaInvocs(ctx context.Context, region, lambdaName string) ([]string, error) {

	// Create a new AWS session with the desired region
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		panic(err)
	}

	// Create a CloudWatch client using the session
	cwClient := cloudwatch.NewFromConfig(cfg)

	// Call the GetMetricData API to get the number of Lambda function invocations for the last week
	input := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: []types.MetricDataQuery{
			{
				Id: aws.String("m1"),
				MetricStat: &types.MetricStat{
					Metric: &types.Metric{
						Namespace:  aws.String("AWS/Lambda"),
						MetricName: aws.String("Invocations"),
						Dimensions: []types.Dimension{
							{
								Name:  aws.String("FunctionName"),
								Value: aws.String(lambdaName),
							},
						},
					},
					Period: aws.Int32(86400), // 1 day
					Stat:   aws.String("Sum"),
				},
				ReturnData: aws.Bool(true),
			},
		},
		// StartTime: aws.Time(time.Now().Add(-3 * time.Hour)), // 3 hours ago for example
		StartTime: aws.Time(time.Now().Add(-7 * 24 * time.Hour)), // 1 week ago
		EndTime:   aws.Time(time.Now()),
		ScanBy:    types.ScanByTimestampDescending,
	}
	output, err := cwClient.GetMetricData(context.Background(), input)
	if err != nil {
		panic(err)
	}

	var unusedLambdas []string
	// Extract the value from the output
	if len(output.MetricDataResults[0].Values) == 0 {
		unusedLambdas = append(unusedLambdas, lambdaName)
	}
	return unusedLambdas, nil
}
