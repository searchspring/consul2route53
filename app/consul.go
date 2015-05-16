package consul2route53

import(
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

// Consul represent information from a consul cluster
type Consul struct {
	*Config
	nodes []*Node 
	services []*Service
}

type Json interface {
}

func (c *Consul) Services() []*Service {
	return c.services
}

func (c *Consul) ServicesMap() map[string]*Service {
	srvmap := make(map[string]*Service)
	for _,service := range c.Services() {
		srvmap[service.Service] = service
	}
	return srvmap
}

func (c *Consul) Nodes() []*Node {
	return c.nodes
}

// getJson submits api requests to consul and returns error.
func (c *Consul) getJson(path string, data interface{}) error {
	resp,err := http.Get("http://"+c.ConsulHost()+":"+c.ConsulPort()+"/v1/"+path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	contents,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	
	if (resp.StatusCode >= 400) {
		var apiError string 
		err = json.Unmarshal(contents, apiError)
		return fmt.Errorf("Request Path: %s\nHTTP Status: %d\nApiError: %s\nError: %s\n",path,resp.StatusCode,apiError, err)
	}

	err = json.Unmarshal(contents,data)
	if err != nil {
		return fmt.Errorf("Unmarshal data error: %s\ndata: %#v\n", err, string(contents))
	}
	return nil
}

// GetNodes connects to a consul server and retrieves information about nodes.  Returns error.
func (c *Consul) GetNodes() error {
	path := "catalog/nodes"
	var nodes []*Node 
	err := c.getJson(path,&nodes)
	if err != nil {
		return err
	}
	c.nodes = nodes
	return nil
}

// GetServices connects to consul and gets a list of services. Returns error.
func (c *Consul) GetServices() error {
	path := "catalog/services"
	var servicenames map[string][]string
	err := c.getJson(path,&servicenames)
	for servicename,_ := range servicenames {
	 	service,err := c.GetServiceInfo(servicename)
		if err != nil {
			return err
		}
		c.services = append(c.services, service)
	}
	return err
}

// GetServiceInfo given a service name, connects to consul and retrieves information about that service.  Returns *Service or error.
func (c *Consul) GetServiceInfo(servicename string) (*Service,error) {
	path := "catalog/service/"+servicename
	var srv []Srv
	var service *Service
	err := c.getJson(path,&srv)
	if len(srv) > 0 {
		srv1 := srv[0]
		service = &Service{
			ID: srv1.ServiceID,
			Service: srv1.ServiceName,
			Tags: srv1.ServiceTags,
			Port: srv1.ServicePort,
			Address: srv1.Address,
		}
	}
	return service, err
}