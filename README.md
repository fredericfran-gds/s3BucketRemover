# Basic S3 Bucket Remover

This is a utility to delete a S3 bucket because there is no direct AWS CLI command
to delete a versioned S3 bucket

## Usage

### Prerequisites

1. Assume that you have either:
  1. exported your AWS credentials or
  2. set your AWS credentials file `~/.aws/credentials` or
  3. configured aws vault

2. have Go installed, tested with version 1.11

### Run

1. execute:
   ```
   go run main.go --name <s3_bucket_name> --region <aws_region>
   ```
   where `<s3_bucket_name>` is the bucket name to be deleted and `<aws_region>` is the name of the AWS region where the S3 bucket is located

