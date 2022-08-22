package main

import (
	"math/rand"
	"os"
	"runtime/debug"
	"time"

	"github.com/urfave/cli"

	"github.com/free5gc/go-upf/internal/logger"
	upfapp "github.com/free5gc/go-upf/pkg/app"
	"github.com/free5gc/go-upf/pkg/factory"
	"github.com/ubrovc/util/version"
)

func main() {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.MainLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	app := cli.NewApp()
	app.Name = "upf"
	app.Usage = "5G User Plane Function (UPF)"
	app.Action = action
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
		cli.StringFlag{
			Name:  "log, l",
			Usage: "Output NF log to `FILE`",
		},
		cli.StringFlag{
			Name:  "log5gc, lc",
			Usage: "Output free5gc log to `FILE`",
		},
	}

	rand.Seed(time.Now().UnixNano())

	if err := app.Run(os.Args); err != nil {
		logger.MainLog.Errorf("UPF Cli Run Error: %v", err)
	}
}

func action(cliCtx *cli.Context) error {
	err := logger.LogFileHook(cliCtx.String("log"), cliCtx.String("log5gc"))
	if err != nil {
		return err
	}

	logger.MainLog.Infoln("UPF version: ", version.GetVersion())

	cfg, err := factory.ReadConfig(cliCtx.String("config"))
	if err != nil {
		return err
	}

	upf, err := upfapp.NewUpf(cfg)
	if err != nil {
		return err
	}

	if err := upf.Run(); err != nil {
		return err
	}

	return nil
}
