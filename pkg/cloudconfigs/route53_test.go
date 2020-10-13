package cloudconfigs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/vkhodor/cdncheck/pkg/config"
	"testing"
)

func TestGetState(t *testing.T) {
	records := []*route53.ResourceRecordSet{
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("A")},
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("CNAME")},
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("A")},
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("A")},
	}

	result, err := getState(records, "content.cdn.personaly.bid")
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

	result, err = getState(records, "content.cdn.personaly.bid")
	if err != nil {
		t.Error()
	}
	if result != "normal" {
		t.Error()
	}

}

func TestRecordsToChanges(t *testing.T) {

	answer := `{
  ResourceRecordSet: {
    GeoLocation: {
      CountryCode: "US"
    },
    Name: "content",
    ResourceRecords: [{
        Value: "1.1.1.1"
      },{
        Value: "2.2.2.2"
      },{
        Value: "3.3.3.3"
      }],
    SetIdentifier: "test",
    TTL: 60,
    Type: "A"
  }
}`


	var records []config.DNSRecord
	values := []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"}

	records = append(records, config.DNSRecord{
		Identifier: aws.String("test"),
		Values: &values,
		Type: aws.String("A"),
		CountryCode: aws.String("US"),
		TTL: aws.Int(60),
		Name: aws.String("content"),
	})

	changes, err := recordsToChanges(records)
	if err != nil {
		t.Error()
	}

	for _, c := range changes {
		if c.String() != answer {
			t.Error()
		}
	}
}
