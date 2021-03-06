package ebs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func Attach(svc *ec2.EC2, instanceID string, volumeID string) {
	input := &ec2.AttachVolumeInput{
		Device:     aws.String("/dev/sdf"),
		InstanceId: aws.String(instanceID),
		VolumeId:   aws.String(volumeID),
	}
	_, err := svc.AttachVolume(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		} else {
			fmt.Println(err.Error())
		}
		return
	}
	describeVolumesInput := &ec2.DescribeVolumesInput{
		VolumeIds: aws.StringSlice([]string{volumeID}),
	}
	if err := svc.WaitUntilVolumeInUse(describeVolumesInput); err != nil {
		panic(err)
	}
}

func Detach(svc *ec2.EC2, volumeID string) {
	input := &ec2.DetachVolumeInput{
		VolumeId: aws.String(volumeID),
	}
	_, err := svc.DetachVolume(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		} else {
			fmt.Println(err.Error())
		}
		return
	}
	describeVolumesInput := &ec2.DescribeVolumesInput{
		VolumeIds: aws.StringSlice([]string{volumeID}),
	}
	if err := svc.WaitUntilVolumeAvailable(describeVolumesInput); err != nil {
		panic(err)
	}
}

func Create(svc *ec2.EC2, snapshotID string, name string) {
	tagList := &ec2.TagSpecification{
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(name),
			},
		},
		ResourceType: aws.String(ec2.ResourceTypeVolume),
	}
	input := &ec2.CreateVolumeInput{
		AvailabilityZone:  aws.String("us-east-1a"),
		SnapshotId:        aws.String(snapshotID),
		VolumeType:        aws.String("sc1"),
		TagSpecifications: []*ec2.TagSpecification{tagList},
	}
	_, err := svc.CreateVolume(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		} else {
			fmt.Println(err.Error())
		}
		return
	}
	describeVolumesInput := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(name)},
			},
		},
	}
	if err := svc.WaitUntilVolumeAvailable(describeVolumesInput); err != nil {
		panic(err)
	}
}

func Delete(svc *ec2.EC2, volumeID string) {
	input := &ec2.DeleteVolumeInput{
		VolumeId: aws.String(volumeID),
	}
	_, err := svc.DeleteVolume(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		} else {
			fmt.Println(err.Error())
		}
		return
	}
	describeVolumesInput := &ec2.DescribeVolumesInput{
		VolumeIds: aws.StringSlice([]string{volumeID}),
	}
	if err := svc.WaitUntilVolumeDeleted(describeVolumesInput); err != nil {
		panic(err)
	}
}

func FetchVolumeID(svc *ec2.EC2, name string) string {
	var volumeID string
	params := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(name)},
			},
		},
	}
	result, err := svc.DescribeVolumes(params)
	if err != nil {
		fmt.Println("Error", err)
	}
	for _, volumes := range result.Volumes {
		volumeID = aws.StringValue(volumes.VolumeId)
	}
	return volumeID
}
