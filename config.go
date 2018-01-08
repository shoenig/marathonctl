// Author Seth Hoenig

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shoenig/config"
)

type flags struct {
	version    bool
	semver     bool
	configfile string
	host       string
	login      string
	format     string
	insecure   bool
}

// cli arguments override configuration file
func cliargs() flags {
	var f flags
	flag.BoolVar(&f.version, "v", false, "display version (git sha1) and exit")
	flag.BoolVar(&f.semver, "s", false, "display semversion and exit")
	flag.StringVar(&f.configfile, "c", "", "path to configfile")
	flag.StringVar(&f.host, "h", "", "override marathon host (with transport and port)")
	flag.StringVar(&f.login, "u", "", "override username and password")
	flag.StringVar(&f.format, "f", "", "override output format (raw, json, jsonpp)")
	flag.BoolVar(&f.insecure, "k", false, "insecure - do not verify certificate authority")
	flag.Parse()
	return f
}

func readConfigfile(filename string) (string, string, string, error) {
	props, err := config.ReadProperties(filename)
	if err != nil {
		return "", "", "", err
	}
	host := props.GetStringOr("marathon.host", "")
	user := props.GetStringOr("marathon.user", "")
	pass := props.GetStringOr("marathon.password", "")
	format := props.GetStringOr("marathon.format", "")

	login := ""
	if user != "" && pass != "" {
		login = user + ":" + pass
	}

	return host, login, format, nil
}

func findBestConfigfile() string {
	locations := [2]string{
		filepath.Join(os.Getenv("HOME"), ".config", "marathonctl", "config"),
		filepath.Join("etc", "marathonctl"),
	}

	for _, location := range locations {
		if _, err := os.Stat(location); err == nil {
			return location
		}
	}

	return ""
}

// loadConfig will parse the CLI flags.
// If --version or --semver are set, no further configuration
// is read. Otherwise, configuration is read from --configfile as
// specified, and then overridden with provided CLI flags.
func loadConfig() (flags, error) {
	f := cliargs()

	if f.version || f.semver {
		return f, nil
	}

	if f.host != "" && f.login != "" {
		return f, nil
	}

	if f.configfile == "" {
		f.configfile = findBestConfigfile()
	}

	if f.configfile != "" {
		fileHost, fileLogin, fileFormat, err := readConfigfile(f.configfile)
		if err != nil {
			return flags{}, fmt.Errorf("failed to read config file: %v", err)
		}

		if f.host == "" && fileHost != "" {
			f.host = fileHost
		}
		if f.login == "" && fileLogin != "" {
			f.login = fileLogin
		}
		if f.format == "" && fileFormat != "" {
			f.format = fileFormat
		}
	}

	return f, nil
}
