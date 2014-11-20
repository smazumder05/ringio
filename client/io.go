package client

import (
	"flag"
	"net/rpc"

	"bitbucket.org/dullgiulio/ringio/server"
	"bitbucket.org/dullgiulio/ringio/utils"
)

type CommandIO struct {
	client   *rpc.Client
	response *server.RpcResp
}

func NewCommandIO() *CommandIO {
	return &CommandIO{
		response: new(server.RpcResp),
	}
}

func (c *CommandIO) Help() string {
	return `List all agents`
}

func (c *CommandIO) Init(fs *flag.FlagSet) error {
	// nothing to do yet.
	return nil
}

func (c *CommandIO) Run(cli *Cli) error {
	if client, err := rpc.Dial("unix", utils.FileInDotpath(cli.Session)); err != nil {
		utils.Fatal(err)
	} else {
		c.client = client
	}

	go addSourceAgentPipe(c.client, c.response, utils.GetRandomDotfile())
	addSinkAgentPipe(c.client, c.response, utils.GetRandomDotfile())

	return nil
}
