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
	Dbg   bool `long:"dbg" env:"DEBUG" description:"debug mode"`
	FtpIn struct {
		Enabled       bool          `long:"enabled" env:"ENABLED" description:"enable ftp watcher"`
		Ip            string        `long:"ip" env:"IP" default:"127.0.0.1" description:"ip address of ftp server"`
		User          string        `long:"user" env:"USER" default:"anonymous" description:"user name"`
		Password      string        `long:"password" env:"PASSWORD" default:"" description:"user password"`
		Dir           string        `long:"dir" env:"DIR" default:"/" description:"ftp directory for recursive read"`
		CheckInterval time.Duration `long:"interval" env:"INTERVAL" default:"3s" description:"ftp check interval"`
	} `group:"in" namespace:"ftp_in" env-namespace:"FTP_IN"`
	FileIn struct {
		Enabled       bool          `long:"enabled" env:"ENABLED" description:"enable file watcher"`
		Dir           string        `long:"dir" env:"DIR" default:"/tmp" description:"file directory for recursive read"`
		CheckInterval time.Duration `long:"interval" env:"INTERVAL" default:"3s" description:"file check interval"`
	} `group:"in" namespace:"file_in" env-namespace:"FILE_IN"`
	HttpIn struct {
		Enabled       bool          `long:"enabled" env:"ENABLED" description:"enable http watcher"`
		Url           string        `long:"url" env:"URL" default:"http://localhost:8080" description:"http url"`
		CheckInterval time.Duration `long:"interval" env:"INTERVAL" default:"3s" description:"http check interval"`
	} `group:"in" namespace:"http_in" env-namespace:"HTTP_IN"`

	Deepstack struct {
		Enabled bool   `long:"enabled" env:"ENABLED" description:"enable deepstack filter"`
		Url     string `long:"url" env:"URL" default:"http://localhost:8080" description:"deepstack url"`
		ApiKey  string `long:"api-key" env:"API_KEY" default:"" description:"deepstack api key"`
	} `group:"filters" namespace:"deepstack" env-namespace:"DEEPSTACK"`

	FileOut struct {
		Enabled bool   `long:"enabled" env:"ENABLED" description:"enable file saver"`
		Dir     string `long:"dir" env:"DIR" default:"/tmp" description:"file directory for saving"`
	} `group:"out" namespace:"file_out" env-namespace:"FILE_OUT"`
	HttpOut struct {
		Enabled bool   `long:"enabled" env:"ENABLED" description:"enable http saver"`
		Url     string `long:"url" env:"URL" default:"http://localhost:8080" description:"http url"`
	} `group:"out" namespace:"http_out" env-namespace:"HTTP_OUT"`
	FtpOut struct {
		Enabled  bool   `long:"enabled" env:"ENABLED" description:"enable ftp saver"`
		Ip       string `long:"ip" env:"IP" default:"127.0.0.1" description:"ip address of ftp server"`
		User     string `long:"user" env:"USER" default:"anonymous" description:"user name"`
		Password string `long:"password" env:"PASSWORD" default:"" description:"user password"`
		Dir      string `long:"dir" env:"DIR" default:"/" description:"ftp directory for recursive write"`
	} `group:"out" namespace:"ftp_out" env-namespace:"FTP_OUT"`
	TelegramOut struct {
		Enabled bool   `long:"enabled" env:"ENABLED" description:"enable telegram filter"`
		Token   string `long:"token" env:"TOKEN" default:"" description:"telegram token"`
		ChatId  string `long:"chat-id" env:"CHAT_ID" default:"" description:"telegram chat id"`
	} `group:"out" namespace:"telegram" env-namespace:"TELEGRAM"`
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

	return err
}

func setupLog(dbg bool) {
	if dbg {
		log.Setup(log.Debug, log.CallerFile, log.CallerFunc, log.Msec, log.LevelBraces)
		return
	}
	log.Setup(log.Msec, log.LevelBraces)
}
