package consul2route53

import (
	"os"
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestReadConfig(t *testing.T) {
	// Test missing config
	_,err := ReadConfig("bogusbogus")
	if err == nil {
		t.Error("Expect missing file error, didn't get it.")
	}
	// Testing broken config
	_,err = ReadConfig("config.go")
	if err == nil {
		t.Error("Expected yaml error, didn't get it.")
	}
	// Testing normal config
	testpath, err := ioutil.TempDir("","consul253")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(testpath)
	testcfg := map[string]string{
		"consulhost": "127.0.0.1",
		"consulport": "8500",
		"zoneid": "bogus",
		"domain": "service.consul",
		"delegrationsetid": "bogus",
	}
	yaml,err := yaml.Marshal(testcfg)
	if err != nil {
		t.Error(err)
	}
	testcfgfile := testpath+"/consul253.yml"
	err = ioutil.WriteFile(testcfgfile, yaml, 0644)
	if err != nil {
		t.Error(err)
	}
	config,err := ReadConfig(testcfgfile)
	if err != nil {
		t.Errorf("Failed to read config file: %s\n", err)
	}
	
	if config.ConsulHost() != "127.0.0.1" {
		t.Errorf("Expected 'Consulhost' to be '127.0.0.1' and got '%s'\n", config.Consulhost)
	}
	
	if config.ConsulPort() != "8500" {
		t.Errorf("Expected 'consulport' to be '8500' and got '%s'\n", config.Consulport)
	}
	
	if config.Zoneid != "bogus" {
		t.Errorf("Expected 'zoneid' to be 'bogus and got '%s'\n", config.Zoneid)
	}

	config.SetConsulHost("localhost")
	if config.ConsulHost() != "localhost" {
		t.Errorf("Expected 'ConsulHost' to be 'localhost' and got '%s'\n", config.ConsulHost())
	}

	config.SetConsulPort("8600")
	if config.ConsulPort() != "8600" {
		t.Errorf("Expected 'ConsulPort' to be '8600' and got '%#v'\n", config.ConsulPort()  )
	}

	config.SetZoneId("testbogus")
	if config.ZoneId() != "testbogus" {
		t.Errorf("Expected 'Zoneid' to be 'testbogus' and got '%#v'\n", config.ZoneId())
	}

	config.SetZone("testbogus.com.")
	if config.Zone() != "testbogus.com." {
		t.Errorf("Expected 'Zone' to be 'testbogus.com.' and got '%#v'\n", config.Zone())
	}

	config.SetTTL(300)
	if config.TTL() != 300 {
		t.Errorf("Expected 'TTL' to be '300' and got '%#v'\n,", config.TTL())
	}



}