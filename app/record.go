package consul2route53

type Record struct {
	Type string
	Name string
	Value string
	Weight int64	
}