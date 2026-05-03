package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go-vless-client/config"
)

func main() {
	vlessURI := flag.String("uri", "", "VLESS connection URI")
	socksPort := flag.Int("port", 1080, "Local SOCKS5 proxy port")
	flag.Parse()

	if *vlessURI == "" {
		fmt.Println("Usage: go-vless-client -uri \"vless://...\" [-port 1080]")
		os.Exit(1)
	}

	params, err := config.ParseVlessURI(*vlessURI)
	if err != nil {
		fmt.Printf("Failed to parse VLESS URI: %v\n", err)
		os.Exit(1)
	}

	instance, err := config.StartXray(params, *socksPort)
	if err != nil {
		fmt.Printf("Failed to start xray: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("SOCKS5 proxy listening on 127.0.0.1:%d\n", *socksPort)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	err = instance.Close()
	if err != nil {
		fmt.Printf("Error closing xray: %v\n", err)
	}
}
