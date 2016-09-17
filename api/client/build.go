package client

// import (
// 	"fmt"
// 	"github.com/appcelerator/amp/api/rpc/build"
// 	"golang.org/x/net/context"
// 	"strings"
// )

// func createProjectRequest(repo string) (request *build.ProjectRequest, err error) {
// 	split := strings.Split(repo, "/")
// 	if len(split) != 2 {
// 		return nil, fmt.Errorf("invalid repo %q", split)
// 	}
// 	owner := split[0]
// 	name := split[1]
// 	request = &build.ProjectRequest{
// 		Owner: owner,
// 		Name:  name,
// 	}
// 	return
// }

// func createBuildRequest(buildid string) (request *build.BuildRequest, err error) {
// 	split := strings.Split(buildid, "/")
// 	if len(split) != 3 {
// 		return nil, fmt.Errorf("invalid build %q", split)
// 	}
// 	owner := split[0]
// 	name := split[1]
// 	sha := split[2]
// 	request = &build.BuildRequest{
// 		Owner: owner,
// 		Name:  name,
// 		Sha:   sha,
// 	}
// 	return
// }

// // RegisterProject registers a project through the amplifier proxy to the build service
// func (a *AMP) RegisterProject(repo string) (*build.Project, error) {
// 	client := build.NewAmpBuildClient(a.Conn)
// 	ctx := context.Background()
// 	request, err := createProjectRequest(repo)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return client.CreateProject(ctx, request)
// }

// // RemoveProject removes a project through the amplifier proxy to the build service
// func (a *AMP) RemoveProject(repo string) (*build.Project, error) {
// 	client := build.NewAmpBuildClient(a.Conn)
// 	ctx := context.Background()
// 	request, err := createProjectRequest(repo)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return client.DeleteProject(ctx, request)
// }

// // ListProjects lists projects through the amplifier proxy to the build service
// func (a *AMP) ListProjects(organization string, latest bool) (*build.ProjectList, error) {
// 	client := build.NewAmpBuildClient(a.Conn)
// 	ctx := context.Background()
// 	query := build.ProjectQuery{
// 		Organization: organization,
// 		Latest:       latest,
// 	}
// 	return client.ListProjects(ctx, &query)
// }

// // ListBuilds lists builds through the amplifier proxy to the build service
// func (a *AMP) ListBuilds(repo string, latest bool) (*build.BuildList, error) {
// 	client := build.NewAmpBuildClient(a.Conn)
// 	ctx := context.Background()
// 	request, err := createProjectRequest(repo)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return client.ListBuilds(ctx, request)
// }

// // BuildLog returns a log stream through the amplifier proxy to the build service
// func (a *AMP) BuildLog(buildid string) (build.AmpBuild_BuildLogClient, error) {
// 	client := build.NewAmpBuildClient(a.Conn)
// 	ctx := context.Background()
// 	request, err := createBuildRequest(buildid)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return client.BuildLog(ctx, request)
// }

// // Rebuild triggers a rebuild through the amplifier proxy to the build service
// func (a *AMP) Rebuild(buildid string) (*build.Build, error) {
// 	client := build.NewAmpBuildClient(a.Conn)
// 	ctx := context.Background()
// 	request, err := createBuildRequest(buildid)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return client.Rebuild(ctx, request)
// }
