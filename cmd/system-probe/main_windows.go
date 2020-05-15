// +build windows

package main

import (
	"flag"
)

func main() {
	flag.StringVar(&opts.configPath, "config", "c:\\programdata\\datadog\\system-probe.yaml", "Path to system-probe config formatted as YAML")
	flag.StringVar(&opts.pidFilePath, "pid", "", "Path to set pidfile for process")
	flag.BoolVar(&opts.version, "version", false, "Print the version and exit")
	flag.Parse()

	runAgent()
}
