package utils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

func FindUnusedElasticIps(ctx context.Context, svc *ec2.Client) ([]string, error) {
	input := &ec2.DescribeAddressesInput{}
	output, err := svc.DescribeAddresses(ctx, input)
	if err != nil {
		return nil, err
	}

	var unusedIps []string
	for _, a := range output.Addresses {
		if aws.ToString(a.InstanceId) == "" && aws.ToString(a.NetworkInterfaceId) == "" {
			unusedIps = append(unusedIps, aws.ToString(a.PublicIp))
		}
	}
	return unusedIps, nil
}

func FindUnusedEBSVolumes(ctx context.Context, svc *ec2.Client) ([]string, error) {
	input := &ec2.DescribeVolumesInput{}
	output, err := svc.DescribeVolumes(ctx, input)
	if err != nil {
		return nil, err
	}

	var unusedVolumes []string
	for _, v := range output.Volumes {
		if len(v.Attachments) == 0 {
			unusedVolumes = append(unusedVolumes, aws.ToString(v.VolumeId))
		}
	}
	return unusedVolumes, nil
}

func FindUnusedNetworkInterfaces(ctx context.Context, svc *ec2.Client) ([]string, error) {
	input := &ec2.DescribeNetworkInterfacesInput{}
	output, err := svc.DescribeNetworkInterfaces(ctx, input)
	if err != nil {
		return nil, err
	}

	var unusedInterfaces []string
	for _, n := range output.NetworkInterfaces {
		if n.Attachment == nil || len(aws.ToString(n.Attachment.AttachmentId)) == 0 {
			unusedInterfaces = append(unusedInterfaces, aws.ToString(n.NetworkInterfaceId))
		}
	}
	return unusedInterfaces, nil
}

func FindUnusedTargetGroups(ctx context.Context, svc *elbv2.Client) ([]string, error) {
	input := &elbv2.DescribeTargetGroupsInput{}
	output, err := svc.DescribeTargetGroups(ctx, input)
	if err != nil {
		return nil, err
	}

	var unusedTargetGroups []string
	for _, tg := range output.TargetGroups {
		if len(tg.LoadBalancerArns) == 0 {
			unusedTargetGroups = append(unusedTargetGroups, aws.ToString(tg.TargetGroupArn))
			continue
		}

		// Check if the target group has any targets
		targets, err := svc.DescribeTargetHealth(ctx, &elbv2.DescribeTargetHealthInput{
			TargetGroupArn: tg.TargetGroupArn,
		})
		if err != nil {
			return nil, err
		}

		if len(targets.TargetHealthDescriptions) == 0 {
			unusedTargetGroups = append(unusedTargetGroups, fmt.Sprintf("%s. Has attached unused LoadBalancer ARN: %s", aws.ToString(tg.TargetGroupArn), tg.LoadBalancerArns[0]))
		}
	}
	return unusedTargetGroups, nil
}

func FindUnusedLoadBalancers(ctx context.Context, svc *elbv2.Client) ([]string, error) {
	input := &elbv2.DescribeLoadBalancersInput{}
	output, err := svc.DescribeLoadBalancers(ctx, input)
	if err != nil {
		return nil, err
	}

	var unusedLBs []string
	for _, lb := range output.LoadBalancers {
		if len(*lb.LoadBalancerArn) > 0 {
			// Check if the load balancer has any listeners
			listenersInput := &elbv2.DescribeListenersInput{
				LoadBalancerArn: lb.LoadBalancerArn,
			}
			listenersOutput, err := svc.DescribeListeners(ctx, listenersInput)
			if err != nil {
				return nil, err
			}
			if len(listenersOutput.Listeners) == 0 {
				unusedLBs = append(unusedLBs, aws.ToString(lb.LoadBalancerArn))
			}
		}
	}
	return unusedLBs, nil
}
