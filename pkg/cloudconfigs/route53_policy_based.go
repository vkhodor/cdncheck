package cloudconfigs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/sirupsen/logrus"
	"github.com/vkhodor/cdncheck/pkg/config"
)

type CloudRoute53PolicyBased struct {
	client          *route53.Route53
	zoneId          string
	recordName      string
	logger          *logrus.Logger
	normalChanges   []*route53.Change
	fallbackChanges []*route53.Change
}

func NewCloudRoute53PolicyBased(zoneId string, recordName string, logger *logrus.Logger) *CloudRoute53PolicyBased {
	mySession := session.Must(session.NewSession())
	svc := route53.New(mySession)
	ret := CloudRoute53PolicyBased{
		client:     svc,
		zoneId:     zoneId,
		recordName: recordName,
		logger:     logger,
	}
	return &ret
}

func (c *CloudRoute53PolicyBased) Normal() (bool, error) {
	input := &route53.CreateTrafficPolicyInstanceInput{
		HostedZoneId:         aws.String(c.zoneId),
		TrafficPolicyId:      aws.String("2d912bec-8e77-44de-8183-070f5227b302"),
		TrafficPolicyVersion: aws.Int64(1),
		Name:                 aws.String("test.algorithmic.bid."),
		TTL:                  aws.Int64(60),
	}
	result, err := c.client.CreateTrafficPolicyInstance(input)
	c.logger.Debug(result)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *CloudRoute53PolicyBased) Fallback() (bool, error) {
	input := &route53.CreateTrafficPolicyInstanceInput{
		HostedZoneId:         aws.String(c.zoneId),
		TrafficPolicyId:      aws.String("ca7be789-cd5f-45af-9e86-310091df7f93"),
		TrafficPolicyVersion: aws.Int64(1),
	}
	result, err := c.client.CreateTrafficPolicyInstance(input)
	c.logger.Debug(result)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *CloudRoute53PolicyBased) State() (string, error) {
	input := &route53.ListTrafficPolicyInstancesByHostedZoneInput{
		HostedZoneId: aws.String(c.zoneId),
	}

	output, err := c.client.ListTrafficPolicyInstancesByHostedZone(input)
	c.logger.Debug(output)
	if err != nil {
		c.logger.Debug(err)
		return "error", err
	}
	return "true", nil
}

func (c *CloudRoute53PolicyBased) LoadChanges(config.Config) error {
	return nil
}
