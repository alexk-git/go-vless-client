package config

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/xtls/xray-core/core"
	_ "github.com/xtls/xray-core/app/dispatcher"
	_ "github.com/xtls/xray-core/app/proxyman/inbound"
	_ "github.com/xtls/xray-core/app/proxyman/outbound"
	_ "github.com/xtls/xray-core/proxy/socks"
	_ "github.com/xtls/xray-core/proxy/vless/outbound"
	_ "github.com/xtls/xray-core/transport/internet/reality"
	_ "github.com/xtls/xray-core/transport/internet/tcp"
	_ "github.com/xtls/xray-core/main/json"
)

type VlessParams struct {
	UUID       string
	Address    string
	Port       int
	Encryption string
	Flow       string
	Security   string
	SNI        string
	Fingerprint string
	PublicKey  string
	ShortID    string
	Network    string
	Remark     string
}

func ParseVlessURI(uri string) (*VlessParams, error) {
	uri = strings.ReplaceAll(uri, "\\u0026", "&")

	if !strings.HasPrefix(uri, "vless://") {
		return nil, fmt.Errorf("not a vless URI")
	}

	// vless://UUID@ADDRESS:PORT?params#remark
	// Replace vless:// with https:// so net/url can parse it.
	parsed, err := url.Parse("https://" + uri[len("vless://"):])
	if err != nil {
		return nil, fmt.Errorf("failed to parse URI: %w", err)
	}

	uuid := parsed.User.Username()
	host := parsed.Hostname()
	portStr := parsed.Port()

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	q := parsed.Query()

	remark := ""
	if parsed.Fragment != "" {
		remark = parsed.Fragment
	}

	return &VlessParams{
		UUID:        uuid,
		Address:     host,
		Port:        port,
		Encryption:  q.Get("encryption"),
		Flow:        q.Get("flow"),
		Security:    q.Get("security"),
		SNI:         q.Get("sni"),
		Fingerprint: q.Get("fp"),
		PublicKey:   q.Get("pbk"),
		ShortID:     q.Get("sid"),
		Network:     q.Get("type"),
		Remark:      remark,
	}, nil
}

func buildJSON(params *VlessParams, socksPort int) string {
	return fmt.Sprintf(`{
  "inbounds": [
    {
      "tag": "socks-in",
      "port": %d,
      "listen": "0.0.0.0",
      "protocol": "socks",
      "settings": {
        "udp": true
      }
    }
  ],
  "outbounds": [
    {
      "tag": "vless-out",
      "protocol": "vless",
      "settings": {
        "vnext": [
          {
            "address": "%s",
            "port": %d,
            "users": [
              {
                "id": "%s",
                "encryption": "%s",
                "flow": "%s"
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "%s",
        "security": "%s",
        "realitySettings": {
          "serverName": "%s",
          "fingerprint": "%s",
          "publicKey": "%s",
          "shortId": "%s"
        }
      }
    }
  ]
}`,
		socksPort,
		params.Address,
		params.Port,
		params.UUID,
		params.Encryption,
		params.Flow,
		params.Network,
		params.Security,
		params.SNI,
		params.Fingerprint,
		params.PublicKey,
		params.ShortID,
	)
}

func StartXray(params *VlessParams, socksPort int) (*core.Instance, error) {
	jsonConfig := buildJSON(params, socksPort)

	config, err := core.LoadConfig("json", strings.NewReader(jsonConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	instance, err := core.New(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create xray instance: %w", err)
	}

	if err := instance.Start(); err != nil {
		return nil, fmt.Errorf("failed to start xray: %w", err)
	}

	return instance, nil
}
