package main

import (
	"context"
	"fmt"
	log "github.com/go-pkgz/lgr"
	"github.com/umputun/go-flags"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var opts struct {
	Dbg bool `long:"dbg" env:"DEBUG" description:"debug mode"`
	Ftp struct {
		Enabled       bool          `long:"enabled" env:"ENABLED" description:"enable ftp watcher"`
		Ip            string        `long:"ip" env:"IP" default:"127.0.0.1" description:"ip address of ftp server"`
		User          string        `long:"user" env:"USER" default:"anonymous" description:"user name"`
		Password      string        `long:"password" env:"PASSWORD" default:"" description:"user password"`
		Dir           string        `long:"dir" env:"DIR" default:"/" description:"ftp directory for recursive read"`
		CheckInterval time.Duration `long:"interval" env:"INTERVAL" default:"3s" description:"ftp check interval"`
	} `group:"ftp" namespace:"ftp" env-namespace:"FTP"`
	File struct {
		Enabled       bool          `long:"enabled" env:"ENABLED" description:"enable file watcher"`
		Dir           string        `long:"dir" env:"DIR" default:"/tmp" description:"file directory for recursive read"`
		CheckInterval time.Duration `long:"interval" env:"INTERVAL" default:"3s" description:"file check interval"`
	} `group:"file" namespace:"file" env-namespace:"FILE"`
	Http struct {
		Enabled       bool          `long:"enabled" env:"ENABLED" description:"enable http watcher"`
		Url           string        `long:"url" env:"URL" default:"http://localhost:8080" description:"http url"`
		CheckInterval time.Duration `long:"interval" env:"INTERVAL" default:"3s" description:"http check interval"`
	} `group:"http" namespace:"http" env-namespace:"HTTP"`
}

var revision = "unknown"

func main() {
	fmt.Printf("detect_bot %s\n", revision)

	p := flags.NewParser(&opts, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag)
	p.SubcommandsOptional = true
	if _, err := p.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			log.Printf("[ERROR] cli error: %v", err)
		}
		os.Exit(2)
	}
	setupLog(opts.Dbg)
	log.Printf("[DEBUG] options: %+v", opts)
	err := run()
	if err != nil {
		log.Fatalf("[ERROR] detect bot failed, %v", err)
	}

}

func run() error {
	var err error
	_, cancel := context.WithCancel(context.Background())

	go func() {
		if x := recover(); x != nil {
			log.Printf("[WARN] run time panic:\n%v", x)
			panic(x)
		}

		// catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Printf("[WARN] interrupt signal")
		cancel()
	}()

	//var conf *config.FileConf
	//var err error
	//
	//if opts.Config != "" {
	//      conf, err = config.ParseConf(opts.Config, opts.SaveDir)
	//      if err != nil {
	//              return err
	//      }
	//}

	return err
}

func setupLog(dbg bool) {
	if dbg {
		log.Setup(log.Debug, log.CallerFile, log.CallerFunc, log.Msec, log.LevelBraces)
		return
	}
	log.Setup(log.Msec, log.LevelBraces)
}
