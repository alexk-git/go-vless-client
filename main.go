package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"go-vless-client/config"
)

func loadEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, "\"'")
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

func main() {
	loadEnv(".env")

	envURI := os.Getenv("VLESS_URI")
	envPort := os.Getenv("SOCKS_PORT")

	defaultPort := 1080
	if envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			defaultPort = p
		}
	}

	vlessURI := flag.String("uri", envURI, "VLESS connection URI")
	socksPort := flag.Int("port", defaultPort, "Local SOCKS5 proxy port")
	flag.Parse()

	if *vlessURI == "" {
		fmt.Println("Usage: go-vless-client -uri \"vless://...\" [-port 1080]")
		fmt.Println("Or set VLESS_URI and SOCKS_PORT in .env file")
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
