package core

import (
	"fmt"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/coreos/etcd/clientv3"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"log"
	"time"
)

const stackRootKey = "amp/stacks"

//ETCDClient etcd struct
type ETCDClient struct {
	client *clientv3.Client
}

var etcdClient ETCDClient

func (inst *ETCDClient) init() error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.etcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Connected to etcd: %v\n", conf.etcdEndpoints)
	inst.client = cli
	return err
}

//Close close ETCD client
func (inst *ETCDClient) Close() error {
	if err := inst.client.Close(); err != nil {
		return err
	}
	inst.client = nil
	return nil
}

func (inst *ETCDClient) watchForServicesUpdate() {
	watchKeys := stackRootKey
	fmt.Println("Waiting for update on ", watchKeys)
	rch := inst.client.Watch(context.Background(), watchKeys, clientv3.WithPrefix())
	wresp := <-rch
	for _, ev := range wresp.Events {
		fmt.Printf("Key updated: %q\n", ev.Kv.Key)
	}
	time.Sleep(5 * time.Second)
	haproxy.loadTry = 0
	haproxy.updateConfiguration(true)
}

func (inst *ETCDClient) getAllMappings() ([]*publicMapping, error) {
	log.Printf("getAllMappings\n")
	resp, err := inst.client.Get(context.Background(), stackRootKey, clientv3.WithPrefix())
	if err != nil {
		log.Printf("getAllMappings error: %v\n", err)
		return nil, err
	}
	log.Printf("getAllMappings result: %+v\n", resp)
	mappingList := []*publicMapping{}
	if len(resp.Kvs) == 0 {
		return mappingList, nil
	}
	for _, kvs := range resp.Kvs {
		mappingGet := &stack.StackMapping{}
		if err := proto.Unmarshal(kvs.Value, mappingGet); err != nil {
			log.Printf("Error unmarcharling mapping: %v\n", err)
		} else {
			fmt.Printf("process mapping:%v\n", mappingGet)
			if data, err := stack.EvalMappingString(mappingGet.Mapping); err != nil {
				log.Printf("Mapping error: %v\n", err)
			} else {
				mapping := &publicMapping{
					stack:   mappingGet.StackName,
					service: mappingGet.ServiceName,
					label:   data[0],
					port:    data[1],
					mode:    data[2],
					portTo:  data[3],
				}
				mappingList = append(mappingList, mapping)
			}
		}
	}
	return mappingList, nil
}
