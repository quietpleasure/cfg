package cfg

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
)

type Config struct {
	Data map[string]interface{}
}

func New() *Config {
	return &Config{
		Data: make(map[string]interface{}),
	}
}

func (c *Config) LoadGlob(pattern string) error {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	for _, filename := range files {
		err = c.LoadFile(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO: add .ini
func (c *Config) LoadFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	switch filepath.Ext(filename) {
	case ".yml", ".yaml":
		return c.LoadYAML(data)
	case ".json":
		return c.LoadJSON(data)
	case ".toml":
		return c.LoadTOML(data)
	}
	return nil
}

// Default tag=cfg
func (c *Config) Decode(out interface{}, tagname ...string) error {
	tag := "cfg"
	if tagname != nil {
		tag = tagname[0]
	}
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          tag,
		WeaklyTypedInput: true,
		Result:           out,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
		),
		ErrorUnused: true,
		ErrorUnset:  true,
	})
	if err != nil {
		return err
	}
	return d.Decode(c.Data)
}
