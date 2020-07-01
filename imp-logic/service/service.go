package service

import (
	"context"
	"github.com/bilibili/discovery/naming"
	"github.com/oofpgDLD/ctpf/imp-logic/conf"
	"github.com/oofpgDLD/ctpf/imp-logic/dao"
)

// Logic struct
type Service struct {
	c   *conf.Config
	dao *dao.Dao
	// online
	totalIPs   int64
	totalConns int64
	roomCount  map[string]int32
	// load balancer
	nodes        []*naming.Instance
	regions      map[string]string // province -> region
}

// New init
func New(c *conf.Config) *Service {
	s := &Service{
		c:            c,
		dao:          dao.New(c),
		regions:      make(map[string]string),
	}
	return s
}

// Ping ping resources is ok.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close close resources.
func (s *Service) Close() {
	s.dao.Close()
}