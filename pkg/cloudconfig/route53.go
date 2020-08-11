package cloudconfig

import (
	"github.com/aws/aws-sdk-go/aws"
)

//	"github.com/aws/aws-sdk-go/aws/session"
//	"github.com/aws/aws-sdk-go/service/athena"
//	"github.com/dustin/go-humanize"

type CloudRoute53 struct {
	client *Route53
}

func NewRoute53() *CloudRoute53 {
	mySession := session.Must(session.NewSession())
	svc := route53.New(mySession)
	ret := CloudRoute53{client: svc}
	return &ret
}
