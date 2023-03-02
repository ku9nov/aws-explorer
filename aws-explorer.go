package main

import (
	"aws-explorer/utils"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/rds"

	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

func main() {
	var logMode aws.ClientLogMode

	// Load AWS configuration from environment variables, shared config, or EC2 instance metadata
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		logMode = aws.LogResponseWithBody | aws.LogRequestWithBody
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithClientLogMode(logMode))

	if err != nil {
		fmt.Println("Error loading AWS configuration:", err)
		return
	}

	// Call the DescribeRegions API to get a list of all available regions
	regions, err := utils.DescribeRegions(context.Background(), ec2.NewFromConfig(cfg))
	if err != nil {
		fmt.Println("Error describing regions:", err)
		return
	}

	// Iterate through each region and find unused Elastic IP addresses
	for _, r := range regions {

		// Find unused ElasticIP
		unusedIps, err := utils.FindUnusedElasticIps(context.Background(), ec2.NewFromConfig(cfg, func(options *ec2.Options) {
			options.Region = r
		}))
		if err != nil {
			fmt.Printf("Error finding unused Elastic IPs in %s: %v\n", r, err)
			continue
		}

		// Print out the list of unused ElasticIPs
		for _, ip := range unusedIps {
			fmt.Printf("Unused Elastic IP address found in %s: %s\n", r, ip)
		}

		// Find unused EBS volumes
		unusedVolumes, err := utils.FindUnusedEBSVolumes(context.Background(), ec2.NewFromConfig(cfg, func(options *ec2.Options) {
			options.Region = r
		}))
		if err != nil {
			fmt.Printf("Error finding unused EBS volumes in %s: %v\n", r, err)
			continue
		}

		// Print out the list of unused EBS volumes
		for _, v := range unusedVolumes {
			fmt.Printf("Unused EBS volume found in %s: %s\n", r, v)
		}

		// Find unused ENI
		unusedENI, err := utils.FindUnusedNetworkInterfaces(context.Background(), ec2.NewFromConfig(cfg, func(options *ec2.Options) {
			options.Region = r
		}))
		if err != nil {
			fmt.Printf("Error finding unused ENI in %s: %v\n", r, err)
			continue
		}

		// Print out the list of unused ENI
		for _, v := range unusedENI {
			fmt.Printf("Unused ENI found in %s: %s\n", r, v)
		}

		// Find unused target groups
		unusedTargetGroups, err := utils.FindUnusedTargetGroups(context.Background(), elbv2.NewFromConfig(cfg, func(options *elbv2.Options) {
			options.Region = r
		}))
		if err != nil {
			fmt.Println("Error finding unused target groups:", err)
			continue
		}

		// Print out the list of unused target groups
		for _, tg := range unusedTargetGroups {
			fmt.Printf("Unused target group found in %s: %s\n", r, tg)
		}

		// Find DynamoDB tables that have the word "terraform" in their name
		terraformTables, err := utils.FindTerraformTables(context.Background(), dynamodb.NewFromConfig(cfg, func(options *dynamodb.Options) {
			options.Region = r
		}))
		if err != nil {
			fmt.Println("Error finding Terraform tables:", err)
			return
		}

		// Print out the list of Terraform tables with non on-demand billing mode
		for _, table := range terraformTables {
			fmt.Printf("Terraform table found: %s\n", table)
		}

		// Find Load Balancers without listeners
		unusedLoadBalancers, err := utils.FindUnusedLoadBalancers(context.Background(), elbv2.NewFromConfig(cfg, func(options *elbv2.Options) {
			options.Region = r
		}))
		if err != nil {
			fmt.Println("Error finding unused Load Balancers:", err)
			continue
		}

		// Print out the list of unused load balancers
		for _, lb := range unusedLoadBalancers {
			fmt.Printf("Unused load balancer found in %s: %s\n", r, lb)
		}

		// Create an RDS client for the current region
		rdsSvc := rds.NewFromConfig(cfg, func(options *rds.Options) {
			options.Region = r
		})

		// Call the DescribeDBInstances API to get a list of all RDS instances in the region
		instances, err := rdsSvc.DescribeDBInstances(context.Background(), &rds.DescribeDBInstancesInput{})
		if err != nil {
			fmt.Println("Error describing DB instances:", err)
			continue
		}

		// Iterate through each RDS instance and count the number of open connections
		for _, instance := range instances.DBInstances {
			unusedRDSInstances, err := utils.CountDBConnections(context.Background(), r, *instance.DBInstanceIdentifier)
			if err != nil {
				fmt.Printf("Error counting connections for DB instance %q: %v\n", *instance.DBInstanceIdentifier, err)
				continue
			}
			// Print out the list of unused RDS
			for _, rds := range unusedRDSInstances {
				fmt.Printf("Unused RDS found in %s: %v. Attention: Just because RDS has 0 connections does not mean it is not being used! Consult with your colleagues and double-check everything before deletion. It would also be useful to take a snapshot.\n", r, rds)
			}
		}

		// Create an Lambda client for the current region
		lambdaSvc := lambda.NewFromConfig(cfg, func(options *lambda.Options) {
			options.Region = r
		})

		// Call ListFunctions API to retrieve all Lambdas in the region
		lambdas, err := lambdaSvc.ListFunctions(context.Background(), &lambda.ListFunctionsInput{})
		if err != nil {
			panic(fmt.Errorf("failed to list functions, %v", err))
		}

		// Iterate through each Lambda functions and count the number of invocations
		for _, lambda := range lambdas.Functions {
			unusedLambdas, err := utils.FindTotalCountOfLambdaInvocs(context.Background(), r, *lambda.FunctionName)
			if err != nil {
				fmt.Printf("Error counting invocations for Lambda %q: %v\n", *lambda.FunctionName, err)
				continue
			}
			// Print out the list of unused Lamda Functions
			for _, lambda := range unusedLambdas {
				fmt.Printf("Unused Lamdas found in %s: %v. Attention: Just because Lambda has 0 invocations for last week does not mean it is not being used! Consult with your colleagues and double-check everything before deletion.\n", r, lambda)
			}
		}

	}
}
