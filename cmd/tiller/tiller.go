package main

import (
	"fmt"
	"net"
	"os"

	"github.com/kubernetes/helm/cmd/tiller/environment"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// rootServer is the root gRPC server.
//
// Each gRPC service registers itself to this server during init().
var rootServer = grpc.NewServer()

// env is the default environment.
//
// Any changes to env should be done before rootServer.Serve() is called.
var env = environment.New()

var addr = ":44134"
var namespace = ""

const globalUsage = `The Kubernetes Helm server.

Tiller is the server for Helm. It provides in-cluster resource management.

By default, Tiller listens for gRPC connections on port 44134.
`

var rootCommand = &cobra.Command{
	Use:   "tiller",
	Short: "The Kubernetes Helm server.",
	Long:  globalUsage,
	Run:   start,
}

func main() {
	pf := rootCommand.PersistentFlags()
	pf.StringVarP(&addr, "listen", "l", ":44134", "The address:port to listen on")
	pf.StringVarP(&namespace, "namespace", "n", "", "The namespace Tiller calls home")
	rootCommand.Execute()
}

func start(c *cobra.Command, args []string) {
	setNamespace()
	lstn, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server died: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Tiller is running on %s\n", addr)

	if err := rootServer.Serve(lstn); err != nil {
		fmt.Fprintf(os.Stderr, "Server died: %s\n", err)
		os.Exit(1)
	}
}

// setNamespace sets the namespace.
//
// It checks for the --namespace flag first, then checks the environment
// (set by Downward API), then goes to default.
func setNamespace() {
	if len(namespace) != 0 {
		fmt.Printf("Setting namespace to %q\n", namespace)
		srv.env.Namespace = namespace
	} else if ns := os.Getenv("DEFAULT_NAMESPACE"); len(ns) != 0 {
		fmt.Printf("Inhereting namespace %q from Downward API\n", ns)
		srv.env.Namespace = ns
	} else {
		fmt.Printf("Using default namespace %q\n", environment.DefaultNamespace)
		srv.env.Namespace = environment.DefaultNamespace
	}
}
