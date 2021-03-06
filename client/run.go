package client

import (
	"flag"
	"net/rpc"

	"github.com/dullgiulio/ringio/server"
	"github.com/dullgiulio/ringio/utils"
)

type CommandRun struct {
	client   *rpc.Client
	response *server.RPCResp
}

func NewCommandRun() *CommandRun {
	return &CommandRun{
		response: new(server.RPCResp),
	}
}

func (c *CommandRun) Help() string {
	return `Run all processes`
}

func (c *CommandRun) Init(fs *flag.FlagSet) bool {
	// nothing to do yet.
	return false
}

func (c *CommandRun) Run(cli *Cli) error {
	c.client = cli.GetClient()

	if err := c.client.Call("RPCServer.Run", &server.RPCReq{}, &c.response); err != nil {
		utils.Fatal(err)
	}

	return nil
}
