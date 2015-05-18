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
	route53srv := &Route53Srv{ Config: config }
	route53srv.SetSrv(new(testutil.MockRoute53))
	mock := new(testutil.MockRoute53)
	mock.Fail = true
	route53srv_fail := &Route53Srv{Config: config}
	route53srv_fail.SetSrv(mock)

	// Test GetZoneInfo
	err := route53srv.GetZoneInfo()
	if err != nil {
		t.Error(err)
	}
	if route53srv.Zone() != "bogus.com." {
		t.Errorf("GetZoneinfo Expected to get 'bogus.com.' got %#v.\n",route53srv.Zone())
	}
	err = route53srv_fail.GetZoneInfo()
	if err == nil {
		t.Error("GetZoneinfo Expected error and didn't get it\n")
	}

	
	// Test GetRecords
	err = route53srv.GetRecords()
	if err != nil {
		t.Error(err)
	}
	if len(route53srv.Records()) != 2 {
		t.Errorf("GetRecords Expected to get 1 record, got %#v.\n", len(route53srv.Records()))
	}
	err = route53srv_fail.GetRecords()
	if err == nil {
		t.Error("GetRecords Expected error and didn't get it\n")
	}
	// Test AddChange
	for _,record := range route53srv.Records() {
		route53srv.AddChange("DELETE",*record)
	}
	if len(route53srv.Changes) != 2 {
		t.Errorf("Addchange Expected to add 1 record to Changes, got %#v.\n", len(route53srv.Changes))
	}

	// Test ChangeRecords
	err = route53srv.ChangeRecords()
	if err != nil {
		t.Error(err)
	}
	err = route53srv_fail.ChangeRecords()
	if err == nil {
		t.Error("ChangeRecords Expected error and didn't get it\n")
	}
}