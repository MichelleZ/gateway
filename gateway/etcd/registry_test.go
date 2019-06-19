/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */

package etcd

import (
	"sync"
	"testing"
	"time"

	etcd3 "go.etcd.io/etcd/clientv3"

	"github.com/xuperchain/xuperunion/common/log"
	loggw "github.com/xuperchain/xuperunion/gateway/log"
)

func TestRegister(t *testing.T) {
	etcdConfig := etcd3.Config{
		Endpoints: []string{"http://10.100.78.75:2379"},
	}

	xlog, _ := log.OpenLog(loggw.CreateLog())

	registry, err := NewRegistry(
		Input{
			EtcdConfig: etcdConfig,
			Prefix:     "gateway/testRegistry", // real: /etchxchain/gateway/testRegistry/
			Addr:       "0.0.0.0:11011",
			TTL:        5,
			Xlog:       xlog,
		})

	if err != nil {
		t.Error("Register failed")
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		registry.Register()
		wg.Done()
	}()

	registry.getInfo()

	time.Sleep(5 * time.Second)
	t.Log("Unregister")
	registry.UnRegister()

	wg.Wait()
}
