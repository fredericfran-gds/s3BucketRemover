package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var name, region *string

func init() {
	fmt.Println("S3 Versioned Bucket Remover")
	flags()
}

func flags() {
	name = flag.String("name", "", "name of the s3 bucket to be removed")
	region = flag.String("region", "", "aws region of the s3 bucket to be removed")

	flag.Parse()

	if err := validateFlag(); err != nil {
		fmt.Printf("error while parsing flags: %v\n", err)
		os.Exit(1)
	}
}

func validateFlag() error {
	errStr := ""

	if *name == "" {
		errStr += ": name flag was not set"
	}

	if *region == "" {
		errStr += ": region flag was not set"
	}

	if errStr != "" {
		return fmt.Errorf("following flags were not set %s", errStr)
	}

	return nil
}

func createSession(region string) *s3.S3 {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region)}))
	svc := s3.New(sess)
	return svc
}

func deleteBucket(svc *s3.S3, bucket string) error {
	iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})

	if err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter); err != nil {
		return fmt.Errorf("unable to delete objects from bucket %s: %v", bucket, err)
	}

	_, err := svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return fmt.Errorf("unable to delete bucket %s: %v", bucket, err)
	}

	// Wait until bucket is deleted before finishing
	fmt.Printf("Waiting for bucket %s to be deleted...\n", bucket)

	err = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		return fmt.Errorf("error occurred when waiting for bucket %s to be deleted: %v", bucket, err)
	}

	return nil
}

func main() {
	fmt.Printf("bucket %v to be deleted in region %v\n", *name, *region)

	svc := createSession(*region)

	err := deleteBucket(svc, *name)
	if err != nil {
		fmt.Printf("failed to delete bucket %s: %v\n", *name, err)
		os.Exit(1)
	}

	fmt.Printf("Bucket %s successfully deleted\n", *name)
}
