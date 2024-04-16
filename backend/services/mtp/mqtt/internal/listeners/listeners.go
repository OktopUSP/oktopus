package listeners

import (
	"broker/internal/config"
	"broker/internal/listeners/http"
	broker "broker/internal/listeners/mqtt"
	"broker/internal/listeners/ws"
	"sync"

	"github.com/mochi-co/mqtt/v2"
	"github.com/rs/zerolog"
)

func StartServers(c config.Config) {

	server := mqtt.New(&mqtt.Options{})

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		mqttServer := newMqttServer(c)
		mqttServer.Start(server)
		wg.Done()
	}()

	go func() {
		if c.WsEnable {
			wsServer := newWsServer(c)
			wsServer.Start(server)
		}
		wg.Done()
	}()

	go func() {
		if c.HttpEnable {
			httpServer := newHttpServer(c)
			httpServer.Start(server)
		}
		wg.Done()
	}()

	server.Log.Level(zerolog.Level(c.LogLevel))

	wg.Wait()

	err := server.Serve()
	if err != nil {
		server.Log.Fatal().Err(err).Msg("server error")
	}
}

func newMqttServer(c config.Config) *broker.Mqtt {
	return &broker.Mqtt{
		Port:       c.MqttPort,
		TlsPort:    c.MqttTlsPort,
		NoTls:      c.NoTls,
		Tls:        c.Tls,
		Fullchain:  c.Fullchain,
		Privkey:    c.Privkey,
		AuthEnable: c.AuthEnable,
		Redis: broker.Redis{
			RedisEnable:   c.RedisEnable,
			RedisAddr:     c.RedisAddr,
			RedisPassword: c.RedisPassword,
		},
		LogLevel: c.LogLevel,
		Nats:     c.Nats,
	}
}

func newWsServer(c config.Config) *ws.Ws {
	return &ws.Ws{
		WsPort: c.WsPort,
	}
}

func newHttpServer(c config.Config) *http.Http {
	return &http.Http{
		HttpPort: c.HttpPort,
	}
}
