package discover

import (
	"context"
	"encoding/json"
	client "github.com/coreos/etcd/clientv3"
	"log"
	"time"
)

const (
	//ETCD服务地址
	ENDPOINTS = "http://106.14.125.55:2379"
	//请求超时时间
	DIALTIMEOUT = 5
	//每次续租时间
	LEASEEXPIRE = 10
	//心跳周期
	CYCLE = 3
)

type Discover Struct {
	Client *client.Client
}

//服务注册
func Register(name string, iport string) {

}

//服务发现
func Query(name string) {

}

//服务监控
func (discover *Discover) Watch() {
	/**
	watcherCh := discover.Client.Watch(context.TODO(), "services", client.WithPrefix())
	for resp := range watcherCh {
		for _, ev := ranage resp.Events {
			key := string(ev.Kv.Key)
			switch ev.Type.String {
				case "PUT":
				case "GET":
				case "DELETE":
			}
		}
	}
	*/
}

//心跳
func (discover *Discover) heartBeat() {
	for {
		_, err := discover.Client.Lease.Grant(context.TODO(), LEASEEXPIRE)
		if err != nil {
			log.Fatalf("设置租约时间失败:%s\n", err.Error())
		}

		time.Sleep(CYCLE * time.Second)
	}
}

//创建Discover
func NewDiscover() (*Discover, error) {
	endPoints := []string{ENDPOINTS}
	cfg := client.Config{
		Endpoints: endPoints,
		DialTimeout: DIALTIMEOUT * time.Second,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Error("Error: cannot connect to etcd: ", err)

		return nil, err
	}

	discover := &Discover{
		Client: etcdClient,
	}

	return discover, nil
}















