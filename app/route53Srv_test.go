 package consul2route53

import(
	"testing"
	"github.com/SearchSpring/consul2route53/testutil"
)



func TestRoute53SrvStruct(t *testing.T) {
	config := &Config{
		Zoneid: "bogusbogus",
		Ttl: 300,
	}
	route53srv := &Route53Srv{
		Config: config, 
		Srv: new(testutil.MockRoute53),
	}

	// Test GetZoneInfo
	err := route53srv.GetZoneInfo()
	if err != nil {
		t.Error(err)
	}
	if route53srv.Zone() != "bogus.com." {
		t.Errorf("GetZoneinfo Expected to get 'bogus.com.' got %#v.\n",route53srv.Zone())
	} 
	
	// Test GetRecords
	err = route53srv.GetRecords()
	if err != nil {
		t.Error(err)
	}
	if len(route53srv.Records()) != 1 {
		t.Errorf("GetRecords Expected to get 1 record, got %#v.\n", len(route53srv.Records()))
	}

	// Test AddChange
	for _,record := range route53srv.Records() {
		route53srv.AddChange("DELETE",*record)
	}
	if len(route53srv.Changes) != 1 {
		t.Errorf("Addchange Expected to add 1 record to Changes, got %#v.\n", len(route53srv.Changes))
	}

	// Test ChangeRecords
	err = route53srv.ChangeRecords()
	if err != nil {
		t.Error(err)
	}
}