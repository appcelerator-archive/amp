package core

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/appcelerator/amp/pkg/docker"
	"github.com/appcelerator/amp/pkg/labels"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

const (
	dockerStackLabelName = "com.docker.stack.namespace"
)

//HAProxy haproxy struct
type HAProxy struct {
	docker             *client.Client
	mappingMap         map[string]*publicMapping
	exec               *exec.Cmd
	isLoadingConf      bool
	dnsRetryLoopID     int
	dnsNotResolvedList []string
	updateID           int
	updateChannel      chan int
	iterationNumber    int
}

type publicMapping struct {
	iterationNumber int    //internal counter to detect removed mapping
	account         string //Account or stack owner namem used in the url host part
	stack           string //Stack name, used in the url host part
	service         string //service to reach
	label           string //label used in the url host part
	port            string //internal service port
}

var (
	haproxy HAProxy
	ctx     context.Context
)

//Set app mate initial values
func (app *HAProxy) init() {
	app.mappingMap = make(map[string]*publicMapping)
	ctx = context.Background()
	if err := dockerInit(); err != nil {
		fmt.Printf("Init error: %v\n", err)
		os.Exit(1)
	}
	app.updateChannel = make(chan int)
	app.isLoadingConf = false
	haproxy.updateConfiguration(true)
}

func dockerInit() error {
	// Connection to Docker
	defaultHeaders := map[string]string{"User-Agent": "haproxy"}
	cli, err := client.NewClient(docker.DefaultURL, docker.DefaultVersion, nil, defaultHeaders)
	if err != nil {
		return err
	}
	haproxy.docker = cli
	fmt.Println("Connected to Docker-engine")
	return nil
}

//Launch a routine to catch SIGTERM Signal
func (app *HAProxy) trapSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		<-ch
		fmt.Println("\namp-haproxy-controller received SIGTERM signal")
		app.docker.Close()
		os.Exit(1)
	}()
}

//Launch HAProxy using cmd command
func (app *HAProxy) start() {
	//launch HAPRoxy
	go func() {
		fmt.Println("launching HAProxy on initial configuration")
		app.exec = exec.Command("haproxy", "-f", "/usr/local/etc/haproxy/haproxy.cfg")
		app.exec.Stdout = os.Stdout
		app.exec.Stderr = os.Stderr
		err := app.exec.Run()
		if err != nil {
			fmt.Printf("HAProxy exit with error: %v\n", err)
			app.docker.Close()
			os.Exit(1)
		}
	}()
	//Launch main docker watch loop
	fmt.Println("launching Docker watch")
	app.dockerWatch()
	//Launch main haproxy update loop
	for {
		uid := <-app.updateChannel
		if uid == app.updateID {
			app.updateConfiguration(true)
		}
	}
}

//Stop HAProxy
func (app *HAProxy) stop() {
	fmt.Println("Send SIGTERM signal to HAProxy")
	if app.exec != nil {
		app.exec.Process.Kill()
	}
}

func (app *HAProxy) dockerWatch() {
	go func() {
		for {
			time.Sleep(time.Duration(conf.dockerWatchPeriod) * time.Second)
			if app.updateServiceMap() {
				app.updateID++
				app.updateChannel <- app.updateID
			}
		}
	}()
}

func (app *HAProxy) updateServiceMap() bool {
	list, err := app.docker.ServiceList(ctx, types.ServiceListOptions{})
	if err != nil || len(list) == 0 {
		fmt.Printf("Error reading docker service list: %v\n", err)
		return false
	}
	isItemUpdated := false
	for _, serv := range list {
		if serv.Spec.Labels != nil {
			stackName, exist := serv.Spec.Labels[dockerStackLabelName]
			if exist {
				serviceName := serv.Spec.Annotations.Name
				mappings, ok := serv.Spec.Labels[labels.LabelsNameMapping]
				if ok {
					if app.addMappings(stackName, serviceName, mappings) {
						isItemUpdated = true
					}
				}
			}
		}

	}
	return isItemUpdated
}

func (app *HAProxy) addMappings(stackName string, serviceName string, mappings string) bool {
	isItemUpdated := false
	app.iterationNumber++
	mappingList := strings.Split(mappings, ",")
	tmpMap := make(map[string]*publicMapping)
	for _, value := range mappingList {
		data, err := app.evalMappingString(stackName, serviceName, strings.TrimSpace(value))
		if err != nil {
			fmt.Printf("Mapping error for service %s: %v\n", serviceName, err)
			return false
		}
		fmt.Printf("Found mapping %s: %v\n", value, data)
		data.iterationNumber = app.iterationNumber
		mappingID := fmt.Sprintf("%s-%s-%s:%s", data.account, data.stack, data.label, data.port)
		mapping, exist := app.mappingMap[mappingID]
		if !exist {
			isItemUpdated = true
			tmpMap[mappingID] = data
		} else {
			mapping.iterationNumber = app.iterationNumber
		}
	}
	for id, item := range app.mappingMap {
		if item.iterationNumber == app.iterationNumber {
			tmpMap[id] = item
		} else {
			isItemUpdated = true
		}
	}
	app.mappingMap = tmpMap
	return isItemUpdated
}

