package utils

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func DescribeRegions(ctx context.Context, svc *ec2.Client) ([]string, error) {
	input := &ec2.DescribeRegionsInput{}
	output, err := svc.DescribeRegions(ctx, input)
	if err != nil {
		return nil, err
	}

	var regions []string
	for _, r := range output.Regions {
		regions = append(regions, aws.ToString(r.RegionName))
	}
	return regions, nil
}
