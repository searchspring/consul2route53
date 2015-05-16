 package consul2route53

import(
	"fmt"
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/route53"
)

type Route53Srv struct {
	*Config
	records []*Record
	Srv SrvRoute53
	Changes []*route53.Change
}

type SrvRoute53 interface {
	ListResourceRecordSets(*route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput,error)
	ChangeResourceRecordSets(*route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput,error)
	GetHostedZone(*route53.GetHostedZoneInput)(*route53.GetHostedZoneOutput,error)
}

func (r *Route53Srv) ChangesNum() int {
	return len(r.Changes)
}

func (r *Route53Srv) Records() []*Record {
	return r.records
}

func (r *Route53Srv) RecordsMap() map[string]*Record {
	recmap := make(map[string]*Record)
	for _,record := range r.records {
		recmap[record.Name+"_"+record.Type] = record
	}
	return recmap
}

func (r *Route53Srv) GetZoneInfo() error {
	if r.Srv == nil {
		r.Srv = route53.New(nil)
	}
	params := &route53.GetHostedZoneInput{ID: aws.String(r.ZoneId())}
	resp, err := r.Srv.GetHostedZone(params)
	if awserr := aws.Error(err); awserr != nil {
	    // A service error occurred.
	    return fmt.Errorf("Service Error: %#s : %#s\n", awserr.Code, awserr.Message)
	} else if err != nil {
	    // A non-service error occurred.
	    return fmt.Errorf("Non-Service Error: %#s\n",err)
	}
	r.SetZone(*resp.HostedZone.Name)
	return nil
}


// GetRecords retrieves DNS records from Route53 and returns error.
func (r *Route53Srv) GetRecords() error {
	// If Srv doesn't exist, create it
	if r.Srv == nil {
		r.Srv = route53.New(nil)
	}
	config := r.Config
	zoneid := config.ZoneId()
	params := &route53.ListResourceRecordSetsInput{HostedZoneID:aws.String(zoneid)}
	resp, err := r.Srv.ListResourceRecordSets(params)
	if awserr := aws.Error(err); awserr != nil {
	    // A service error occurred.
	    return fmt.Errorf("Service Error: %#s : %#s\n", awserr.Code, awserr.Message)
	} else if err != nil {
	    // A non-service error occurred.
	    return fmt.Errorf("Non-Service Error: %#s\n",err)
	}
	for _,rec := range resp.ResourceRecordSets {
		record := Record{
			Type: *rec.Type,
			Name: *rec.Name,
			Value: *rec.ResourceRecords[0].Value,
		}
		r.records = append(r.records, &record)
	}
	return nil
}

// ChangeRecords will execute a set of DNS changes in Route53, returns error.
func (r *Route53Srv) ChangeRecords() (error) {
	// If Srv doesn't exist, create it
	if r.Srv == nil {
		r.Srv = route53.New(nil)
	}
	changebatch := &route53.ChangeBatch{
		Changes: r.Changes,
		Comment: aws.String("Change sync from consul"),
	}
	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: changebatch,
		HostedZoneID: aws.String(r.ZoneId()),
	}
	_, err := r.Srv.ChangeResourceRecordSets(params)
	if awserr := aws.Error(err); awserr != nil {
		return fmt.Errorf("Error: %#s : %#s\n", awserr.Code, awserr.Message)
	}
	if err != nil {
		return err
	}
	return nil
}

// AddChange will handed "CHANGE","DELETE" or "UPSERT" and a pointer to a Record will add that change to Changes.
func (r *Route53Srv) AddChange(changetype string, rec Record){
	recordchange := &route53.ResourceRecordSet{
		Name: &rec.Name,
		ResourceRecords: []*route53.ResourceRecord{ &route53.ResourceRecord{Value: &rec.Value} },
		TTL: aws.Long(r.TTL()),
		Type: &rec.Type,
	}
	change := &route53.Change{
		Action: aws.String(changetype),
		ResourceRecordSet: recordchange,
	}
	r.Changes = append(r.Changes,change)
}
