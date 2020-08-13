package cloudconfigs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"testing"
)

func TestGetStatus(t *testing.T) {
	records := []*route53.ResourceRecordSet{
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("A")},
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("CNAME")},
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("A")},
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("A")},
	}

	result, err := getStatus(records, "content.cdn.personaly.bid")
	if err != nil {
		t.Error()
	}
	if result != "fallback" {
		t.Error()
	}

	records = []*route53.ResourceRecordSet{
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("A")},
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("A")},
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("A")},
	}

	result, err = getStatus(records, "content.cdn.personaly.bid")
	if err != nil {
		t.Error()
	}
	if result != "normal" {
		t.Error()
	}

}
