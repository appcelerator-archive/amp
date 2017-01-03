package main

import (
	"fmt"
	ampClient "github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/cmd/adm-server/servergrpc"
	dockerclient "github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

const (
	//ClearScreen ANSI Escape code
	ClearScreen = "\033[2J\033[0;0H"
	//MoveCursorHome ANSI Escape code
	MoveCursorHome = "\033[0;0H"
)

type clusterClient struct {
	id            string
	client        servergrpc.ClusterServerServiceClient
	stream        servergrpc.ClusterServerService_GetClientStreamClient
	conn          *grpc.ClientConn
	ctx           context.Context
	verbose       bool
	silence       bool
	debug         bool
	nodeName      string
	nodeHost      string
	configuration *ampClient.Configuration
	clusterLoader *ClusterLoader
	recvChan      chan *servergrpc.ClientMes
	printColor    [7]*color.Color
	fcolTitle     func(...interface{}) string
	fcolLines     func(...interface{}) string
	// bootstrap properties
	dockerClient *dockerclient.Client
}

var currentColorTheme = "default"
var (
	colRegular = 0
	colInfo    = 1
	colWarn    = 2
	colError   = 3
	colSuccess = 4
	colUser    = 5
	colDebug   = 6
)

var (
	//RootCmd root command
	RootCmd = &cobra.Command{
		Use:   `amp-cluster`,
		Short: "amp-cluster",
		Run: func(cmd *cobra.Command, args []string) {
			//fmt.Println(cmd.UsageString())
		},
	}
)

func (g *clusterClient) init() error {
	g.setColors()
	g.clusterLoader = &ClusterLoader{}
	g.ctx = context.Background()
	g.recvChan = make(chan *servergrpc.ClientMes)
	return nil
}

func (g *clusterClient) initConfiguration(configFile string, serverAddr string) {
	g.configuration = &ampClient.Configuration{}
	InitConfig(g, configFile, g.configuration, g.debug, serverAddr)
	g.setColors()
}

func (g *clusterClient) initConnection() error {
	//---Treatement only for local usage
	if g.isLocalhostServer() {
		if err := g.clusterLoader.init(g, ""); err != nil {
			g.fatalc("init error: %v\n", err)
		}
		if !g.clusterLoader.isServiceRunning("adm-server") {
			g.fatalc("Cluster services are not available locally\n")
		}
	}
	//---
	if err := g.connectServer(); err != nil {
		g.fatalc("Connection to cluster server error: %v\n", err)
	}
	if err := g.startServerReader(); err != nil {
		g.fatalc("Connection to cluster server stream error: %v\n", err)
	}
	return nil
}

func (g *clusterClient) connectServer() error {
	conn, err := grpc.Dial(g.configuration.AdminServerAddress,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second*20))
	if err != nil {
		return err
	}
	g.conn = conn
	g.client = servergrpc.NewClusterServerServiceClient(conn)
	return nil
}

func (g *clusterClient) createSendMessageNoAnswer(target string, functionName string, args ...string) error {
	mes := &servergrpc.ClientMes{} //TODO
	_, err := g.sendMessage(mes, true)
	return err
}

func (g *clusterClient) createSendMessage(target string, waitForAnswer bool, functionName string, args ...string) (*servergrpc.ClientMes, error) {
	mes := &servergrpc.ClientMes{} //TODO
	return g.sendMessage(mes, waitForAnswer)
}

func (g *clusterClient) sendMessage(mes *servergrpc.ClientMes, wait bool) (*servergrpc.ClientMes, error) {
	err := g.stream.Send(mes)
	if err != nil {
		return nil, err
	}
	g.printfc(colDebug, "Message sent: %v\n", mes)
	if wait {
		ret := <-client.recvChan
		return ret, nil
	}
	return nil, nil
}

func (g *clusterClient) getNextAnswer() *servergrpc.ClientMes {
	mes := <-g.recvChan
	return mes
}

