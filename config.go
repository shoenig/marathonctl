// Author Seth Hoenig

package main

import (
	"errors"
	"flag"

	"github.com/shoenig/config"
)

// cli arguments override configuration file
func cliargs() (config, host, login string) {
	flag.StringVar(&config, "c", "", "config file")
	flag.StringVar(&host, "h", "", "marathon host with transport and port")
	flag.StringVar(&login, "u", "", "username and password")
	flag.Parse()
	return
}

func readConfigfile(filename string) (host, login string, e error) {
	c, e := config.ReadProperties(filename)
	if e != nil {
		return "", "", e
	}
	h := c.GetStringOr("marathon.host", "")
	u := c.GetStringOr("marathon.user", "")
	p := c.GetStringOr("marathon.password", "")

	l := ""
	if u != "" && p != "" {
		l = u + ":" + p
	}

	return h, l, nil
}

// todo(someday) read $HOME/.config/marathonctl/config
// Read -config file
// Then override with cli args
func Config() (string, string, error) {
	config, host, login := cliargs()
	if host != "" && login != "" {
		return host, login, nil
	}

	if config != "" {
		h, l, e := readConfigfile(config)
		if e != nil {
			return "", "", e
		}
		if host == "" {
			host = h
		}
		if login == "" {
			login = l
		}
	}

	if host == "" || login == "" {
		return "", "", errors.New("no host or login info provided")
	}

	return host, login, nil
}
