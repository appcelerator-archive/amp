package build

// import (
// 	"fmt"
// 	"io"
// 	"os"

// 	"golang.org/x/net/context"
// 	"google.golang.org/grpc"
// )

// // Proxy is used to implement build.BuildServer
// type Proxy struct{}

// var (
// 	client AmpBuildClient
// )

// // Init ititializes the grpc connection to the build server
// func init() {
// 	conn, err := grpc.Dial("build.amp.appcelerator.io:50052", grpc.WithInsecure())
// 	if err != nil {
// 		fmt.Println("build connection failure", err)
// 		os.Exit(1)
// 	}
// 	client = NewAmpBuildClient(conn)
// }

// // PingPong forwards to the build PingPong service
// func (p *Proxy) PingPong(ctx context.Context, ping *Ping) (pong *Pong, err error) {
// 	return client.PingPong(ctx, ping)
// }

// // CreateProject fowards to the build CreateProject service
// func (p *Proxy) CreateProject(ctx context.Context, request *ProjectRequest) (project *Project, err error) {
// 	return client.CreateProject(ctx, request)
// }

// // ListProjects fowards to the build ListProjects service
// func (p *Proxy) ListProjects(ctx context.Context, query *ProjectQuery) (projects *ProjectList, err error) {
// 	return client.ListProjects(ctx, query)
// }

// // DeleteProject forwards to the build DeleteProject service
// func (p *Proxy) DeleteProject(ctx context.Context, request *ProjectRequest) (project *Project, err error) {
// 	return client.DeleteProject(ctx, request)
// }

// // ListBuilds forwards to the build ListBuilds service
// func (p *Proxy) ListBuilds(ctx context.Context, request *ProjectRequest) (builds *BuildList, err error) {
// 	return client.ListBuilds(ctx, request)
// }

// // BuildLog forwards to the build BuildLog service
// func (p *Proxy) BuildLog(request *BuildRequest, server AmpBuild_BuildLogServer) (err error) {
// 	var log *Log
// 	logs, err := client.BuildLog(context.Background(), request)
// 	if err != nil {
// 		return
// 	}
// 	for {
// 		log, err = logs.Recv()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return
// 		}
// 		err = server.Send(log)
// 		if err != nil {
// 			return
// 		}
// 	}
// 	return
// }

// // Rebuild forwards to the build Rebuild service
// func (p *Proxy) Rebuild(ctx context.Context, request *BuildRequest) (build *Build, err error) {
// 	return client.Rebuild(ctx, request)
// }
