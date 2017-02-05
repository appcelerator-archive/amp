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
)

//HAProxy haproxy struct
type HAProxy struct {
	exec               *exec.Cmd
	isLoadingConf      bool
	dnsRetryLoopID     int
	dnsNotResolvedList []string
	updateID           int
	updateChannel      chan int
}

type publicMapping struct {
	stack   string //Stack name, used in the url host part
	service string //service to reach
	label   string //label used in the url host part
	port    string //internal service port
	mode    string //only for tcp mode (grpc)
	portTo  string //internal grpc port
}

var (
	haproxy HAProxy
)

//Set app mate initial values
func (app *HAProxy) init() {
	app.updateChannel = make(chan int)
	app.isLoadingConf = false
	app.dnsNotResolvedList = []string{}
	haproxy.updateConfiguration(true)
}

//Launch a routine to catch SIGTERM Signal
func (app *HAProxy) trapSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		<-ch
		fmt.Println("\namp-haproxy-controller received SIGTERM signal")
		etcdClient.Close()
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
			etcdClient.Close()
			os.Exit(1)
		}
	}()
	//Lanch main update loop
	go func() {
		for {
			uid := <-app.updateChannel
			if uid == app.updateID {
				app.updateConfiguration(true)
			}
		}
	}()
}

//Stop HAProxy
func (app *HAProxy) stop() {
	fmt.Println("Send SIGTERM signal to HAProxy")
	if app.exec != nil {
		app.exec.Process.Kill()
	}
}

//Launch HAProxy using cmd command
func (app *HAProxy) reloadConfiguration() {
	app.isLoadingConf = true
	fmt.Println("reloading HAProxy configuration")
	if app.exec == nil {
		fmt.Printf("HAProxy not started yet, waiting it started\n")
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
	fmt.Println("update HAProxy configuration")
	mappingList, err := etcdClient.getAllMappings()
	if err != nil {
		fmt.Println("Erreur on get maapingList: ", err)
		return err
	}
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
	skip := false
	for scanner.Scan() {
		line := scanner.Text()
		skip = hasToBeSkipped(line, skip)
		if conf.debug {
			fmt.Printf("line: %t: %s\n", skip, line)
		}
		if !skip {
			if strings.HasPrefix(strings.Trim(line, " "), "[frontendInline]") {
				app.writeFrontendInline(file, mappingList)
			} else if strings.HasPrefix(strings.Trim(line, " "), "[backends]") {
				app.writeBackend(file, mappingList)
			} else if strings.HasPrefix(strings.Trim(line, " "), "[frontend]") {
				app.writeFrontend(file, mappingList)
			} else {
				file.WriteString(line + "\n")
			}
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

// compute if the line is the begining of a block which should be skipped or not
func hasToBeSkipped(line string, skip bool) bool {
	if line == "" {
		//if blanck line then end of skip
		return false
	} else if skip {
		// if skipped mode and line not black then continue to skip
		return true
	}

	ref := strings.Trim(line, " ")
	if conf.noDefaultBackend {
		if strings.HasPrefix(ref, "frontend main_") || strings.HasPrefix(ref, "backend main_") || strings.HasPrefix(ref, "backend stack_") {
			//if nodefaultbackend activated and default frontend or backend then skip
			return true
		}
	} else {
		if conf.stackName != "" {
			if strings.HasPrefix(ref, "frontend main_") || strings.HasPrefix(ref, "backend main_") {
				//if stack mode and main frontend or backend then skip
				return true
			}
		} else {
			if strings.HasPrefix(ref, "frontend stack_") || strings.HasPrefix(ref, "backend stack_") {
				//if main mode and stack frontend or backend then skip
				return true
			}
		}
	}
	//if line not "" and not skip then continue to not skip
	return false
}

// write backends for main service configuration
func (app *HAProxy) writeFrontendInline(file *os.File, mappingList []*publicMapping) error {
	for _, mapping := range mappingList {
		if mapping.mode != "tcp" {
			line := fmt.Sprintf("    use_backend bk_%s_%s-%s if { hdr_beg(host) -i %s.%s. }\n", mapping.stack, mapping.service, mapping.port, mapping.label, mapping.stack)
			file.WriteString(line)
			fmt.Printf(line)
		}
	}
	return nil
}

// write backends for main service configuration
func (app *HAProxy) writeFrontend(file *os.File, mappingList []*publicMapping) error {
	for _, mapping := range mappingList {
		if mapping.mode == "tcp" {
			lines := []string{
				fmt.Sprintf("\nfrontend %s_%s_grpc\n", mapping.stack, mapping.service),
				"    mode tcp\n",
				fmt.Sprintf("    bind *:%s npn spdy/2 alpn h2,http/1.1\n", mapping.portTo),
				fmt.Sprintf("    default_backend bk_%s_%s-%s\n\n", mapping.stack, mapping.service, mapping.port)}
			for _, line := range lines {
				file.WriteString(line)
				fmt.Printf(line)
			}
		}
	}
	return nil
}

// write backends for stack haproxy configuration
func (app *HAProxy) writeBackend(file *os.File, mappingList []*publicMapping) error {
	for _, mapping := range mappingList {
		dnsResolved := app.tryToResolvDNS(mapping.service)
		line := fmt.Sprintf("\nbackend bk_%s_%s-%s\n", mapping.stack, mapping.service, mapping.port)
		file.WriteString(line)
		fmt.Printf(line)
		if mapping.mode == "tcp" {
			if dnsResolved {
				line1 := fmt.Sprintf("    mode tcp\n    server %s_1 %s_%s:%s resolvers docker resolve-prefer ipv4\n", mapping.service, mapping.stack, mapping.service, mapping.port)
				file.WriteString(line1)
				fmt.Printf(line1)
			} else {
				line1 := "    #dns name not resolved\n"
				file.WriteString(line1)
				fmt.Printf(line1)
				line2 := fmt.Sprintf("#mode tcp\n#    server %s_1 %s_%s:%s resolvers docker resolve-prefer ipv4\n", mapping.service, mapping.stack, mapping.service, mapping.port)
				file.WriteString(line2)
				fmt.Printf(line2)
				app.addDNSNameInRetryList(mapping.service)
			}
		} else {

			//if dns name is not resolved haproxy (v1.6) won't start or accept the new configuration so server is disabled
			//to be removed when haproxy will fixe this bug
			if dnsResolved {
				line1 := fmt.Sprintf("    server %s_1 %s_%s:%s resolvers docker resolve-prefer ipv4\n", mapping.service, mapping.stack, mapping.service, mapping.port)
				file.WriteString(line1)
				fmt.Printf(line1)
			} else {
				line1 := "    #dns name not resolved\n"
				file.WriteString(line1)
				fmt.Printf(line1)
				line2 := fmt.Sprintf("    #server %s_1 %s_%s:%s resolvers docker resolve-prefer ipv4\n", mapping.service, mapping.stack, mapping.service, mapping.port)
				file.WriteString(line2)
				fmt.Printf(line2)
				app.addDNSNameInRetryList(mapping.service)
			}
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
