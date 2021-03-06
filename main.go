package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BrosSquad/hosts/cmd"
	"github.com/BrosSquad/hosts/host"
	_ "github.com/jessevdk/go-flags"
)

const Localhost = "127.0.0.1"

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	options, err := cmd.NewOptions()
	handleError(err)

	switch options.Command() {
	case "add":
		if options.AddOptions.Ip == "" {
			options.AddOptions.Ip = Localhost
		}

		file, err := os.OpenFile(options.AddOptions.File, os.O_RDWR|os.O_APPEND, 0644)
		handleError(err)
		p := host.NewParser(file)
		handleError(p.Add(options.AddOptions.Host, options.AddOptions.Ip))
		fmt.Printf("New host added to file: %s %s\n", options.AddOptions.Host, options.AddOptions.Ip)
		handleError(file.Close())
	case "remove":
		tmp, err := ioutil.TempFile("", "hosts_copy")
		handleError(err)
		file, err := os.OpenFile(options.RemoveOptions.File, os.O_RDWR, 0644)
		handleError(err)
		p := host.NewParser(file)
		handleError(p.Remove(tmp, options.RemoveOptions.Host))
		fmt.Printf("Host removed from file: %s\n", options.RemoveOptions.Host)
		handleError(tmp.Close())
		handleError(p.Close())
	case "list":
		file, err := os.OpenFile(options.ListOptions.File, os.O_RDONLY, 0644)
		handleError(err)
		p := host.NewParser(file)
		handleError(p.List(func(host, ip string, isComment bool) error {
			fmt.Printf("Host: %s, IP: %s\n", host, ip)
			return nil
		}))
		handleError(p.Close())
	}
}
