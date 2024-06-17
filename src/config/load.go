package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"strings"
)

func Load(configPath string) Config {
	// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
	var k = koanf.New(".")

	// Load default values using the confmap provider.
	// We provide a flat map with the "." delimiter.
	// A nested map can be loaded by setting the delimiter to an empty string "".
	err := k.Load(confmap.Provider(defaultConfig, "."), nil)
	if err != nil {
		panic(err)
	}

	// Load YAML config and merge into the previously loaded config (because we can).
	err = k.Load(file.Provider(configPath), yaml.Parser())
	if err != nil {
		panic(err)
	}

	err = k.Load(env.Provider("MESSENGER_", ".", func(s string) string {
		str := strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "MESSENGER_")), "_", ".", -1)

		return strings.Replace(str, "..", "_", -1)
	}), nil)
	if err != nil {
		panic(err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		panic(err)
	}

	return cfg
}
