package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/weiiwang01/wpex/internal/config"
	"github.com/weiiwang01/wpex/internal/relay"
	"golang.org/x/time/rate"
)

var version string

type pubKeys []string

func (ks *pubKeys) String() string {
	return strings.Join(*ks, ",")
}

func (ks *pubKeys) Set(s string) error {
	*ks = append(*ks, s)
	return nil
}

func main() {
	configFile := flag.String("config", "config.yaml", "config file path")

	conf, err := config.ReadConfig(*configFile)

	if err != nil {
		log.Fatal("unable to read config file", "error", err)
	}

	// bind := flag.String("bind", "", "address to bind to")
	// port := flag.Uint("port", 40000, "port number to listen")
	port := conf.Port
	bind := conf.Bind
	debug := flag.Bool("debug", false, "enable debug messages")
	// broadcastRate := flag.Uint("broadcast-rate", 0, "broadcast rate limit in packet per second")
	broadcastRate := conf.BroadcastRate
	versionFlag := flag.Bool("version", false, "show version number and quit")
	var allows pubKeys
	allows = conf.Allows
	// flag.Var(&allows, "allow", "allow a wireguard public key. --allow can be used multiple times for allowing multiple public keys")
	flag.Parse()

	if *versionFlag {
		fmt.Println("wpex", version)
		os.Exit(0)
	}

	loggingLevel := new(slog.LevelVar)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: loggingLevel}))
	if *debug {
		loggingLevel.Set(slog.LevelDebug)
	}
	slog.SetDefault(logger)
	address := fmt.Sprintf("%s:%d", bind, port)
	var allowKeys [][]byte
	for _, allow := range allows {
		k, err := base64.StdEncoding.DecodeString(allow)
		if err != nil || len(k) != 32 {
			log.Fatalf("invalid wireguard public key: '%s'", allow)
		}
		logger.Debug("allow wireguard public key", "key", allow)
		allowKeys = append(allowKeys, k)
	}
	limit := rate.Limit(broadcastRate)
	if broadcastRate == 0 {
		slog.Debug("broadcast rate limit is set to +Inf")
		limit = rate.Inf
	} else {
		slog.Debug(fmt.Sprintf("broadcast rate limit is set to %d", broadcastRate))
	}
	relay.Start(address, allowKeys, rate.NewLimiter(limit, int((broadcastRate)*5)))
}
