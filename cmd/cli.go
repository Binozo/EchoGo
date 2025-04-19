package main

import (
	"context"
	"errors"
	"github.com/Binozo/EchoGo/v2/pkg/echo"
	"github.com/urfave/cli/v3"
	"log"
	"os"
	"os/exec"
)

func main() {
	alexa, err := echo.New()
	if err != nil {
		log.Fatal(err)
	}

	cmd := &cli.Command{
		Name:                  "EchoGo",
		Usage:                 "Control your echo easily",
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			{
				Name:           "boot",
				Usage:          "Boot your echo",
				DefaultCommand: "boot preloader_no_hdr.bin",
				Action: func(ctx context.Context, command *cli.Command) error {
					args := command.Args()
					preloaderPath := echo.DefaultPreloaderPath
					if args.Len() > 0 {
						preloaderPath = args.First()
					}

					log.Println("Using", preloaderPath, "to boot echo. This may take a few seconds. Do not press any buttons on your echo.")
					if err := alexa.Boot(preloaderPath); err != nil {
						return err
					}
					log.Println("Your echo is ready. Happy hacking!")
					return nil
				},
			},
			{
				Name:  "shutdown",
				Usage: "Shutdown your echo",
				Action: func(ctx context.Context, command *cli.Command) error {
					return alexa.Shutdown()
				},
			},
			{
				Name:  "restart",
				Usage: "Restart your echo",
				Action: func(ctx context.Context, command *cli.Command) error {
					log.Println("Restarting your echo. Do not press any buttons on the echo.")
					if err := alexa.Restart(echo.DefaultPreloaderPath); err != nil {
						return err
					}

					log.Println("Your echo is ready. Happy hacking!")
					return nil
				},
			},
			{
				Name:  "compile",
				Usage: "Compile the server app for your echo",
				Action: func(ctx context.Context, command *cli.Command) error {
					return cliCompile()
				},
			},
			{
				Name:  "run",
				Usage: "Run the server app on your echo for debugging purposes",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "compile",
						Aliases: []string{"c"},
					},
				},
				Action: func(ctx context.Context, command *cli.Command) error {
					if len(command.Flags) > 0 {
						if command.Flags[0].Get() == true {
							if err := cliCompile(); err != nil {
								return err
							}
						}
					}

					log.Println("Running server")
					return alexa.Run(echo.DefaultServerPath)
				},
			},
			{
				Name:  "deploy",
				Usage: "Deploy the server app to your echo",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "compile",
						Aliases: []string{"c"},
					},
				},
				Action: func(ctx context.Context, command *cli.Command) error {
					serverPath := "build/server"
					args := command.Args()
					if args.Len() > 0 {
						serverPath = args.First()
					}

					if len(command.Flags) > 0 {
						if command.Flags[0].Get() == true {
							if err := cliCompile(); err != nil {
								return err
							}
						}
					}

					if _, err := os.Stat(serverPath); errors.Is(err, os.ErrNotExist) {
						return err
					}
					log.Println("Deploying server from", serverPath)

					return alexa.Deploy(serverPath)
				},
			},
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func cliCompile() error {
	cmd := exec.Command("bash", "compile.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Println("Successfully compiled echo server at build/server")
	log.Println("Deploy using the 'deploy' command")
	return nil
}
