// Remote implementation of injectable functions
package remote

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/rpc"

	"github.com/antha-lang/antha/inject"
)

// Run a function remotely
type Runner struct {
	Name    string             // Name of type on remote
	Address string             // Address of RPC server
	Runner  inject.TypedRunner // Signature of function
}

func (a *Runner) Run(ctx context.Context, value inject.Value) (inject.Value, error) {
	client, err := rpc.DialHTTP("tcp", a.Address)
	if err != nil {
		return nil, err
	}

	in := a.Runner.Input()
	if err := inject.AssignableTo(value, &in); err != nil {
		return nil, fmt.Errorf("input value not assignable to %T: %s", in, err)
	}

	reply := inject.MakeValue(a.Runner.Output())
	if err := client.Call(fmt.Sprintf("%s.Call", a.Name), value, &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

func (a *Runner) Input() interface{} {
	return a.Runner.Input()
}

func (a *Runner) Output() interface{} {
	return a.Runner.Output()
}

type Server struct {
	Name        string                 // Name of type to use
	Listener    net.Listener           // Connections to listen to
	Runner      inject.TypedRunner     // Function to call
	MakeContext func() context.Context // Contexts for inject.Runner.Run()
}

// Wrapper for Server that only exposes RPC methods, so rpc.Register() doesn't
// warn about unmatched methods.
type S struct {
	Server *Server
}

// Function that will be called remotely
func (a *S) Call(value inject.Value, reply *inject.Value) error {
	ctx := context.Background()
	if a.Server.MakeContext != nil {
		ctx = a.Server.MakeContext()
	}
	r, err := a.Server.Runner.Run(ctx, value)
	if err != nil {
		return err
	}
	*reply = r
	return nil
}

func (a *Server) Serve() error {
	s := rpc.NewServer()
	s.RegisterName(a.Name, &S{Server: a})
	s.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	return http.Serve(a.Listener, s)
}
