package cloudconfigs

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/sirupsen/logrus"
	"github.com/vkhodor/cdncheck/pkg/config"
)

type CloudRoute53PolicyBased struct {
	client              *route53.Route53
	zoneId              string
	recordName          string
	logger              *logrus.Logger
	normalPolicyBased   *config.DNSRecord
	fallbackPolicyBased *config.DNSRecord
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
	policyInstances, err := c.ListPolicyInstances()
	if err != nil {
		return false, err
	}

	for _, p := range policyInstances.TrafficPolicyInstances {
		if *p.TrafficPolicyId == *c.fallbackPolicyBased.TrafficPolicyId {
			deleteInput := &route53.DeleteTrafficPolicyInstanceInput{Id: p.Id}
			deleteOutput, err := c.client.DeleteTrafficPolicyInstance(deleteInput)
			c.logger.Debug(deleteOutput)
			if err != nil {
				return false, err
			}
			break
		}
	}

	// если таки удаляли фолбэк то в цикле получаем Instance по Id до тех пор пока попытка получаения успешная.
	// как только запись получить не удалось добавляем новую

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

func (c *CloudRoute53PolicyBased) ListPolicyInstances() (*route53.ListTrafficPolicyInstancesByHostedZoneOutput, error) {
	input := &route53.ListTrafficPolicyInstancesByHostedZoneInput{
		HostedZoneId: aws.String(c.zoneId),
	}
	return c.client.ListTrafficPolicyInstancesByHostedZone(input)
}

func (c *CloudRoute53PolicyBased) State() (string, error) {
	output, err := c.ListPolicyInstances()
	if err != nil {
		return "unknown:Unknown", err
	}

	for _, p := range output.TrafficPolicyInstances {
		if *p.TrafficPolicyId == *c.normalPolicyBased.TrafficPolicyId && *p.TrafficPolicyVersion == *c.normalPolicyBased.TrafficPolicyVersion {
			return fmt.Sprintf("normal:%v", *p.State), nil
		}
		if *p.TrafficPolicyId == *c.fallbackPolicyBased.TrafficPolicyId && *p.TrafficPolicyVersion == *c.fallbackPolicyBased.TrafficPolicyVersion {
			return fmt.Sprintf("fallback:%v", *p.State), nil
		}
	}

	return "unknown:Unknown", nil
}

func (c *CloudRoute53PolicyBased) LoadChanges(config config.Config) error {
	fallbackPolicyBased, err := config.GetPolicyBasedFallbackRecord()
	if err != nil {
		return err
	}
	c.fallbackPolicyBased = &fallbackPolicyBased

	normalPolicyBased, err := config.GetPolicyBasedNormalRecord()
	if err != nil {
		return err
	}
	c.normalPolicyBased = &normalPolicyBased
	return nil
}
