package main

import (
	"cam_cp/app/filter"
	"cam_cp/app/frame"
	"cam_cp/app/sender"
	"cam_cp/app/watcher"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/go-pkgz/lgr"
	"github.com/umputun/go-flags"
)

var opts struct {
	Dbg bool `long:"dbg" env:"DEBUG" description:"debug mode"`
	In  struct {
		Ftp struct {
			Enable        bool          `long:"enable" env:"ENABLE" description:"enable ftp"`
			Ip            string        `long:"ip" env:"IP" default:"127.0.0.1" description:"ip address of ftp server"`
			User          string        `long:"user" env:"USER" default:"anonymous" description:"user name"`
			Password      string        `long:"password" env:"PASSWORD" default:"" description:"user password"`
			Dir           string        `long:"dir" env:"DIR" default:"/" description:"ftp directory for recursive read"`
			CheckInterval time.Duration `long:"interval" env:"INTERVAL" default:"30s" description:"ftp check interval"`
		} `group:"ftp" namespace:"ftp" env-namespace:"FTP"`
		File struct {
			Enable        bool          `long:"enable" env:"ENABLE" description:"enable file"`
			Dir           string        `long:"dir" env:"DIR" default:"/tmp" description:"file directory for recursive read"`
			CheckInterval time.Duration `long:"interval" env:"INTERVAL" default:"30s" description:"file check interval"`
		} `group:"file" namespace:"file" env-namespace:"FILE"`
		HttpJpeg struct {
			Enable        bool   `long:"enable" env:"ENABLE" description:"enable http"`
			Url           string `long:"url" env:"URL" default:"http://localhost:8080" description:"http url"`
			CheckInterval int64  `long:"interval" env:"INTERVAL" default:"500" description:"http check interval in milliseconds"`
			Framebuffer   int    `long:"framebuffer" env:"FRAMEBUFFER" default:"15" description:"Frame buffer size"`
		} `group:"Jpeg" namespace:"Jpeg" env-namespace:"JPEG"`
		HttpMJpeg struct {
			Enable      bool   `long:"enable" env:"ENABLE" description:"enable http"`
			Url         string `long:"url" env:"URL" default:"http://localhost:8080" description:"http url"`
			Framebuffer int    `long:"framebuffer" env:"FRAMEBUFFER" default:"15" description:"Frame buffer size"`
		} `group:"MJpeg" namespace:"MJpeg" env-namespace:"MJPEG"`
	} `group:"in" namespace:"in" env-namespace:"IN"`
	Filter struct {
		Deepstack struct {
			Enable     bool    `long:"enable" env:"ENABLE" description:"enable deepstack"`
			Url        string  `long:"url" env:"URL" default:"http://localhost:8080" description:"deepstack url"`
			ApiKey     string  `long:"api-key" env:"API_KEY" default:"" description:"deepstack api key"`
			Confidence float64 `long:"confidence" env:"CONFIDENCE" default:"0.5" description:"confidence level"`
			Labels     string  `long:"labels" env:"LABELS" default:"person" description:"comma separated labels to detect"`
		} `group:"deepstack" namespace:"deepstack" env-namespace:"DEEPSTACK"`
		Yolo struct {
			Enable      bool    `long:"enable" env:"ENABLE" description:"enable yolo"`
			Config      string  `long:"config" env:"CONFIG" default:"./yolov3.cfg" description:"yolo config file"`
			Weights     string  `long:"weights" env:"WEIGHTS" default:"./yolov3.weights" description:"yolo weights file"`
			Threshold   float32 `long:"threshold" env:"THRESHOLD" default:"0.25" description:"threshold level"`
			Probability float32 `long:"probability" env:"PROBABILITY" default:"75.0" description:"probability in %"`
			GPUIndex    int     `long:"gpu-index" env:"GPU_INDEX" default:"0" description:"gpu device index"`
			Labels      string  `long:"labels" env:"LABELS" default:"person" description:"comma separated labels to detect"`
		} `group:"yolo" namespace:"yolo" env-namespace:"YOLO"`
	} `group:"filter" namespace:"filter" env-namespace:"FILTER"`
	Out struct {
		File struct {
			Enable bool   `long:"enable" env:"ENABLE" description:"enable file"`
			Dir    string `long:"dir" env:"DIR" default:"/tmp" description:"file directory for saving"`
		} `group:"file" namespace:"file" env-namespace:"FILE"`
		Http struct {
			Url string `long:"url" env:"URL" default:"http://localhost:8080" description:"http url"`
		} `group:"http" namespace:"http" env-namespace:"HTTP"`
		Ftp struct {
			Enable   bool   `long:"enable" env:"ENABLE" description:"enable ftp"`
			Ip       string `long:"ip" env:"IP" default:"127.0.0.1" description:"ip address of ftp server"`
			User     string `long:"user" env:"USER" default:"anonymous" description:"user name"`
			Password string `long:"password" env:"PASSWORD" default:"" description:"user password"`
			Dir      string `long:"dir" env:"DIR" default:"/" description:"ftp directory for recursive write"`
		} `group:"ftp" namespace:"ftp" env-namespace:"FTP"`
		Telegram struct {
			Enable bool   `long:"enable" env:"ENABLE" description:"enable telegram"`
			Token  string `long:"token" env:"TOKEN" default:"" description:"telegram token"`
			ChatId int64  `long:"chat-id" env:"CHAT_ID" default:"" description:"telegram chat id"`
		} `group:"telegram" namespace:"telegram" env-namespace:"TELEGRAM"`
		Email struct {
			Enable   bool   `long:"enable" env:"ENABLE" description:"enable email"`
			Host     string `long:"host" env:"HOST" default:"localhost" description:"smtp host"`
			Port     int    `long:"port" env:"PORT" default:"25" description:"smtp port"`
			Password string `long:"password" env:"PASSWORD" default:"" description:"smtp password"`
			To       string `long:"to" env:"TO" default:"" description:"email to"`
			From     string `long:"from" env:"FROM" default:"" description:"email from"`
			Subject  string `long:"subject" env:"SUBJECT" default:"" description:"email subject"`
			Body     string `long:"body" env:"BODY" default:"" description:"email body"`
		} `group:"email" namespace:"email" env-namespace:"EMAIL"`
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

// Watcher->FrameOutChannel->Read in main->Filter->Sender
func run() error {
	var (
		err      error
		watchers []watcher.Watcher
		filters  []filter.Filter
		senders  []sender.Sender
	)

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

	if opts.In.Ftp.Enable {
		var ftpWatcher watcher.Watcher
		ftpWatcher, err = watcher.NewFtp(opts.In.Ftp.Ip, opts.In.Ftp.Dir,
			opts.In.Ftp.User, opts.In.Ftp.Password,
			opts.In.Ftp.CheckInterval)
		if err != nil {
			return err
		}
		watchers = append(watchers, ftpWatcher)
	}

	if opts.In.HttpMJpeg.Enable {
		var httpMJpegWatcher watcher.Watcher
		httpMJpegWatcher, err = watcher.NewHttpMJpeg(opts.In.HttpMJpeg.Url, opts.In.HttpMJpeg.Framebuffer)
		watchers = append(watchers, httpMJpegWatcher)
	}

	if opts.In.HttpJpeg.Enable {
		var httpJpegWatcher watcher.Watcher
		httpJpegWatcher, err = watcher.NewHttpJpeg(opts.In.HttpJpeg.Url, opts.In.HttpJpeg.Framebuffer,
			opts.In.HttpJpeg.CheckInterval)
		watchers = append(watchers, httpJpegWatcher)
	}

	if opts.Out.File.Enable {
		var fileSender sender.Sender
		fileSender, err = sender.NewFile(opts.Out.File.Dir)
		if err != nil {
			return err
		}
		senders = append(senders, fileSender)
	}

	if opts.Out.Telegram.Enable {
		var telegramSender sender.Sender
		telegramSender, err = sender.NewTelegram(opts.Out.Telegram.Token, opts.Out.Telegram.ChatId)
		if err != nil {
			return err
		}
		senders = append(senders, telegramSender)
	}

	if opts.Out.Email.Enable {
		var emailSender sender.Sender
		emailSender, err = sender.NewEmail(opts.Out.Email.Host, opts.Out.Email.Port, opts.Out.Email.Password,
			opts.Out.Email.To, opts.Out.Email.From, opts.Out.Email.Subject, opts.Out.Email.Body)
		if err != nil {
			return err
		}
		senders = append(senders, emailSender)
	}

	if opts.Filter.Deepstack.Enable {
		var deepstackFilter filter.Filter
		deepstackFilter, err = filter.NewDeepstack(opts.Filter.Deepstack.Url,
			opts.Filter.Deepstack.ApiKey, opts.Filter.Deepstack.Labels,
			opts.Filter.Deepstack.Confidence)
		if err != nil {
			return err
		}
		filters = append(filters, deepstackFilter)
	}
	if opts.Filter.Yolo.Enable {
		var yoloFilter filter.Filter
		yoloFilter, err = filter.NewYolo(opts.Filter.Yolo.GPUIndex,
			opts.Filter.Yolo.Config,
			opts.Filter.Yolo.Weights,
			opts.Filter.Yolo.Labels,
			opts.Filter.Yolo.Threshold,
			opts.Filter.Yolo.Probability)

		if err != nil {
			return err
		}
		filters = append(filters, yoloFilter)
		defer yoloFilter.Close()
	}

	framesChan := make(chan []frame.Frame)

	if len(watchers) == 0 {
		return errors.New("no watcher(s) defined")
	}

	for _, wtc := range watchers {
		//Run watcher
		go func(w watcher.Watcher) {
			err = w.Watch(ctx, framesChan)
			if err != nil {
				log.Printf("[ERROR]  watcher failed, %v", err)
				return
			}
		}(wtc)

	}
	if err != nil {
		return err
	}

	// Run main loop
	for {
		select {
		case inFrames := <-framesChan:
			var outFrames []frame.Frame
			for _, fltr := range filters {
				outFrames = append(outFrames, fltr.Filter(inFrames)...)
			}
			if outFrames != nil {
				for _, sndr := range senders {
					err = sndr.Send(outFrames)
					if err != nil {
						log.Printf("[ERROR] sender failed, %v", err)
						continue
					}
				}
			}
		case <-ctx.Done():
			log.Printf("[WARN] detect bot stopped")
			return nil
		}
	}

}

func setupLog(dbg bool) {
	if dbg {
		log.Setup(log.Debug, log.CallerFile, log.CallerFunc, log.Msec, log.LevelBraces)
		return
	}
	log.Setup(log.Msec, log.LevelBraces)
}
