package databases

type Database struct {
	ConnectionName  string `toml:"connectionName"`
	Engine          string `toml:"engine"`
	ConnectionMode  string `toml:"connectionMode"`
	Host            string `toml:"host"`
	Port            string `toml:"port"`
	Username        string `toml:"username"`
	Password        string `toml:"password"`
	DefaultDatabase string `toml:"defaultDatabase"`
}
