package config


type LbxConfig struct {
	Servers []ServerConfig `json:"servers"`
	Port    int            `json:"port"`
    RetryTimeInMinutes int  `json:"retryTimeInMinutes"`
}
