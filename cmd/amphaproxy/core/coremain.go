package core

import (
	"fmt"
	"time"
)

//Run launch main loop
func Run(version string, build string) {
	//ampHAProxyControllerVersion = version
	conf.load(version, build)
	haproxy.init()
	haproxy.trapSignal()
	initAPI()
	haproxy.waitForAmplifier()
	haproxy.start()
}

func (app *HAProxy) waitForAmplifier() {
	amplifierName := fmt.Sprintf("%s_amplifier", conf.ampStackName)
	fmt.Printf("Waiting for %s\n", amplifierName)
	for !app.tryToResolvDNS(amplifierName) {
		time.Sleep(3 * time.Second)
	}
	fmt.Printf("%s is ready\n", amplifierName)
}
