package consul2route53

type Service struct {
	ID string
	Service string
	Tags []string
	Port int64
	Address string
}

type Srv struct{
	Node string
	Address string
	ServiceID string
	ServiceName string
	ServiceTags []string
	ServiceAddress string
	ServicePort int64
}