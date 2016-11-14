package build_test

//import (
//	"fmt"
//	"github.com/appcelerator/amp/api/rpc/build"
//	"github.com/appcelerator/amp/api/server"
//	"golang.org/x/net/context"
//	"google.golang.org/grpc"
//	"io"
//	"os"
//	"strconv"
//	"testing"
//)
//
//var (
//	ctx    context.Context
//	client build.AmpBuildClient
//)
//
//func TestMain(m *testing.M) {
//	_, conn := server.StartTestServer()
//	client = build.NewAmpBuildClient(conn)
//	ctx = context.Background()
//	os.Exit(m.Run())
//}
//
//func TestPing(t *testing.T) {
//	ping := build.Ping{
//		Message: "ping",
//	}
//	_, err := client.PingPong(ctx, &ping)
//	if err != nil {
//		t.Fatal(err)
//	}
//}
//
//func TestBadPing(t *testing.T) {
//	ping := build.Ping{
//		Message: "not ping",
//	}
//	_, err := client.PingPong(ctx, &ping)
//	if err == nil {
//		t.Fatalf("bad ping should have failed")
//	}
//}
//
//func TestCreateProject(t *testing.T) {
//	fmt.Println("amp build register amp/fake")
//	request := build.ProjectRequest{
//		Owner: "amp",
//		Name:  "fake",
//	}
//	project, err := client.CreateProject(ctx, &request)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println("registered https://build.amp.appcelerator.io/p/" + project.Owner + "/" + project.Name)
//}
//
//func TestListProjects(t *testing.T) {
//	fmt.Println("amp build listprojects")
//	query := build.ProjectQuery{}
//	projects, err := client.ListProjects(ctx, &query)
//	if err != nil {
//		t.Fatal(err)
//	}
//	for _, p := range projects.Projects {
//		fmt.Println("https://build.amp.appcelerator.io/p/"+p.Owner+"/"+p.Name, p.Status)
//	}
//}
//
//func TestListProjectsByOrg(t *testing.T) {
//	fmt.Println("amp build listprojects -o appcelerator")
//	query := build.ProjectQuery{
//		Organization: "appcelerator",
//	}
//	projects, err := client.ListProjects(ctx, &query)
//	if err != nil {
//		t.Fatal(err)
//	}
//	for _, p := range projects.Projects {
//		fmt.Println("https://build.amp.appcelerator.io/p/"+p.Owner+"/"+p.Name, p.Status)
//	}
//}
//
//func TestListProjectsByLatest(t *testing.T) {
//	fmt.Println("amp build listprojects -l")
//	query := build.ProjectQuery{
//		Latest: true,
//	}
//	projects, err := client.ListProjects(ctx, &query)
//	if err != nil {
//		t.Fatal(err)
//	}
//	if len(projects.Projects) != 1 {
//		t.Fatalf("should only return one project")
//	}
//	for _, p := range projects.Projects {
//		fmt.Println("https://build.amp.appcelerator.io/p/"+p.Owner+"/"+p.Name, p.Status)
//	}
//}
//
//func TestDeleteProject(t *testing.T) {
//	fmt.Println("amp build remove amp/fake")
//	request := build.ProjectRequest{
//		Owner: "amp",
//		Name:  "fake",
//	}
//	project, err := client.DeleteProject(ctx, &request)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println("removed https://build.amp.appcelerator.io/p/" + project.Owner + "/" + project.Name)
//}
//
//func TestListBuilds(t *testing.T) {
//	fmt.Println("amp build listbuilds appcelerator/amp")
//	request := build.ProjectRequest{
//		Owner: "appcelerator",
//		Name:  "amp",
//	}
//	builds, err := client.ListBuilds(ctx, &request)
//	if err != nil {
//		t.Fatal(err)
//	}
//	for i, b := range builds.Builds {
//		fmt.Println("https://build.amp.appcelerator.io/p/"+b.Owner+"/"+b.Name+"/"+strconv.Itoa(i), b.CommitMessage, b.Status)
//	}
//}
//
//func TestStreamLogs(t *testing.T) {
//	fmt.Println("amp build logs appcelerator/amp/c4015d02fbc60583a4cd82187eb99d3aac3b36e4")
//	request := build.BuildRequest{
//		Owner: "appcelerator",
//		Name:  "amp",
//		Sha:   "c4015d02fbc60583a4cd82187eb99d3aac3b36e4",
//	}
//	logs, err := client.BuildLog(ctx, &request)
//	if err != nil {
//		t.Fatal(err)
//	}
//	for {
//		log, err := logs.Recv()
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			if grpc.ErrorDesc(err) == "EOF" {
//				break
//			} else {
//				t.Fatal(err)
//			}
//		}
//		fmt.Print(log.Message)
//	}
//}
//
//func TestRebuild(t *testing.T) {
//	fmt.Println("amp build rebuild appcelerator/amp/c4015d02fbc60583a4cd82187eb99d3aac3b36e4")
//	request := build.BuildRequest{
//		Owner: "appcelerator",
//		Name:  "amp",
//		Sha:   "c4015d02fbc60583a4cd82187eb99d3aac3b36e4",
//	}
//	build, err := client.Rebuild(ctx, &request)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println("rebuilding https://build.amp.appcelerator.io/p/" + build.Owner + "/" + build.Name + "/" + build.Sha)
//}
