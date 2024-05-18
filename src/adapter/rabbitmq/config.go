package rabbitmq

type Config struct {
	User                        string `koanf:"user"`
	Password                    string `koanf:"password"`
	Host                        string `koanf:"host"`
	Port                        int    `koanf:"port"`
	Vhost                       string `koanf:"vhost"`
	ReconnectSecond             int    `koanf:"reconnect_second"`
	BufferSize                  int    `koanf:"buffer_size"`
	MaxRetryPolicy              int    `koanf:"max_retry_policy"`
	ChannelCleanerTimerInSecond int    `koanf:"channel_cleaner_timer_in_seconds"`
}