//Launch HAProxy using cmd command
func (app *HAProxy) reloadConfiguration() {
	app.isLoadingConf = true
	fmt.Println("reloading HAProxy configuration")
	if app.exec == nil {
		fmt.Printf("HAProxy is not started yet, waiting for it\n")
		return
	}
	pid := app.exec.Process.Pid
	fmt.Printf("Execute: %s %s %s %s %d\n", "haproxy", "-f", "/usr/local/etc/haproxy/haproxy.cfg", "-sf", pid)
	app.exec = exec.Command("haproxy", "-f", "/usr/local/etc/haproxy/haproxy.cfg", "-sf", fmt.Sprintf("%d", pid))
	app.exec.Stdout = os.Stdout
	app.exec.Stderr = os.Stderr
	go func() {
		err := app.exec.Run()
		app.isLoadingConf = false
		if err == nil {
			fmt.Println("HAProxy configuration reloaded")
			return
		}
		fmt.Printf("HAProxy reload configuration error: %v\n", err)
		os.Exit(1)
	}()
}

// update configuration managing the isUpdatingConf flag
func (app *HAProxy) updateConfiguration(reload bool) error {
	app.dnsRetryLoopID++
	app.dnsNotResolvedList = []string{}
	err := app.updateConfigurationEff(reload)
	if err == nil {
		app.startDNSRevolverLoop(app.dnsRetryLoopID)
	}
	return err
}

//update HAProxy configuration for master regarding ETCD keys values and make HAProxy reload its configuration if reload is true
func (app *HAProxy) updateConfigurationEff(reload bool) error {
	fmt.Printf("update HAProxy configuration, mapping len=%d\n", len(app.mappingMap))
	fileNameTarget := "/usr/local/etc/haproxy/haproxy.cfg"
	fileNameTpt := "/usr/local/etc/haproxy/haproxy.cfg.tpt"
	file, err := os.Create(fileNameTarget + ".new")
	if err != nil {
		fmt.Printf("Error creating new haproxy conffile for creation: %v\n", err)
		return err
	}
	filetpt, err := os.Open(fileNameTpt)
	if err != nil {
		fmt.Printf("Error opening conffile template: %s : %v\n", fileNameTpt, err)
		return err
	}
	scanner := bufio.NewScanner(filetpt)
	if conf.debug {
		fmt.Printf("Updating with mappings: %v\n", app.mappingMap)
	}
	for scanner.Scan() {
		line := scanner.Text()
		if conf.debug {
			fmt.Printf("line: %s\n", line)
		}
		if strings.HasPrefix(strings.Trim(line, " "), "[frontendInline]") {
			app.writeFrontendInline(file)
		} else if strings.HasPrefix(strings.Trim(line, " "), "[backends]") {
			app.writeBackend(file)
		} else {
			file.WriteString(line + "\n")
		}
	}
	if err = scanner.Err(); err != nil {
		fmt.Printf("Error reading haproxy conffile template: %s %v\n", fileNameTpt, err)
		file.Close()
		return err
	}
	file.Close()
	os.Remove(fileNameTarget)
	err2 := os.Rename(fileNameTarget+".new", fileNameTarget)
	if err2 != nil {
		fmt.Printf("Error renaming haproxy conffile .new: %v\n", err)
		return err
	}
	fmt.Println("HAProxy configuration updated")
	if reload {
		app.reloadConfiguration()
	}
	return nil
}

// write backends for main service configuration
func (app *HAProxy) writeFrontendInline(file *os.File) error {
	for _, mapping := range app.mappingMap {
		if mapping.account != "" {
			line := fmt.Sprintf("    use_backend bk_%s-%s-%s if { hdr_beg(host) -i %s.%s.%s }\n", mapping.account, mapping.service, mapping.port, mapping.label, mapping.stack, mapping.account)
			file.WriteString(line)
			fmt.Printf(line)
		} else {
			line := fmt.Sprintf("    use_backend bk_-%s-%s if { hdr_beg(host) -i %s.%s }\n", mapping.service, mapping.port, mapping.label, mapping.stack)
			file.WriteString(line)
			fmt.Printf(line)
		}
	}
	return nil
}

