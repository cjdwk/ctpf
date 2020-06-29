module github.com/oofpgDLD/ctpf

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Terry-Mao/goim v0.0.0-00010101000000-000000000000
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.2
	github.com/inconshreveable/log15 v0.0.0-20200109203555-b30bc20e4fd1
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins v1.5.1
	github.com/uber/jaeger-client-go v2.23.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	github.com/zhenjl/cityhash v0.0.0-20131128155616-cdd6a94144ab
	google.golang.org/grpc v1.25.1
)

replace github.com/Terry-Mao/goim => ../../Terry-Mao/goim //v1.0.0
