package consul2route53

import(
	"testing"
	"github.com/searchspring/consul2route53/testutil"
)

func TestRun(t *testing.T) {
	mockconsul := testutil.Mockconsul{}
	err :=	mockconsul.MockConsul()
	if err != nil {
		t.Error(err)
	}
	defer mockconsul.Close()
	config := Config{
		Consulhost: mockconsul.Host,
		Consulport: mockconsul.Port,
		Zoneid: "bogusbogus",
		Ttl: 300,
	}
	consul := New(config)
	consul.SetZone("bogus.com.")
	consul.SetSrv(new(testutil.MockRoute53))
	err = consul.Run()
	if err != nil {
		t.Error(err)
	}

	mockroute53_fail := new(testutil.MockRoute53)
	mockroute53_fail.Fail = true
	consul.SetSrv(mockroute53_fail)
	err = consul.Run()
	if err == nil {
		t.Error("Expected consul2route53 to fail and it did not.")
	}

}
