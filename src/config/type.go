package config

type TwitchRedeem struct {
	Name   string
	Status string
}

type TwitchChannel struct {
	Id   string
	Name string
}

type TwitchAuth struct {
	ClientID     string
	ClientSecret string
	OAuth        string
	RefreshToken string
}

type Twitch struct {
	Channel TwitchChannel
	Auth    TwitchAuth
	Redeem  TwitchRedeem
}

type Output struct {
	Type string
	Tapo Tapo
}

type Tapo struct {
	IP       string
	Username string
	Password string
}

type GPIO struct {
	Enabled   bool
	RedeemPin int
	TapoPin   int
}

type Logging struct {
	LogFile   string
	UseSyslog bool
}

type Web struct {
	Enabled bool
	Port    string
}

type Config struct {
	Twitch  Twitch
	Output  Output
	GPIO    GPIO
	Logging Logging
	Web     Web
}
