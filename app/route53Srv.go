 package consul2route53

import(
	"fmt"
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/route53"
)

type Route53Srv struct {
	*Config
	records []*Record
	srv SrvRoute53
	zone string
	Changes []*route53.Change
}

type SrvRoute53 interface {
	ListResourceRecordSets(*route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput,error)
	ChangeResourceRecordSets(*route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput,error)
	GetHostedZone(*route53.GetHostedZoneInput)(*route53.GetHostedZoneOutput,error)
}

func (r *Route53Srv) SetSrv(srvroute53 SrvRoute53) {
	r.srv = srvroute53
}

func (r *Route53Srv) Srv() SrvRoute53 {
	if r.srv == nil {
		r.SetSrv(route53.New(nil))
	}
	return r.srv
}

func (r *Route53Srv) SetZone(zone string){
	r.zone = zone
}
func (r *Route53Srv) Zone() string {
	if r.zone == "" {
		err := r.GetZoneInfo()
		if err != nil {
			panic(err)
		}
	}
	return r.zone
}


func (r *Route53Srv) ChangesNum() int {
	return len(r.Changes)
}

func (r *Route53Srv) Records() []*Record {
	return r.records
}

func (r *Route53Srv) RecordsMap() (map[string]*Record, error) {
	recmap := make(map[string]*Record)
	err := r.GetRecords()
	if err != nil {
		return recmap, err
	}
	for _,record := range r.records {
		recmap[record.Name+"_"+record.Type] = record
	}
	return recmap, err
}

func (r *Route53Srv) GetZoneInfo() error {
	srv := r.Srv()
	params := &route53.GetHostedZoneInput{ID: aws.String(r.ZoneId())}
	resp, err := srv.GetHostedZone(params)
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
	srv := r.Srv()
	config := r.Config
	zoneid := config.ZoneId()
	params := &route53.ListResourceRecordSetsInput{HostedZoneID:aws.String(zoneid),MaxItems:aws.String("100")}
	resp, err := srv.ListResourceRecordSets(params)
	loop_resp := true
	for loop_resp {	
		if awserr := aws.Error(err); awserr != nil {
		    // A service error occurred.
		    return fmt.Errorf("Service Error: %#s : %#s\n", awserr.Code, awserr.Message)
		} else if err != nil {
		    // A non-service error occurred.
		    return fmt.Errorf("Non-Service Error: %#s\n",err)
		}

		// AWS will max give us 100 records at a time (100_000 max), so loop if there's more
		r.records = make([]*Record,0,100000)
		for _,rec := range resp.ResourceRecordSets {
			record := Record{
				Type: *rec.Type,
				Name: *rec.Name,
				Value: *rec.ResourceRecords[0].Value,
			}
			r.records = append(r.records, &record)
		}
		if *resp.IsTruncated {
			loop_resp = true
		} else {
			loop_resp = false
		}
		params.StartRecordIdentifier = resp.NextRecordIdentifier
		resp,err = srv.ListResourceRecordSets(params)
	}
	return nil
}

// ChangeRecords will execute a set of DNS changes in Route53, returns error.
func (r *Route53Srv) ChangeRecords() (error) {
	// If Srv doesn't exist, create it
	srv := r.Srv()
	changebatch := &route53.ChangeBatch{
		Changes: r.Changes,
		Comment: aws.String("Change sync from consul"),
	}
	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: changebatch,
		HostedZoneID: aws.String(r.ZoneId()),
	}
	_, err := srv.ChangeResourceRecordSets(params)
	if awserr := aws.Error(err); awserr != nil {
		return fmt.Errorf("Error: %#s : %#s\n", awserr.Code, awserr.Message)
	}
	if err != nil {
		return err
	}
	r.Changes = nil
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
