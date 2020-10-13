package cloudconfigs

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/sirupsen/logrus"
	"github.com/vkhodor/cdncheck/pkg/config"
)

type CloudRoute53 struct {
	client          *route53.Route53
	zoneId          string
	recordName      string
	logger          *logrus.Logger
	normalRecords   []*route53.Change
	fallbackRecords []*route53.Change
}

func NewCloudRoute53(zoneId string, recordName string, logger *logrus.Logger) *CloudRoute53 {
	mySession := session.Must(session.NewSession())
	svc := route53.New(mySession)
	ret := CloudRoute53{client: svc, zoneId: zoneId, recordName: recordName, logger: logger}
	return &ret
}

func (c *CloudRoute53) State() (string, error) {
	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(c.zoneId),
	}

	output, err := c.client.ListResourceRecordSets(input)
	if err != nil {
		c.logger.Debug(err)
		return "error", err
	}
	return getState(output.ResourceRecordSets, c.recordName)
}

func (c *CloudRoute53) Fallback() (bool, error) {
	result, err := c.makeChanges("DELETE")
	if err != nil {
		return false, err
	}
	c.logger.Debug("Fallback(): ", result, err)
	return true, nil
}

func (c *CloudRoute53) Normal() (bool, error) {
	result, err := c.makeChanges("CREATE")
	if err != nil {
		return false, err
	}
	c.logger.Debug("Normal(): ", result, err)
	return true, nil
}

func (c *CloudRoute53) makeChanges(changesType string) (bool, error) {
	action := aws.String(changesType)
	name := aws.String(c.recordName)

	NorthAmericaRecords.Action = action
	AsiaRecords.Action = action
	JapanRecords.Action = action
	EuropeRecords.Action = action
	DefaultRecords.Action = action
	CloudFrontRecords.Action = aws.String("CREATE")

	NorthAmericaRecords.ResourceRecordSet.Name = name
	AsiaRecords.ResourceRecordSet.Name = name
	JapanRecords.ResourceRecordSet.Name = name
	EuropeRecords.ResourceRecordSet.Name = name
	DefaultRecords.ResourceRecordSet.Name = name
	CloudFrontRecords.ResourceRecordSet.Name = name

	batch := []*route53.Change{
		NorthAmericaRecords,
		AsiaRecords,
		JapanRecords,
		EuropeRecords,
		DefaultRecords,
	}

	if changesType == "CREATE" {
		CloudFrontRecords.Action = aws.String("DELETE")
		batch = append([]*route53.Change{CloudFrontRecords}, batch...)
	} else {
		batch = append(batch, CloudFrontRecords)
	}

	var input = &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: batch,
		},
		HostedZoneId: aws.String(c.zoneId),
	}
	_, err := c.client.ChangeResourceRecordSets(input)
	if err != nil {
		return false, err
	}

	return true, nil
}

func recordsToChanges(records []config.DNSRecord) ([]*route53.Change, error) {
	var changes []*route53.Change
	for _, r := range records {
		var values []*route53.ResourceRecord
		for _, v := range *r.Values {
			values = append(values, &route53.ResourceRecord{
				Value: aws.String(v),
			})

		}

		changes = append(changes, &route53.Change{
			ResourceRecordSet: &route53.ResourceRecordSet{
				Name:            r.Name,
				ResourceRecords: values,
				TTL:             aws.Int64(int64(*r.TTL)),
				Type:            r.Type,
				GeoLocation: &route53.GeoLocation{
					CountryCode: r.CountryCode,
				},
				SetIdentifier: r.Identifier,
			},
		})
	}

	return changes, nil
}

func (c *CloudRoute53) LoadRecords(config config.Config) error {
	normalRecords, err := config.GetNormalRecords()
	if err != nil {
		return err
	}

	fallbackRecords, err := config.GetFallbackRecords()
	if err != nil {
		return err
	}

	c.normalRecords, err = recordsToChanges(normalRecords)
	if err != nil {
		return err
	}

	c.fallbackRecords, err = recordsToChanges(fallbackRecords)

	return nil
}

func getState(records []*route53.ResourceRecordSet, recordName string) (string, error) {
	if len(records) == 0 {
		return "error", errors.New("len of found records is null")
	}
	for _, record := range records {
		if *record.Name == recordName && *record.Type == "CNAME" {
			return "fallback", nil
		}
	}
	return "normal", nil
}
