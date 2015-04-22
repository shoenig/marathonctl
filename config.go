// Author Seth Hoenig

package main

import (
	"errors"
	"flag"
	"os"

	"github.com/shoenig/config"
)

// cli arguments override configuration file
func cliargs() (config, host, login, format string) {
	flag.StringVar(&config, "c", "", "config file")
	flag.StringVar(&host, "h", "", "marathon host with transport and port")
	flag.StringVar(&login, "u", "", "username and password")
	flag.StringVar(&format, "f", "", "output format")
	flag.Parse()
	return
}

func readConfigfile(filename string) (host, login, format string, e error) {
	c, e := config.ReadProperties(filename)
	if e != nil {
		return "", "", "", e
	}
	h := c.GetStringOr("marathon.host", "")
	u := c.GetStringOr("marathon.user", "")
	p := c.GetStringOr("marathon.password", "")
	f := c.GetStringOr("marathon.format", "")

	l := ""
	if u != "" && p != "" {
		l = u + ":" + p
	}

	return h, l, f, nil
}

func configFile() string {
	configLocations := [2]string{os.Getenv("HOME") + "/.config/marathonctl/config", "/etc/marathonctl"}
	for _, location := range configLocations {
		if _, err := os.Stat(location); err == nil {
			return location
		}
	}
	return ""
}

// todo(someday) read $HOME/.config/marathonctl/config
// Read -config file
// Then override with cli args
func Config() (string, string, string, error) {
	config, host, login, format := cliargs()

	if host != "" && login != "" {
		return host, login, format, nil
	}

	if config == "" {
		config = configFile()
	}

	if config != "" {
		h, l, f, e := readConfigfile(config)
		if e != nil {
			return "", "", "", e
		}
		if host == "" {
			host = h
		}
		if login == "" {
			login = l
		}
		if format == "" {
			format = f
			if format == "" {
				format = "human"
			}
		}
	}

	if host == "" {
		return "", "", "", errors.New("no host info provided")
	}

	return host, login, format, nil
}
