package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
	"strings"
)


const (
	PluginDirectoryFromRoot = "/wp-content/plugins/"
)

var (
	DefaultDirectory, _ = os.Getwd();
)

func main() {
	// Init cli app
	app := cli.NewApp()
	app.Name = "ED Composer to Git"
	app.Usage = "Removes composer from projects!"
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "dir, d",
			Value: DefaultDirectory,
			Usage: "Removes composer from the supplied `project directory`",
		},
	}
	app.Action = composerToGitCli

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}


func composerToGitCli(c *cli.Context) error {
	// Ensure we're in the right place..
	if _, err := os.Stat(c.String("dir")+"/composer.json"); os.IsNotExist(err) {
		return cli.NewExitError("Could not find composer.json file.", 69)
	}


	plugins, err := filepath.Glob(c.String("dir") + PluginDirectoryFromRoot + "/*")
	if err != nil || len(plugins) == 0 {
		return cli.NewExitError("Error finding plugins.", 69)
	}

	themes, err := filepath.Glob(c.String("dir") + ThemesDir+ "/*")
	if err != nil || len(themes) == 0 {
		return cli.NewExitError("Error finding themes.", 69)
	}


	cl := &client{
		log: log.New(os.Stdout, "composerToGit: ", 0),
		plugins: plugins,
		themes: themes,
		root: c.String("dir"),
	}

	cl.Convert()


	return nil
}




func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}