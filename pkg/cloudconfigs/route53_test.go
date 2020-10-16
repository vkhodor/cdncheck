package cloudconfigs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/sirupsen/logrus"
	"github.com/vkhodor/cdncheck/pkg/config"
	"testing"
)

func TestGetState(t *testing.T) {

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{DisableColors: false, FullTimestamp: true})
	logger.SetLevel(logrus.DebugLevel)

	records := []*route53.ResourceRecordSet{
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("A"), SetIdentifier: aws.String("fallback:xxx")},
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("CNAME"), SetIdentifier: aws.String("fallback:yyy")},
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("A"), SetIdentifier: aws.String("fallback:zzz")},
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("A"), SetIdentifier: aws.String("fallback:aaaa")},
	}

	result, err := getState(records, logger)
	if err != nil {
		t.Error()
	}
	if result != "fallback" {
		t.Error()
	}

	records = []*route53.ResourceRecordSet{
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("A"), SetIdentifier: aws.String("normal:aaaa")},
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("A"), SetIdentifier: aws.String("normal:bbbb")},
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("A"), SetIdentifier: aws.String("normal:cccc")},
	}

	result, err = getState(records, logger)
	if err != nil {
		t.Error()
	}
	if result != "normal" {
		t.Error()
	}

	records = []*route53.ResourceRecordSet{
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("A"), SetIdentifier: aws.String("normal:aaaaa")},
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("A"), SetIdentifier: aws.String("fallback:b")},
		&route53.ResourceRecordSet{Name: aws.String("content.cdn.personaly.bid"), Type: aws.String("A"), SetIdentifier: aws.String("normal:default-content")},
	}

	result, err = getState(records, logger)
	if err == nil {
		t.Error()
	}
	if result != "error" {
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
		Identifier:  aws.String("test"),
		Values:      &values,
		Type:        aws.String("A"),
		CountryCode: aws.String("US"),
		TTL:         aws.Int(60),
		Name:        aws.String("content"),
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
