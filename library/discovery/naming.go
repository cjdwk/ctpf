package discovery

import "github.com/oofpgDLD/ctpf/library/environment"


func ServerName(env *environment.Env) string{
	//TODO convert env to server name
	return "ctpf-comet-" + env.Host
}

func ServerName2(env *environment.Env) string{
	//TODO convert env to server name
	return "ctpf-logic-" + env.Host
}