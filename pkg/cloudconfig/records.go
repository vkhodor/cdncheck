package cloudconfig

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
)

var (
	DefaultRecords = &route53.Change{
		ResourceRecordSet: &route53.ResourceRecordSet{
			ResourceRecords: []*route53.ResourceRecord{
				{
					Value: aws.String("199.115.113.118"),
				},
				{
					Value: aws.String("199.115.113.105"),
				},
			},
			TTL:  aws.Int64(60),
			Type: aws.String("A"),
			GeoLocation: &route53.GeoLocation{
				CountryCode: aws.String("*"),
			},
			SetIdentifier: aws.String("default-content"),
		},
	}

	NorthAmericaRecords = &route53.Change{
		ResourceRecordSet: &route53.ResourceRecordSet{
			ResourceRecords: []*route53.ResourceRecord{
				{
					Value: aws.String("199.115.113.118"),
				},
				{
					Value: aws.String("199.115.113.105"),
				},
			},
			TTL:  aws.Int64(60),
			Type: aws.String("A"),
			GeoLocation: &route53.GeoLocation{
				ContinentCode: aws.String("NA"),
			},
			SetIdentifier: aws.String("us-content"),
		},
	}

	JapanRecords = &route53.Change{
		ResourceRecordSet: &route53.ResourceRecordSet{
			ResourceRecords: []*route53.ResourceRecord{
				{
					Value: aws.String("23.106.248.66"),
				},
			},
			TTL:  aws.Int64(60),
			Type: aws.String("A"),
			GeoLocation: &route53.GeoLocation{
				CountryCode: aws.String("JP"),
			},
			SetIdentifier: aws.String("jp-content"),
		},
	}

	AsiaRecords = &route53.Change{
		ResourceRecordSet: &route53.ResourceRecordSet{
			ResourceRecords: []*route53.ResourceRecord{
				{
					Value: aws.String("23.106.248.66"),
				},
			},
			TTL:  aws.Int64(60),
			Type: aws.String("A"),
			GeoLocation: &route53.GeoLocation{
				ContinentCode: aws.String("AS"),
			},
			SetIdentifier: aws.String("asia-content"),
		},
	}

	EuropeRecords = &route53.Change{
		ResourceRecordSet: &route53.ResourceRecordSet{
			ResourceRecords: []*route53.ResourceRecord{
				{
					Value: aws.String("95.168.161.57"),
				},
			},
			TTL:  aws.Int64(60),
			Type: aws.String("A"),
			GeoLocation: &route53.GeoLocation{
				ContinentCode: aws.String("EU"),
			},
			SetIdentifier: aws.String("eu-content"),
		},
	}

	CloudFrontRecords = &route53.Change{
		ResourceRecordSet: &route53.ResourceRecordSet{
			ResourceRecords: []*route53.ResourceRecord{
				{
					Value: aws.String("d8kbfjpasyqym.cloudfront.net"),
				},
			},
			TTL:  aws.Int64(60),
			Type: aws.String("CNAME"),
		},
	}
)
