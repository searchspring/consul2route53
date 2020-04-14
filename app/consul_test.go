package consul2route53

import (
	"testing"
	"github.com/searchspring/consul2route53/testutil"

)

func TestGetNodes(t *testing.T) {
	
	mockconsul := testutil.Mockconsul{}
	err :=	mockconsul.MockConsul()
	if err != nil {
		t.Error(err)
	}
	defer mockconsul.Close()

	consul := &Consul{
		Config: &Config{
			Consulhost: mockconsul.Host,
			Consulport: mockconsul.Port,
		},
	}
	// Test GetNodes and GetNodeServices
	err = consul.GetNodes()
	if err != nil {
		t.Error(err)
	}
	if len(consul.Nodes()) != 2 {
		t.Errorf("Expected to get 2 nodes, got %#v.\n", len(consul.Nodes()))
	}

	// Test GetServices
	err = consul.GetServices()
	if err != nil {
		t.Error(err)
	}
	if len(consul.Services()) != 1 {
		t.Errorf("Expected to get 1 service, got %#v.\n", len(consul.Services()))
	}

	consul.SetConsulPort("8600")
	err = consul.GetNodes()
	if err == nil {
		t.Errorf("Expected to get error with 8600 for consul port\n")
	}
	
}
