package cloudconfigs

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/sirupsen/logrus"
	"github.com/vkhodor/cdncheck/pkg/config"
	"time"
)

type CloudRoute53PolicyBased struct {
	client              *route53.Route53
	zoneId              string
	recordName          string
	logger              *logrus.Logger
	normalPolicyBased   *config.DNSRecord
	fallbackPolicyBased *config.DNSRecord
	fallbackPrefix      *string
	normalPrefix        *string
}

func NewCloudRoute53PolicyBased(zoneId string, recordName string, logger *logrus.Logger, fallbackPrefix *string, normalPrefix *string) *CloudRoute53PolicyBased {
	mySession := session.Must(session.NewSession())
	svc := route53.New(mySession)
	ret := CloudRoute53PolicyBased{
		client:         svc,
		zoneId:         zoneId,
		recordName:     recordName,
		logger:         logger,
		fallbackPrefix: fallbackPrefix,
		normalPrefix:   normalPrefix,
	}
	return &ret
}

func (c *CloudRoute53PolicyBased) Normal() (bool, error) {
	policyInstances, err := c.ListPolicyInstances()
	if err != nil {
		return false, err
	}

	var id *string
	for _, p := range policyInstances.TrafficPolicyInstances {
		if *p.TrafficPolicyId == *c.fallbackPolicyBased.TrafficPolicyId {
			id = p.Id
			deleteInput := &route53.DeleteTrafficPolicyInstanceInput{Id: p.Id}
			deleteOutput, err := c.client.DeleteTrafficPolicyInstance(deleteInput)
			c.logger.Debug(deleteOutput)
			if err != nil {
				return false, err
			}
			break
		}
	}

	if id != nil {
		var t time.Duration = 0
		policyRemoved := false
		for t < 5*time.Minute {
			time.Sleep(10 * time.Second)
			t += 10 * time.Second
			c.logger.Debug("Sleeping: ", t)
			_, err := c.client.GetTrafficPolicyInstance(&route53.GetTrafficPolicyInstanceInput{Id: id})
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case route53.ErrCodeNoSuchTrafficPolicyInstance:
						policyRemoved = true
					default:
						c.logger.Error("GetTrafficPolicyInstance: ", err)
						continue
					}
					if policyRemoved {
						c.logger.Debug("Policy removed. Going next.")
						break
					}

				}
			}
			c.logger.Debug("Found relevant policy instance. Continue")
			continue
		}
		if t >= 5*time.Minute {
			c.logger.Error("Timeout for policy removing.")
			return false, nil
		}
	}

	input := &route53.CreateTrafficPolicyInstanceInput{
		HostedZoneId:         aws.String(c.zoneId),
		TrafficPolicyId:      c.normalPolicyBased.TrafficPolicyId,
		TrafficPolicyVersion: c.normalPolicyBased.TrafficPolicyVersion,
		Name:                 aws.String(c.recordName),
		TTL:                  aws.Int64(int64(*c.normalPolicyBased.TTL)),
	}
	result, err := c.client.CreateTrafficPolicyInstance(input)
	c.logger.Debug(result)
	if err != nil {
		return false, err
	}

	var t time.Duration = 0
	for t < 5*time.Minute {
		time.Sleep(10 * time.Second)
		t += 10 * time.Second
		c.logger.Debug("Sleeping: ", t)
		state, _ := c.State()
		if state == "normal:Applied" {
			c.logger.Debug("done")
			break
		}
		c.logger.Debug("Policy is not applied yet. ", state)
	}
	return true, nil
}

func (c *CloudRoute53PolicyBased) Fallback() (bool, error) {
	policyInstances, err := c.ListPolicyInstances()
	if err != nil {
		return false, err
	}

	var id *string
	for _, p := range policyInstances.TrafficPolicyInstances {
		if *p.TrafficPolicyId == *c.normalPolicyBased.TrafficPolicyId {
			id = p.Id
			deleteInput := &route53.DeleteTrafficPolicyInstanceInput{Id: p.Id}
			deleteOutput, err := c.client.DeleteTrafficPolicyInstance(deleteInput)
			c.logger.Debug(deleteOutput)
			if err != nil {
				return false, err
			}
			break
		}
	}

	if id != nil {
		var t time.Duration = 0
		policyRemoved := false
		for t < 5*time.Minute {
			time.Sleep(10 * time.Second)
			t += 10 * time.Second
			c.logger.Debug("Sleeping: ", t)
			_, err := c.client.GetTrafficPolicyInstance(&route53.GetTrafficPolicyInstanceInput{Id: id})
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case route53.ErrCodeNoSuchTrafficPolicyInstance:
						policyRemoved = true
					default:
						c.logger.Error("GetTrafficPolicyInstance: ", err)
						continue
					}
					if policyRemoved {
						c.logger.Debug("Policy removed. Going next.")
						break
					}
				}
			}
			c.logger.Debug("Found relevant policy instance. Continue")
			continue
		}
		if t >= 5*time.Minute {
			c.logger.Error("Timeout for policy removing.")
			return false, nil
		}
	}

	input := &route53.CreateTrafficPolicyInstanceInput{
		HostedZoneId:         aws.String(c.zoneId),
		TrafficPolicyId:      c.fallbackPolicyBased.TrafficPolicyId,
		TrafficPolicyVersion: c.fallbackPolicyBased.TrafficPolicyVersion,
		Name:                 aws.String(c.recordName),
		TTL:                  aws.Int64(int64(*c.fallbackPolicyBased.TTL)),
	}
	result, err := c.client.CreateTrafficPolicyInstance(input)
	c.logger.Debug(result)
	if err != nil {
		return false, err
	}

	var t time.Duration = 0
	for t < 5*time.Minute {
		time.Sleep(10 * time.Second)
		t += 10 * time.Second
		c.logger.Debug("Sleeping: ", t)
		state, _ := c.State()
		if state == "fallback:Applied" {
			c.logger.Debug("done")
			break
		}
		c.logger.Debug("Policy is not applied yet. ", state)
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
			return fmt.Sprintf("%v:%v", *c.normalPrefix, *p.State), nil
		}
		if *p.TrafficPolicyId == *c.fallbackPolicyBased.TrafficPolicyId && *p.TrafficPolicyVersion == *c.fallbackPolicyBased.TrafficPolicyVersion {
			return fmt.Sprintf("%v:%v", *c.fallbackPrefix, *p.State), nil
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