// write backends for stack haproxy configuration
func (app *HAProxy) writeBackend(file *os.File) error {
	for _, mapping := range app.mappingMap {
		dnsResolved := app.tryToResolvDNS(mapping.service)
		line := fmt.Sprintf("\nbackend bk_%s-%s-%s\n", mapping.account, mapping.service, mapping.port)
		file.WriteString(line)
		fmt.Printf(line)
		//if dns name is not resolved haproxy (v1.6) won't start or accept the new configuration so server is disabled
		//to be removed when haproxy will fixe this bug
		if dnsResolved {
			line1 := fmt.Sprintf("    server %s_1 %s:%s resolvers docker resolve-prefer ipv4\n", mapping.service, mapping.service, mapping.port)
			file.WriteString(line1)
			fmt.Printf(line1)
		} else {
			line1 := fmt.Sprintf("    #dns name %s not resolved\n", mapping.service)
			file.WriteString(line1)
			fmt.Printf(line1)
			line2 := fmt.Sprintf("    #server %s_1 %s:%s resolvers docker resolve-prefer ipv4\n", mapping.service, mapping.service, mapping.port)
			file.WriteString(line2)
			fmt.Printf(line2)
			app.addDNSNameInRetryList(mapping.service)
		}
	}
	return nil
}

// test if a dns name is resolved or not
func (app *HAProxy) tryToResolvDNS(name string) bool {
	_, err := net.LookupIP(name)
	if err != nil {
		return false
	}
	return true
}

// add unresolved dns name in list to be retested later
func (app *HAProxy) addDNSNameInRetryList(name string) {
	app.dnsNotResolvedList = append(app.dnsNotResolvedList, name)
}

// on regular basis try to see if one of the unresolved dns become resolved, if so execute a configuration update.
// need to have only one loop at a time, if the id change then the current loop should stop
// id is incremented at each configuration update which can be trigger by ETCD wash also
func (app *HAProxy) startDNSRevolverLoop(loopID int) {
	//if no unresolved DNS name then not needed to start the loop
	if len(haproxy.dnsNotResolvedList) == 0 {
		return
	}
	fmt.Printf("Start DNS resolver id: %d\n", loopID)
	go func() {
		for {
			for _, name := range haproxy.dnsNotResolvedList {
				if app.tryToResolvDNS(name) {
					if haproxy.dnsRetryLoopID == loopID {
						fmt.Printf("DNS %s resolved, update configuration\n", name)
						app.updateConfiguration(true)
					}
					fmt.Printf("Stop DNS resolver id: %d\n", loopID)
					return
				}
			}
			time.Sleep(10)
			if haproxy.dnsRetryLoopID != loopID {
				fmt.Printf("Stop DNS resolver id: %d\n", loopID)
				return
			}
		}
	}()
}

// parse a value of io.amp.mapping label
func (app *HAProxy) evalMappingString(stack string, service string, labelValue string) (*publicMapping, error) {
	mapping := &publicMapping{
		stack:   app.getSimpleName(stack),
		service: service,
		label:   app.getDefaultName(service),
	}
	if labelValue == "" {
		return mapping, nil
	}
	items := strings.Split(labelValue, " ")
	for _, item := range items {
		param := strings.Split(item, "=")
		if len(param) == 0 {
			param = strings.Split(item, ":")
		}
		if len(param) != 2 {
			return nil, fmt.Errorf("Error in label %s: parameter %s should have '=' or ':' separator", labelValue, item)
		}
		name := strings.TrimSpace(param[0])
		value := strings.TrimSpace(param[1])
		if name == "account" {
			mapping.account = value
		} else if name == "name" {
			mapping.label = value
		} else if name == "port" {
			mapping.port = value
		}
	}
	if mapping.port == "" {
		return nil, fmt.Errorf("Error in label %s: need to define the targetted port (port=xxx)", labelValue)
	}
	return mapping, nil
}

func (app *HAProxy) getDefaultName(name string) string {
	list := strings.Split(name, "_")
	if len(list) == 1 {
		return name
	}
	ret := list[1]
	if len(list) > 2 {
		for _, item := range list[2:len(list)] {
			ret = fmt.Sprintf("%s_%s", ret, item)
		}
	}
	return ret
}

func (app *HAProxy) getSimpleName(name string) string {
	list := strings.Split(name, "-")
	ret := list[0]
	if len(list) > 2 {
		for _, item := range list[1 : len(list)-1] {
			ret = fmt.Sprintf("%s-%s", ret, item)
		}
	}
	return ret
}
