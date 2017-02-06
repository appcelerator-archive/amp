package core

import (
	"fmt"
	"time"
)

var ampHAProxyControllerVersion string

//Run launch main loop
func Run(version string, build string) {
	ampHAProxyControllerVersion = version
	conf.load(version, build)
	err := etcdClient.init()
	for err != nil {
		fmt.Printf("Waiting for ETCD connection\n")
		time.Sleep(5 * time.Second)
		err = etcdClient.init()
	}
	haproxy.init()
	haproxy.trapSignal()
	initAPI()
	haproxy.start()
	time.Sleep(10 * time.Second)
	for {
		time.Sleep(3 * time.Second)
		etcdClient.watchForServicesUpdate()
	}
}
