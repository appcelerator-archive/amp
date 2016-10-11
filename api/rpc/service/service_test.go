package service_test

import (
        "fmt"
        "os"
        "testing"
        "time"

        "github.com/appcelerator/amp/api/rpc/service"
        "github.com/appcelerator/amp/api/server"
        "github.com/stretchr/testify/assert"
        "golang.org/x/net/context"
        "google.golang.org/grpc"
)

var (
        config           server.Config
        port             string
        etcdEndpoints    string
        elasticsearchURL string
        kafkaURL         string
        influxURL        string
        dockerURL        string
        dockerVersion    string
        client           service.ServiceClient
        ctx              context.Context
)


        var service1 = service.ServiceSpec{
                Image: "appcelerator/pinger",
                Mode:  &service.ServiceSpec_Replicated{
                        Replicated: &service.ReplicatedService{Replicas: 2}, 
                },
        }

        var service2 = service.ServiceSpec{
                Image: "appcelerator/pinger",
                Labels: map[string]string{
                        "label1":"value1",
                        "label2":"value2",
                },

        }

        var service3 = service.ServiceSpec{
                Image: "appcelerator/pinger",
                ContainerLabels: map[string]string{
                        "label1":"value1",
                        "label2":"value2",
                },
        }

        var service4 = service.ServiceSpec{
                Image: "appcelerator/pinger",
                Env: []string{
                        "Var1=value1",
                        "Var2=value2",
                },
        }

        var service5 = service.ServiceSpec{
                Image: "appcelerator/pinger",
                PublishSpecs: []*service.PublishSpec{
                        {
                                PublishPort: 3001,
                                InternalPort: 3000,
                        },
                },
        }       

        var service6 = service.ServiceSpec{
                Image: "appcelerator/pinger",
                PublishSpecs: []*service.PublishSpec{
                        {
                                Name: "www",
                                InternalPort: 80,
                        },
                },
        }  

        var serviceList = []*service.ServiceSpec{
                &service1,
                &service2,
                &service3,
                &service4,
                &service5,
                &service6,
        }


func TestMain(m *testing.M) {
        ctx = context.Background()
       _, conn := server.StartTestServer()
        // there is no event when the server starts listening, so we just wait a second
        time.Sleep(1 * time.Second)
        conn, err := grpc.Dial("localhost:50101", grpc.WithInsecure())
        if err != nil {
                fmt.Println("connection failure")
                os.Exit(1)
        }
        client = service.NewServiceClient(conn)
        os.Exit(m.Run())
}

//Test two stacks life cycle in the same time
func TestServices(t *testing.T) {
        for _, serv := range serviceList{
                name := fmt.Sprintf("service-test-%d", time.Now().Unix())
                serv.Name = name
                respc, errc := client.Create(ctx, &service.ServiceCreateRequest{
                        ServiceSpec: serv,
                })
                if errc != nil {
                        t.Fatal(errc)
                }
                assert.NotEmpty(t, respc.Id, "returned service id should not be empty after create")
                respr, errr := client.Remove(ctx, &service.RemoveRequest{
                        Ident: respc.Id,
                })
                if errr != nil {
                        t.Fatal(errr)
                }
                assert.NotEmpty(t, respr.Ident, "returned service id should not be empty after remove")
        }

}
