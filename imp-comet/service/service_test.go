package service

import (
	"github.com/oofpgDLD/ctpf/imp-comet/conf"
	"testing"
	"os"
)

var (
	s *Service
)

func TestMain(m *testing.M) {
	s = New("test", conf.Default())
	os.Exit(m.Run())
}
