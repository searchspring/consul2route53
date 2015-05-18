package testutil

import(
	"time"
	"fmt"
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/route53"
)

type MockRoute53 struct{
	Fail bool
}

func ( r *MockRoute53 ) ListResourceRecordSets(*route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error) {
	if r.Fail {
		return &route53.ListResourceRecordSetsOutput{},fmt.Errorf("Bogus failure")
	}
	rec1 := &route53.ResourceRecord{ Value: aws.String("127.0.0.1")}
	var ttl int64
	ttl = 3600
	recset1 := &route53.ResourceRecordSet{
		AliasTarget: &route53.AliasTarget{},
		Failover: aws.String("failover"),
		GeoLocation: &route53.GeoLocation{},
		HealthCheckID: aws.String("bogus"),
		Name: aws.String("bogus"),
		Region: aws.String("us-east-1"),
		ResourceRecords: []*route53.ResourceRecord{rec1},
		SetIdentifier: aws.String("bogus"),
		TTL: &ttl,
		Type: aws.String("A"),
		Weight:  nil,

	}
	rec2 := &route53.ResourceRecord{ Value: aws.String("127.0.0.2")}
	recset2 := &route53.ResourceRecordSet{
		AliasTarget: &route53.AliasTarget{},
		Failover: aws.String("failover"),
		GeoLocation: &route53.GeoLocation{},
		HealthCheckID: aws.String("bogus"),
		Name: aws.String("redis.bogus.com."),
		Region: aws.String("us-east-1"),
		ResourceRecords: []*route53.ResourceRecord{rec2},
		SetIdentifier: aws.String("redis.bogus.com."),
		TTL: &ttl,
		Type: aws.String("A"),
		Weight:  nil,

	}
	istruncated := false
	ret := &route53.ListResourceRecordSetsOutput{
		IsTruncated: &istruncated,
		MaxItems: aws.String("100"),
		ResourceRecordSets: []*route53.ResourceRecordSet{recset1,recset2},
	}
	return ret, nil
}

func ( r *MockRoute53 ) ChangeResourceRecordSets(*route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput,error) {
	if r.Fail {
		return &route53.ChangeResourceRecordSetsOutput{},fmt.Errorf("Bogus failure")
	}
	now := time.Now()
	changeinfo := &route53.ChangeInfo{
		Comment: aws.String("Nothing"),
		ID: aws.String("bogus"),
		Status: aws.String("INSYC"),
		SubmittedAt: &now,
	}
	ret := &route53.ChangeResourceRecordSetsOutput{
		ChangeInfo: changeinfo,
	}
	return ret, nil
}

func ( r *MockRoute53 ) GetHostedZone(*route53.GetHostedZoneInput) (*route53.GetHostedZoneOutput,error) {
	if r.Fail {
		return &route53.GetHostedZoneOutput{},fmt.Errorf("Bogus failure")
	}
	ret := &route53.GetHostedZoneOutput{ HostedZone: &route53.HostedZone{ Name: aws.String("bogus.com.")}}
	return ret, nil
}