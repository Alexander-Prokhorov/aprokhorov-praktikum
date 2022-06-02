package config

type Config struct {
	Server         string   `yaml:"SERVER"`
	Port           string   `yaml:"PORT"`
	PollInterval   string   `yaml:"POOL_INTERVAL"`
	ReportInterval string   `yaml:"REPORT_INTERVAL"`
	MemStatMetrics []string `yaml:"MEMSTAT_METRICS"`
}

func (c *Config) InitDefaults() {
	c.Server = "127.0.0.1"
	c.Port = "8080"
	c.PollInterval = "2s"
	c.ReportInterval = "10s"
	c.MemStatMetrics = []string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
	}
}
