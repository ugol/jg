package wamprpc

import (
	"context"
	"encoding/json"
	"os"

	"github.com/gammazero/nexus/v3/client"
	"github.com/gammazero/nexus/v3/wamp"
	"github.com/rs/zerolog/log"
)

type Config struct {
	WampURI   string `json:"wamp_uri"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Realm     string `json:"realm"`
	Procedure string `json:"procedure"`
	SerType   string `json:"serType"`
	Compress  bool   `json:"compress"`
	Authid    string `json:"authid"`
}

type Producer struct {
	client    client.Client
	realm     string
	procedure string
	authid    string
}

func (p *Producer) Initialize(ctx context.Context, configFile string) {
	var config Config
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read configuration file")
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse configuration parameters")
	}
	var wampclient *client.Client

	// Get requested serialization.
	serialization := client.JSON
	switch config.SerType {
	case "json":
	case "msgpack":
		serialization = client.MSGPACK
	case "cbor":
		serialization = client.CBOR
	default:
		log.Fatal().Err(err).Msg("Invalid serialization, muse be one of: json, msgpack, cbor")
	}

	cfg := client.Config{
		Realm:         config.Realm,
		Serialization: serialization,
		HelloDetails: wamp.Dict{
			"authid": config.Authid,
		},
	}

	if config.Compress {
		cfg.WsCfg.EnableCompression = true
	}

	addr := config.WampURI

	wampclient, err = client.ConnectNet(context.Background(), addr, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't connect to WAMP Router")
	}
	// defer wampclient.Close()

	p.realm = config.Realm
	p.procedure = config.Procedure
	p.authid = config.Authid

	p.client = *wampclient
}

func (p *Producer) Produce(ctx context.Context, k []byte, v []byte, _ any) {
	data := string(v)
	args := wamp.List{data}
	opts := wamp.Dict{
		"authid": p.authid,
	}
	_, err := p.client.Call(ctx, p.procedure, opts, args, nil, nil)
	if err != nil {
		log.Fatal().Err(err).Msgf("call error: %s", err)
	}
}

func (p *Producer) Close(ctx context.Context) error {
	err := p.client.Close()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to close WAMP connection")
	}
	return err
}
