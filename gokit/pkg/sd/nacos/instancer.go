package nacos

import (
    "github.com/go-kit/kit/log"
    kitsd "github.com/go-kit/kit/sd"
    "github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
    "github.com/ncraft-io/ncraft/gokit/pkg/sd/internal/instance"
)

// Instancer yields instances stored in a certain etcd keyspace. Any kind of
// change in that keyspace is watched and will update the Instancer's Instancers.
type Instancer struct {
	cache       *instance.Cache
	client      *Client
	serviceName string
	groupName   string
	clusters    []string
	logger      log.Logger
	quitChan    chan struct{}
}

// NewInstancer returns an etcd instancer. It will start watching the given
// prefix for changes, and update the subscribers.
func NewInstancer(c *Client, serviceName, groupName string, clusters []string, logger log.Logger) (*Instancer, error) {
	s := &Instancer{
		client:      c,
		serviceName: serviceName,
		groupName:   groupName,
		clusters:    clusters,
		cache:       instance.NewCache(),
		logger:      logger,
		quitChan:    make(chan struct{}),
	}

	res, err := s.client.GetInstance(s.serviceName)
	if err == nil {
		logs.Info("serviceName", s.serviceName, "instances", len(res))
	} else {
		logs.Info("serviceName", s.serviceName, "err", err)
	}
	s.cache.Update(kitsd.Event{Instances: res, Err: err})
	go s.loop()
	return s, nil
}

func (s *Instancer) loop() {
	ch := make(chan struct{})
	go s.client.WatchService(s.serviceName, s.groupName, s.clusters, ch)

	for {
		select {
		case <-ch:
			instances, err := s.client.GetInstance(s.serviceName)
			if err != nil {
				s.logger.Log("msg", "failed to retrieve entries", "err", err)
				s.cache.Update(kitsd.Event{Err: err})
				continue
			}
			s.cache.Update(kitsd.Event{Instances: instances})

		case <-s.quitChan:
			return
		}
	}
}

// Stop terminates the Instancer.
func (s *Instancer) Stop() {
	close(s.quitChan)
}

// Register implements Instancer.
func (s *Instancer) Register(ch chan<- kitsd.Event) {
	s.cache.Register(ch)
}

// Deregister implements Instancer.
func (s *Instancer) Deregister(ch chan<- kitsd.Event) {
	s.cache.Deregister(ch)
}
