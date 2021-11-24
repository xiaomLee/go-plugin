package mosquitto

import (
	"crypto/tls"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"github.com/xiaomLee/go-plugin/mq"
)

type Client struct {
	mqtt.Client
	options mq.Options
}

func (c *Client) Init(opts ...mq.Option) error {
	c.configure(opts...)
	connOpts := mqtt.NewClientOptions()

	if c.options.ClientId != "" {
		connOpts.SetClientID(c.options.ClientId)
	}
	if c.options.Username != "" {
		connOpts.SetUsername(c.options.Username)
	}
	if c.options.Password != "" {
		connOpts.SetPassword(c.options.Password)
	}
	if c.options.OnConnect != nil {
		connOpts.SetOnConnectHandler(c.options.OnConnect)
	} else {
		connOpts.SetOnConnectHandler(func(c mqtt.Client) {
			logrus.Infoln("mqtt OnConnect ")
		})
	}
	connOpts.OnConnectionLost = func(client mqtt.Client, err error) {
		logrus.Infoln("mqtt connect lost ", err)
	}
	connOpts.SetCleanSession(c.options.CleanSession)
	connOpts.SetResumeSubs(c.options.ResumeSubs)

	for _, broker := range c.options.Brokers {
		connOpts.AddBroker(broker)
	}

	connOpts.AutoReconnect = true
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	c.Client = mqtt.NewClient(connOpts)
	if token := c.Client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (c *Client) configure(opts ...mq.Option) {
	for _, o := range opts {
		o(&c.options)
	}
}
