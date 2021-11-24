package mq

import mqtt "github.com/eclipse/paho.mqtt.golang"

type Options struct {
	ClientId string
	Username string
	Password string
	Brokers  []string
	Debug    bool

	// mqtt clean session
	// https://mcxiaoke.gitbooks.io/mqtt-cn/content/mqtt/0301-CONNECT.html
	CleanSession bool
	ResumeSubs   bool

	// mqtt qos
	// The Quality of Service 0,1,2 (default 0)
	// https://mcxiaoke.gitbooks.io/mqtt-cn/content/mqtt/0303-PUBLISH.html
	Qos int
	// Persistence
	Retained  bool
	OnConnect mqtt.OnConnectHandler

	// publish timeout
	PublishTimeout int
}

func ClientId(id string) Option {
	return func(o *Options) {
		o.ClientId = id
	}
}

func Auth(username, password string) Option {
	return func(o *Options) {
		o.Username = username
		o.Password = password
	}
}

func Brokers(brokers ...string) Option {
	return func(o *Options) {
		o.Brokers = brokers
	}
}

func Debug(debug bool) Option {
	return func(o *Options) {
		o.Debug = debug
	}
}

func CleanSession(clean bool) Option {
	return func(o *Options) {
		o.CleanSession = clean
	}
}

func ResumeSubs(resume bool) Option {
	return func(o *Options) {
		o.ResumeSubs = resume
	}
}

func Qos(qos int) Option {
	return func(o *Options) {
		o.Qos = qos
	}
}

func Retained(retained bool) Option {
	return func(o *Options) {
		o.Retained = retained
	}
}

func OnConnect(cb mqtt.OnConnectHandler) Option {
	return func(o *Options) {
		o.OnConnect = cb
	}
}

func PublishTimeout(timeout int) Option {
	return func(o *Options) {
		o.PublishTimeout = timeout
	}
}
