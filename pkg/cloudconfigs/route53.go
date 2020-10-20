package cloudconfigs

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/sirupsen/logrus"
	"github.com/vkhodor/cdncheck/pkg/config"
	"strings"
)

type CloudRoute53 struct {
	client          *route53.Route53
	zoneId          string
	recordName      string
	logger          *logrus.Logger
	normalChanges   []*route53.Change
	fallbackChanges []*route53.Change
}

func NewCloudRoute53(zoneId string, recordName string, logger *logrus.Logger) *CloudRoute53 {
	mySession := session.Must(session.NewSession())
	svc := route53.New(mySession)
	ret := CloudRoute53{
		client:     svc,
		zoneId:     zoneId,
		recordName: recordName,
		logger:     logger,
	}
	return &ret
}

func (c *CloudRoute53) State() (string, error) {
	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(c.zoneId),
		StartRecordName: &c.recordName,
	}

	output, err := c.client.ListResourceRecordSets(input)
	if err != nil {
		c.logger.Debug(err)
		return "error", err
	}
	return getState(output.ResourceRecordSets, c.logger)
}

func (c *CloudRoute53) setNormalAction(action string) {
	for _, r := range c.normalChanges {
		r.Action = aws.String(action)
	}
}

func (c *CloudRoute53) setFallbackAction(action string) {
	for _, r := range c.fallbackChanges {
		r.Action = aws.String(action)
	}
}

func (c *CloudRoute53) Normal() (bool, error) {
	c.setFallbackAction("DELETE")
	c.setNormalAction("CREATE")

	var batch []*route53.Change
	batch = c.fallbackChanges
	for _, r := range c.normalChanges {
		batch = append(batch, r)
	}
	return c.makeChanges(batch)
}

func (c *CloudRoute53) Fallback() (bool, error) {
	c.setNormalAction("DELETE")
	c.setFallbackAction("CREATE")

	var batch []*route53.Change
	batch = c.normalChanges
	for _, r := range c.fallbackChanges {
		batch = append(batch, r)
	}

	return c.makeChanges(batch)
}

func (c *CloudRoute53) makeChanges(batch []*route53.Change) (bool, error) {
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
					CountryCode:   r.CountryCode,
					ContinentCode: r.ContinentCode,
				},
				SetIdentifier: r.Identifier,
			},
		})
	}

	return changes, nil
}

func (c *CloudRoute53) LoadChanges(config config.Config) error {
	normalRecords, err := config.GetNormalRecords()
	if err != nil {
		return err
	}

	fallbackRecords, err := config.GetFallbackRecords()
	if err != nil {
		return err
	}

	c.normalChanges, err = recordsToChanges(normalRecords)
	if err != nil {
		return err
	}

	c.fallbackChanges, err = recordsToChanges(fallbackRecords)

	return nil
}

func getState(records []*route53.ResourceRecordSet, logger *logrus.Logger) (string, error) {
	state := ""
	if len(records) == 0 {
		error := errors.New("len of found records is null")
		logger.Debug("error ", error)
		return "error", error
	}
	for i, record := range records {
		recordState := strings.Split(*record.SetIdentifier, ":")[0]
		if i == 0 {
			state = recordState
			continue
		}

		if state != recordState {
			e := errors.New("mismatch in record identifiers")
			logger.Debug("error ", e)
			return "error", e
		}
	}
	logger.Debug(state)
	return state, nil
}
