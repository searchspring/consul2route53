package main

import(
	"os/signal"
	"os"
	"fmt"
	"syscall"
	"time"
	"flag"
	"log"
	"github.com/SearchSpring/consul2route53/app"
)

var consulhost string
var consulport string
var zoneid string
var configfile string
var sleeptime int64
var ttl int64
var version bool
var once bool
const (
	versioninfo = "v0.0.1"
)

type loopobject interface{
	Run() error
}

func main() {
	flag.StringVar(&consulhost, "host", "localhost", "Address of consul server")
	flag.StringVar(&consulport, "port", "8500", "Port of consul server")
	flag.StringVar(&zoneid, "zoneid", "", "Route53 ZoneID")
	flag.StringVar(&configfile, "config", "./consul2route53.conf", "Path to configfile")
	flag.Int64Var(&sleeptime, "sleeptime", 1000, "Sleep time in loop")
	flag.BoolVar(&once, "once", false, "Run Once")
	flag.Int64Var(&ttl, "ttl", 300, "TTL")
	flag.BoolVar(&version, "version", false, "consul253 version")
	flag.Parse()

	for {
		config,_ := consul2route53.ReadConfig(configfile) 
		config.SetConsulHost(consulhost)
		config.SetConsulPort(consulport)
		config.SetZoneId(zoneid)
		config.SetTTL(ttl)
		consul2route53 := consul2route53.New(config)
		switch {
			case once:
				err := consul2route53.Run()
				if err != nil {
					log.Fatal(err)
				}
				return
			default:
				loop(consul2route53)
		}
	}
}


// Loop until SIGHUP
func loop(obj loopobject) {
	c := make(chan os.Signal, 1)
	r := make(chan bool)

	signal.Notify(c, syscall.SIGHUP)
 
	go func(){
		for sig := range c {
			fmt.Println(sig)
			r <- true
		}
	}()
	for {
		select {
			case msg := <-r:
				log.Printf("Got A HUP. Reloading: %#v\n", msg)
				return
			default:
		}
		err := obj.Run()
		if err != nil {
			log.Println(err)
		}
		time.Sleep(time.Duration(sleeptime * int64(time.Millisecond)))
	}

}
