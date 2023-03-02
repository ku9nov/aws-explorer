# aws-explorer

The application is designed to scan an AWS infrastructure and identify unused resources. It has the ability to scan for various resources, such as RDS instances with 0 connections, Lambda functions that haven't executed in over a week, load balancers without listeners, target groups without targets (and it can recognize if a load balancer is connected to it), unused Elastic IPs, unused EBS volumes, and unused network interfaces.

The scanning process is performed by leveraging the AWS SDK and making API calls to AWS services. The application utilizes filters and queries to identify resources that meet the specified criteria for each resource type.

The application can identify unused resources across all regions in one pass, saving the user time and effort. This feature can be particularly useful for large and complex infrastructures with resources distributed across multiple regions. The user can then choose to remove the identified unused resources to optimize their infrastructure and reduce costs.

:warning: Once the scanning is complete, the application will display a report of the identified unused resources. It's important to carefully review the identified resources before removing them to avoid accidentally deleting important resources.

### DynamoDB
The application also includes a feature to scan DynamoDB tables in the AWS account that have the word "terraform" in their name. If a table with this name is found, the application will check if it's in PROVISIONED capacity mode rather than ON-DEMAND. Users can consider switching to On-demand capacity mode if their usage patterns allow it, which can result in significant cost savings.

:warning: Exploring and controlling additional costs in DynamoDB tables works only if your DynamoDB table has at least once been set to on-demand capacity mode.
## Getting Started

These instructions will help you to run this application on your local machine.

### Prerequisites

- Go version 1.18 or above
- A program that you want to run in the background

### Installing

- Clone the repository
- Run `go build` command to build the application

### Running the application

```
./aws-explorer
```
When running the application with this command, it will use the default settings from the system.

```
AWS_PROFILE=OTHER_PROFILE_NAME ./aws-explorer
```
If you have multiple AWS profiles on your system and would like to use a different profile than the default one, you can use this command. Replace "OTHER_PROFILE_NAME" with the name of the profile you want to use.

```
LOG_LEVEL=DEBUG ./aws-explorer
```
If you would like to enable debug mode, you can use this command. The application will output more information about its processes and will be more verbose. This can be helpful for troubleshooting issues.