package configs

type Configs struct {
	PostgreSQL PostgreSQL
}

type PostgreSQL struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}