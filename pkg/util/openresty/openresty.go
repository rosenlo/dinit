package openresty

import (
	"fmt"

	consul "github.com/hashicorp/consul/api"
	"github.com/rosenlo/toolkits/log"
)

type NginxUpstream struct {
	SRVS     []*UpstreamNode `json:"srvs"`
	LBMethod string          `json:"lb_method"`
}

type UpstreamNode struct {
	Available  bool   `json:"available"`
	Weight     int    `json:"weight"`
	ServerType string `json:"server_type"`
	Host       string `json:"host"`
	Backup     bool   `json:"backup"`
	Port       int32  `json:"port"`
}

type Client struct {
	Cli  *consul.Client
	lock *consul.Lock
}

func NewClient(config *consul.Config) (*Client, error) {
	cli, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Client{
		Cli: cli,
	}, nil
}

func (c *Client) Lock(key string) error {
	lock, err := c.Cli.LockKey(key)
	if err != nil {
		return err
	}
	c.lock = lock

	stopCh := make(chan struct{})
	lockCh, err := lock.Lock(stopCh)
	if err != nil {
		return err
	}
	if lockCh == nil {
		return fmt.Errorf("Failed to get lockCh")
	}
	log.Info("Successfully obtained the lockCh")

	go func() {
		select {
		case <-lockCh:
		default:
			close(stopCh)
			log.Info("lockCh is closed")
		}
	}()

	return nil
}

func (c *Client) UnLock() {
	err := c.lock.Unlock()
	if err != nil {
		log.Warnf("Unlock raise an exception %v", err)
		return
	}
	log.Info("Unlock success")
}

func (c *Client) Put(key string, value []byte) error {
	log := log.WithField("key", key)
	// kvPair, err := c.Get(key)
	// if err != nil {
	// 	return fmt.Errorf("Failed to get the key %s", err)
	// }
	// log.Infof("kvPair %v", kvPair)

	kv := c.Cli.KV()

	writeMeta, err := kv.Put(
		&consul.KVPair{Key: key, Value: value},
		nil,
	)
	if err != nil {
		return fmt.Errorf("Failed to put the key due to %s", err)
	}
	log.WithField("latency", writeMeta.RequestTime).Info("Put the key success")

	return nil
}

func (c *Client) Get(key string) (*consul.KVPair, error) {
	kvPair, _, err := c.Cli.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}

	return kvPair, nil
}
