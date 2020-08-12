package cloudconfigs

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/sirupsen/logrus"
)

type CloudRoute53 struct {
	client     *route53.Route53
	zoneId     string
	recordName string
	logger     *logrus.Logger
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
	return getStatus(output.ResourceRecordSets, c.recordName)
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

func getStatus(records []*route53.ResourceRecordSet, recordName string) (string, error) {
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
