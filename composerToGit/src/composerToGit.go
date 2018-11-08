package main

import (
	"fmt"
	copy2 "github.com/otiai10/copy"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	PluginDir = "/wp-content/plugins/"
	ThemesDir = "/wp-content/themes/"
)

type client struct {
	log *log.Logger
	root string // Project root
	plugins []string
	themes []string
	pushToGit bool // Should push to git?
}

// ComposerToGit converts project away from composer
func (c *client) Convert() {
	// Remove composer.json and composer.lock files.. ignore errors as files are already verified
	c.deleteFile(c.root + "/composer.json")

	c.deleteFile(c.root + "/composer.lock")


	// Delete plugins .git
	for _, plugin := range c.plugins {
		ppath, _  := os.Lstat(plugin);
		ps := strings.Split(plugin, "/")
		pluginName := ps[len(ps)-1]


		switch mode := ppath.Mode(); {
		case mode.IsDir(): // Is regular plugin
			err := c.deleteFile(plugin + "/.git")
			if err != nil {
				c.log.Println(err)
			}

		case mode&os.ModeSymlink != 0: // Is symlink
			// Get symlink src
			l, err := os.Readlink(plugin)
			if err != nil {
				c.log.Printf("Error: %s", err)
			}
			// Remove symlink
			err = os.RemoveAll(plugin)
			if err != nil {
				c.log.Printf("Error: %s", err)
			}
			// Copy from symlink to dest
			err = copy2.Copy(l, c.root+PluginDir+pluginName)
			if err != nil {
				c.log.Printf("Error: ", err)
			}
			err = c.deleteFile(plugin + "/.git")
			if err != nil {
				c.log.Println(err)
			}
		}
	}


	// Delete themes .git
	for _, theme := range c.themes {
		ppath, _  := os.Lstat(theme);
		ts := strings.Split(theme, "/")
		themeName := ts[len(ts)-1]


		switch mode := ppath.Mode(); {
		case mode.IsDir(): // Is regular plugin
			err := c.deleteFile(theme+ "/.git")
			if err != nil {
				c.log.Println("Error: ",err)
			}

		case mode&os.ModeSymlink != 0: // Is symlink
			// Get symlink src
			l, err := os.Readlink(theme)
			if err != nil {
				c.log.Printf("Error: ", err)
			}
			// Remove symlink
			err = os.Remove(theme)
			if err != nil {
				c.log.Printf("Error: ", err)
			}
			// Copy from symlink to dest
			err = copy2.Copy(l, c.root+ThemesDir+themeName)
			if err != nil {
				c.log.Printf("Error: ", err)
			}
		}
	}


	// Remove themes and plugins from .gitignore
	file, err := ioutil.ReadFile(c.root+"/.gitignore")
	if err != nil {
		c.log.Fatal(err)
	}

	ignores := strings.Split(string(file), "\n")
	for i, ignore := range ignores{
		if strings.Contains(ignore, "themes") || strings.Contains(ignore, "plugins") {
			ignores[i] = ""
		}else {
			ignores[i] = ignore
		}
	}
	output := strings.Join(ignores, "\n")
	err = ioutil.WriteFile(c.root+"/.gitignore", []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}


	fmt.Println("-------------------------------------------")
	c.pushToGit = askForConfirmation("Commit changes?")


	if(c.pushToGit) {
		cmd := exec.Command("git", "rm", "-r", "--cached", ".")
		cmd.Dir = c.root
		out, err := cmd.Output()
		if err != nil {
			c.log.Println("Error: %s", err)
		}
		fmt.Printf("%s\n\n",out)

		cmd = exec.Command("git", "add", ".")
		cmd.Dir = c.root
		out, err = cmd.Output()
		if err != nil {
			c.log.Println("Error: %s", err)
		}
		fmt.Printf("%s\n\n",out)

		cmd = exec.Command("git", "commit", "-am", "composerToGit: convert project to git from composer")
		cmd.Dir = c.root
		out, err = cmd.Output()
		if err != nil {
			c.log.Println("Error: %s", err)
		}
		fmt.Printf("%s\n\n",out)
	}
	fmt.Println("-------------------------------------------")
	fmt.Println("Finished conversion")

}





func (c *client) deleteFile(f string) error{
	if v := verifyFile(f); v != false {
		err := os.RemoveAll(f)
		if err != nil {
			return err
		}
		c.log.Printf("Deleting %s", f)
		return nil
	}
	return fmt.Errorf("%s not found, continuing..", f)
}

// verifyFile verifies a file exists..
func verifyFile(f string) bool {
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		return true
	}
	return false
}

//exitOnError exits on an error
func exitOnError(err error) {
	if err != nil {
		fmt.Printf("Error -> %s", err)
		os.Exit(69)
	}
}

