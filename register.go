package srvreg

import (
	"context"
	"fmt"
	elog "github.com/EntropyPool/entropy-logger"
	clientv3 "github.com/coreos/etcd/clientv3"
	"time"
)

var cli *clientv3.Client = nil

var etcdEndpoints = []string{"etcd.npool.com:2379"}

const (
	ActionPut    = "put"
	ActionDelete = "delete"
	LeaseExpire  = 30  //租约有效期，单位s
	EtcdPrefix   = "root/"
)

type Event struct {
	Key   string
	Val   string
	Event string
}

func init() {
	var err error
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if nil != err {
		elog.Fatalf(elog.Fields{}, "fail to create etcd client")
	}
}

func Register(k, v string) error {
	_, err := cli.Put(context.TODO(), EtcdPrefix + k, v)
	if nil != err {
		return err
	}
	return nil
}

func Query(k string) ([]string, error) {
	resp, err := cli.Get(context.TODO(), EtcdPrefix + k)
	if nil != err {
		return nil, err
	}
	vals := make([]string, 0)
	for _, kv := range resp.Kvs {
		if string(kv.Key) == (EtcdPrefix + k) {
			vals = append(vals, string(kv.Value))
		}
	}

	if 0 == len(vals) {
		return nil, fmt.Errorf("miss value for %v", k)
	}

	return vals, nil
}

func eventFromEtcdEvent(ev *clientv3.Event) (*Event, error) {
	var event Event

	if "PUT" == ev.Type.String() {
		event.Event = ActionPut
	} else if "DELETE" == ev.Type.String() {
		event.Event = ActionDelete
	} else {
		return nil, fmt.Errorf("unknown event %v", ev.Type.String())
	}

	event.Key = string(ev.Kv.Key)
	event.Val = string(ev.Kv.Value)

	return &event, nil
}

func watchHandler(ch clientv3.WatchChan, cb func(ev Event)) {
	for resp := range ch {
		for _, ev := range resp.Events {
			event, err := eventFromEtcdEvent(ev)
			if nil != err {
				elog.Errorf(elog.Fields{}, "fail: %v", err)
				continue
			}
			cb(*event)
		}
	}
}

func Watch(k string, cb func(ev Event)) {
	k = EtcdPrefix + k
	if nil == cb {
		elog.Fatalf(elog.Fields{}, "miss key watch handler for %v", k)
	}
	ch := cli.Watch(context.TODO(), k)
	go watchHandler(ch, cb)
}

func HeartBeat() {
	for {
		_, err := cli.Lease.Grant(context.TODO(), LeaseExpire)
		if err != nil {
			elog.Errorf(elog.Fields{}, "lease grant failed: " + err.Error())
		}

		time.Sleep(time.Second * 3)
	}
}
