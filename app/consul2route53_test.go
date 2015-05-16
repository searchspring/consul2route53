package consul2route53

import(
	"testing"
	"github.com/SearchSpring/consul2route53/testutil"
)

func TestRun(t *testing.T) {
	mockconsul := testutil.Mockconsul{}
	err :=	mockconsul.MockConsul()
	if err != nil {
		t.Error(err)
	}
	defer mockconsul.Close()
	config := &Config{
		Consulhost: mockconsul.Host,
		Consulport: mockconsul.Port,
		Zoneid: "bogusbogus",
		Domain: "bogus.com.",
		Ttl: 300,
	}
	consul := &Consul2Route53{
		Config: config,
		Consul: &Consul{Config: config},
		Route53Srv: &Route53Srv{
			Config: config,
			Srv: new(testutil.MockRoute53),
		},
	}
	err = consul.Run()
	if err != nil {
		t.Error(err)
	}
}