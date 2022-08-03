package main

import (
	"cam_cp/app/sender"
	"cam_cp/app/watcher"
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
	In  struct {
		Ftp struct {
			Enabled       bool          `long:"enabled" env:"ENABLED"  description:"enable ftp watcher"`
			Ip            string        `long:"ip" env:"IP" default:"127.0.0.1" description:"ip address of ftp server"`
			User          string        `long:"user" env:"USER" default:"anonymous" description:"user name"`
			Password      string        `long:"password" env:"PASSWORD" default:"" description:"user password"`
			Dir           string        `long:"dir" env:"DIR" default:"/" description:"ftp directory for recursive read"`
			CheckInterval time.Duration `long:"interval" env:"INTERVAL" default:"30s" description:"ftp check interval"`
		} `group:"ftp" namespace:"ftp" env-namespace:"FTP"`
		File struct {
			Dir           string        `long:"dir" env:"DIR" default:"/tmp" description:"file directory for recursive read"`
			CheckInterval time.Duration `long:"interval" env:"INTERVAL" default:"30s" description:"file check interval"`
		} `group:"file" namespace:"file" env-namespace:"FILE"`
		Http struct {
			Url           string        `long:"url" env:"URL" default:"http://localhost:8080" description:"http url"`
			CheckInterval time.Duration `long:"interval" env:"INTERVAL" default:"30s" description:"http check interval"`
		} `group:"http" namespace:"http" env-namespace:"HTTP"`
	} `group:"in" namespace:"in" env-namespace:"IN"`
	Filter struct {
		Deepstack struct {
			Url    string `long:"url" env:"URL" default:"http://localhost:8080" description:"deepstack url"`
			ApiKey string `long:"api-key" env:"API_KEY" default:"" description:"deepstack api key"`
		} `group:"deepstack" namespace:"deepstack" env-namespace:"DEEPSTACK"`
	} `group:"filter" namespace:"filter" env-namespace:"FILTER"`
	Out struct {
		File struct {
			Enabled bool   `long:"enabled" env:"ENABLED"  description:"enable file sender"`
			Dir     string `long:"dir" env:"DIR" default:"/tmp" description:"file directory for saving"`
		} `group:"file" namespace:"file" env-namespace:"FILE"`
		Http struct {
			Url string `long:"url" env:"URL" default:"http://localhost:8080" description:"http url"`
		} `group:"http" namespace:"http" env-namespace:"HTTP"`
		Ftp struct {
			Ip       string `long:"ip" env:"IP" default:"127.0.0.1" description:"ip address of ftp server"`
			User     string `long:"user" env:"USER" default:"anonymous" description:"user name"`
			Password string `long:"password" env:"PASSWORD" default:"" description:"user password"`
			Dir      string `long:"dir" env:"DIR" default:"/" description:"ftp directory for recursive write"`
		} `group:"ftp" namespace:"ftp" env-namespace:"FTP"`
		Telegram struct {
			Token  string `long:"token" env:"TOKEN" default:"" description:"telegram token"`
			ChatId string `long:"chat-id" env:"CHAT_ID" default:"" description:"telegram chat id"`
		} `group:"telegram" namespace:"telegram" env-namespace:"TELEGRAM"`
	} `group:"out" namespace:"out" env-namespace:"OUT"`
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
	ctx, cancel := context.WithCancel(context.Background())

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
	out := make(watcher.ExChan)
	if opts.In.Ftp.Enabled {
		runFtpWatcher(ctx, out)
	}
	if opts.Out.File.Enabled {
		runFileSender(ctx, out)
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("[WARN] detect bot stopped")
			return nil
			//case exchanges := <-out:
			//	for _, exchange := range exchanges {
			//		log.Printf("[INFO] detected %s", exchange.Name)
			//	}
		}
	}
	return err
}

func runFtpWatcher(ctx context.Context, out watcher.ExChan) {
	f := watcher.NewFtp(opts.In.Ftp.Ip, opts.In.Ftp.Dir,
		opts.In.Ftp.User, opts.In.Ftp.Password,
		opts.In.Ftp.CheckInterval)

	go func() {
		err := f.Run(ctx, out)
		if err != nil {
			log.Printf("[ERROR] Run ftp watcher failed, %v", err)
		}
	}()

}

func runFileSender(ctx context.Context, out watcher.ExChan) {
	f := sender.NewFile(opts.Out.File.Dir)
	go func() {
		err := f.Run(ctx, out)
		if err != nil {
			log.Printf("[ERROR] Run file sender failed, %v", err)
		}
	}()
}

func setupLog(dbg bool) {
	if dbg {
		log.Setup(log.Debug, log.CallerFile, log.CallerFunc, log.Msec, log.LevelBraces)
		return
	}
	log.Setup(log.Msec, log.LevelBraces)
}
