package client

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"

	"github.com/dullgiulio/ringio/agents"
	"github.com/dullgiulio/ringio/msg"
	"github.com/dullgiulio/ringio/onexit"
	"github.com/dullgiulio/ringio/pipe"
	"github.com/dullgiulio/ringio/server"
	"github.com/dullgiulio/ringio/utils"
)

func addErrorsAgentPipe(client *rpc.Client, filter *msg.Filter, response *server.RpcResp, pipeName string) {
	_addSinkAgentPipe(client, filter, response, pipeName, agents.AgentRoleErrors)
}

func addSinkAgentPipe(client *rpc.Client, filter *msg.Filter, response *server.RpcResp, pipeName string) {
	_addSinkAgentPipe(client, filter, response, pipeName, agents.AgentRoleSink)
}

func addLogAgentPipe(client *rpc.Client, response *server.RpcResp, pipeName string) {
	_addSinkAgentPipe(client, nil, response, pipeName, agents.AgentRoleLog)
}

func _addSinkAgentPipe(client *rpc.Client, filter *msg.Filter,
	response *server.RpcResp, pipeName string, role agents.AgentRole) {
	var id int

	p := pipe.New(pipeName)

	if err := client.Call("RpcServer.Add", &server.RpcReq{
		Agent: &agents.AgentDescr{
			Args: []string{pipeName},
			Meta: agents.AgentMetadata{Role: role, Filter: filter},
			Type: agents.AgentTypePipe,
		},
	}, &id); err != nil {
		utils.Fatal(err)
	}

	if err := p.OpenReadErr(); err != nil {
		utils.Fatal(fmt.Errorf("Couldn't open pipe for reading: %s", err))
	}

	p.Remove()

	onexit.Defer(func() {
		if err := client.Call("RpcServer.Stop", id, &response); err != nil {
			utils.Fatal(err)
		}
	})

	r := bufio.NewReader(p)

	if _, err := r.WriteTo(os.Stdout); err != nil {
		utils.Fatal(err)
	}
}

func addErrorsAgentCmd(client *rpc.Client, filter *msg.Filter, response *server.RpcResp, args []string) {
	_addSinkAgentCmd(client, filter, response, args, agents.AgentRoleErrors)
}

func addSinkAgentCmd(client *rpc.Client, filter *msg.Filter, response *server.RpcResp, args []string) {
	_addSinkAgentCmd(client, filter, response, args, agents.AgentRoleSink)
}

func _addSinkAgentCmd(client *rpc.Client, filter *msg.Filter, response *server.RpcResp, args []string, role agents.AgentRole) {
	var id int

	if err := client.Call("RpcServer.Add", &server.RpcReq{
		Agent: &agents.AgentDescr{
			Args: args,
			Meta: agents.AgentMetadata{Role: role, Filter: filter},
			Type: agents.AgentTypeCmd,
		},
	}, &id); err != nil {
		utils.Fatal(err)
	}

    fmt.Printf("Added agent %%%d\n", id)
}
