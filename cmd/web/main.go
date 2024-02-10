package main

import (
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/tarro-dev/discord-oauth/internal/config"
	"github.com/tarro-dev/discord-oauth/internal/container"
	"github.com/tarro-dev/discord-oauth/pkg/option"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name:   "discord-oauth",
		Usage:  "A simple web server to authenticate with Discord",
		Action: run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func run(cliCtx *cli.Context) error {
	// Load the configuration
	config := config.NewConfig(getFlagString(cliCtx, "config"))
	ctx := config.Logger().WithContext(cliCtx.Context)

	if _, err := config.DiscordClientID(); err != nil {
		zerolog.Ctx(ctx).Fatal().Err(err).Msg("failed to get client ID")
	}

	if _, err := config.DiscordClientSecret(); err != nil {
		zerolog.Ctx(ctx).Fatal().Err(err).Msg("failed to get client secret")
	}

	// Create the container
	container := container.New(config)

	// Get the router
	router, err := container.Router(ctx)
	if err != nil {
		zerolog.Ctx(ctx).Fatal().Err(err).Msg("failed to get router")
	}

	// Start the server
	zerolog.Ctx(ctx).Info().Str("addr", config.Addr()).Msg("listening")
	if err := http.ListenAndServe(config.Addr(), router.Handler(ctx)); err != nil {
		zerolog.Ctx(ctx).Fatal().Err(err).Msg("failed to listen and serve")
	}

	return nil
}

func getFlagString(ctx *cli.Context, name string) option.Option[string] {
	if ctx.IsSet(name) {
		return option.Some(ctx.String(name))
	}
	return option.None[string]()
}
