package config

type MUTConfig struct {
	AppName            string `json:"app_name,omitempty"`
	HTTPPort           int    `json:"http_port,omitempty"`
	HTTPMaxRequestTime int    `json:"max_request_time,omitempty"`
}

type PgsqlConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Schema          string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
}

type MUTDatabase struct {
	MUTEngine          string `json:"engine,omitempty"`
	MUTHost            string `json:"host,omitempty"`
	MUTPort            int    `json:"port,omitempty"`
	MUTUsername        string `json:"username,omitempty"`
	MUTPassword        string `json:"-"`
	MUTDBName          string `json:"database_name,omitempty"`
	MUTSchema          string `json:"schema,omitempty"`
	MUTMaxIdle         int    `json:"max_idle,omitempty"`
	MUTMaxConn         int    `json:"max_conn,omitempty"`
	MUTConnMaxLifetime int    `json:"conn_max_lifetime,omitempty"`
}

type Redis struct {
	Host         string `json:"host,omitempty"`
	Port         int    `json:"port,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"-"`
	DB           int    `json:"db,omitempty"`
	UseTLS       bool   `json:"use_tls,omitempty"`
	MaxRetries   int    `json:"max_retries"`
	MinIdleConns int    `json:"min_idle_conns"`
	PoolSize     int    `json:"pool_size"`
	PoolTimeout  int    `json:"pool_timeout"`
	MaxConnAge   int    `json:"max_conn_age"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}
