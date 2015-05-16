consul2route53
==============

Sync your Consul DNS records to a private zone in Route53.

Usage
=====
```
  -config="./consul2route53.conf": Path to configfile
  -host="localhost": Address of consul server
  -once=false: Run Once
  -port="8500": Port of consul server
  -sleeptime=1000: Sleep time in loop
  -ttl=300: TTL
  -version=false: consul253 version
  -zoneid="": Route53 ZoneID
```

License and Author
==================

* Author:: Greg Hellings (<greg@thesub.net>)


Copyright 2015, B7 Interactive, LLC.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.