func (g *clusterClient) startServerReader() error {
	stream, err := g.client.GetClientStream(g.ctx)
	if err != nil {
		return err
	}
	g.stream = stream
	g.printfc(colDebug, "Connected to adm-server, waiting server ack\n")
	mes, err := g.stream.Recv()
	if err != nil {
		return fmt.Errorf("Server stream error: %v", err)
	}
	if mes.Function != "ClientAck" {
		return fmt.Errorf("Server-client handcheck error: Receive bad message type")
	}
	client.id = mes.ClientId
	g.printfc(colDebug, "Server ack, client id=%s\n", client.id)
	go func() {
		for {
			mes, err := g.stream.Recv()
			if err == io.EOF {
				g.printf("Server stream EOF\n")
				close(g.recvChan)
				return
			}
			if err != nil {
				g.printf("Server stream error: %v\n", err)
				return
			}
			if mes.Function == "Print" {
				g.printfc(int(mes.Output.OutputType), "%s\n", mes.Output.Output)
			}
			//
		}
	}()
	return nil
}

func (g *clusterClient) isLocalhostServer() bool {
	if strings.HasPrefix(g.configuration.AdminServerAddress, "127.0.0.1") || strings.HasPrefix(g.configuration.AdminServerAddress, "localhost") {
		return true
	}
	return false
}

func (g *clusterClient) startLocalClusterService() {
	if err := g.clusterLoader.init(g, ""); err != nil {
		g.fatalc("%v\n", err)
	}
	if err := g.clusterLoader.startClusterServices(); err != nil {
		g.fatalc("%v\n", err)
	}
	time.Sleep(2 * time.Second)
}

func (g *clusterClient) setColors() {
	theme := ""
	if g.configuration != nil { //configuration is nil between client init and cli init during this short time, default theme is used.
		theme = g.configuration.CmdTheme
	}
	if theme == "dark" {
		g.printColor[0] = color.New(color.FgHiWhite)
		g.printColor[1] = color.New(color.FgHiBlack)
		g.printColor[2] = color.New(color.FgYellow)
		g.printColor[3] = color.New(color.FgRed)
		g.printColor[4] = color.New(color.FgGreen)
		g.printColor[5] = color.New(color.FgHiGreen)
		g.printColor[6] = color.New(color.FgHiBlack)
	} else {
		g.printColor[0] = color.New(color.FgMagenta)
		g.printColor[1] = color.New(color.FgHiBlack)
		g.printColor[2] = color.New(color.FgYellow)
		g.printColor[3] = color.New(color.FgRed)
		g.printColor[4] = color.New(color.FgGreen)
		g.printColor[5] = color.New(color.FgHiGreen)
		g.printColor[6] = color.New(color.FgHiBlack)
	}
	//add theme as you want.
	g.fcolTitle = g.printColor[colRegular].SprintFunc()
	g.fcolLines = g.printColor[colSuccess].SprintFunc()
}

func (g *clusterClient) printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (g *clusterClient) getColor(col int) *color.Color {
	if g.silence {
		return nil
	}
	colorp := g.printColor[0]
	if col > 0 && col < len(g.printColor) {
		colorp = g.printColor[col]
	}
	if !g.verbose && col == colInfo {
		return nil
	}
	if !g.debug && col == colDebug {
		return nil
	}
	return colorp
}

func (g *clusterClient) printfc(col int, format string, args ...interface{}) {
	colorp := g.getColor(col)
	if colorp != nil {
		colorp.Printf(format, args...)
	}
}

func (g *clusterClient) fatalc(format string, args ...interface{}) {
	g.printfc(colError, format, args...)
	os.Exit(1)
}

func (g *clusterClient) displayInOrder(title1 string, title2 string, lines []string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	sort.Strings(lines)
	if title1 != "" {
		fmt.Fprintln(w, g.fcolTitle(title1))
	}
	if title2 != "" {
		fmt.Fprintln(w, g.fcolTitle(title2))
	}
	for _, line := range lines {
		fmt.Fprintf(w, "%s\n", g.fcolLines(line))
	}
	w.Flush()
}

func (g *clusterClient) followClearScreen(follow bool) {
	if follow {
		fmt.Println(ClearScreen)
	}
}

func (g *clusterClient) followMoveCursorHome(follow bool) {
	if follow {
		fmt.Println(MoveCursorHome)
	}
}
