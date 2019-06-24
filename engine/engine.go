package engine

import (
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-auth/api"
	"github.com/nuts-foundation/nuts-auth/pkg"
	nutsGo "github.com/nuts-foundation/nuts-go/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewAuthEngine() *nutsGo.Engine {

	authBackend := pkg.AuthInstance()

	return &nutsGo.Engine{
		Cmd: cmd(),
		Config: &authBackend.Config,
		ConfigKey: "auth",
		Configure: authBackend.Configure,
		FlagSet: flagSet(),
		Name: "Auth",
		Routes: func(router runtime.EchoRouter) {
			api.RegisterHandlers(router, &api.ApiWrapper{Auth: authBackend})
		},
	}
}


func cmd() *cobra.Command {

	cmd := &cobra.Command{
		Use: "auth",
		Short: "commands related to authentication",
	}

	cmd.AddCommand(&cobra.Command{
		Use: "server",
		Short: "Run standalone auth server",
		Run: func(cmd *cobra.Command, args []string) {
			authEngine := pkg.AuthInstance()
			echoServer := echo.New()
			echoServer.HideBanner = true
			api.RegisterHandlers(echoServer, &api.ApiWrapper{Auth:authEngine})
			logrus.Fatal(echoServer.Start(authEngine.Config.Address))
		},
	})

	return cmd
}

func flagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("auth", pflag.ContinueOnError)

	flags.String(pkg.ConfAddress, "localhost:1323", "Interface and port for http server to bind to")

	// TODO: add all the global auth command flags

	return flags
}