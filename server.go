package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"encoding/json"
	"time"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

var version = "untagged"

type Config struct {
	Database struct {
		DSN string
	}
	Server struct {
		Addr         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
	}
	Auth struct {
		IntrospectURL string
	}
}

func load(file string) (cfg Config, err error) {
	cfg.Server.Addr = ":8080"
	cfg.Server.ReadTimeout = 5 * time.Second
	cfg.Server.WriteTimeout = 5 * time.Second
	cfg.Server.IdleTimeout = 5 * time.Second

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return
	}

	return
}

func main() {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
	   sig := <-gracefulStop
	   log.Printf("caught sig: %+v", sig)
           os.Exit(0)
        }()
	
	log.Printf("Starting Server")
	
	filename := flag.String("config", "config.yml", "Configuration file")
	flag.Parse()

	log.Printf("Config file is %v", filename)
	config, err := load(*filename)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	log.Printf("Config is %v", config)
	handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)

		json.NewEncoder(rw).Encode(map[string]interface{}{
			"version": version,
		})
	})

	srv := &http.Server{
		ReadHeaderTimeout: config.Server.ReadTimeout,
		IdleTimeout:       config.Server.IdleTimeout,
		ReadTimeout:       config.Server.ReadTimeout,
		WriteTimeout:      config.Server.WriteTimeout,
		Addr:              config.Server.Addr,
		Handler:           handler,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
