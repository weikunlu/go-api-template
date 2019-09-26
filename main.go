package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
	"github.com/weikunlu/go-api-template/cmd"
	"github.com/weikunlu/go-api-template/cmd/database"
	"github.com/weikunlu/go-api-template/config"
	"os"
)

// Variables to identify the build
var (
	Version   string
	Build     string
	BuildDate string

	cliApp *cli.App
)

func init() {
	// Initialise a CLI app
	cliApp = cli.NewApp()
	cliApp.Name = "App Phone Server"
	cliApp.Usage = "appphone-service"
	cliApp.Author = "Lucas Lu"
	cliApp.Email = "lucas.lu@kkday.com"
	cliApp.Version = fmt.Sprintf("%s-%s-%s", Version, Build, BuildDate)

}

func main() {
	config.SetBuildConfig(Version, Build, BuildDate)

	if os.Getenv("DEVELOP") == "true" {
		cmd.RunServer()
		return
	}

	// Set the CLI app commands
	cliApp.Commands = []cli.Command{
		{
			Name:  "runserver",
			Usage: "run web server",
			Action: func(c *cli.Context) error {
				return cmd.RunServer()
			},
		},
		{
			Name:  "migrate",
			Usage: "run database migrate",
			Action: func(c *cli.Context) error {
				return database.RunMigration(c.Args().Get(0), c.Args().Get(1))
			},
		},
	}

	// Run the CLI app
	if err := cliApp.Run(os.Args); err != nil {

		fmt.Printf("startup error: %v", err.Error())
	}
}
