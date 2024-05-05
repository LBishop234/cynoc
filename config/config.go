package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"main/domain"
	"main/log"

	"gopkg.in/yaml.v3"
)

func ReadConfig(fPath string) (domain.SimConfig, error) {
	var config domain.SimConfig
	var err error

	log.Log.Debug().Msg("reading config file")

	switch filepath.Ext(fPath) {
	case ".yaml", ".yml":
		config, err = readYaml(fPath)
	case ".json":
		config, err = readJson(fPath)
	default:
		log.Log.Error().Err(domain.ErrInvalidFilepath).Str("ext", filepath.Ext(fPath)).Msg("invalid config file extension")
		return domain.SimConfig{}, domain.ErrInvalidFilepath
	}

	if err != nil {
		log.Log.Error().Err(err).Str("path", fPath).Msg("error reading config file")
		return domain.SimConfig{}, err
	}

	log.Log.Info().Msg("loaded config from file")
	return config, validate(config)
}

func validate(conf domain.SimConfig) error {
	if conf.CycleLimit < 1 {
		log.Log.Error().Int("cycle_limit", conf.CycleLimit).Msg("invalid cycle limit")
		return domain.ErrInvalidConfig
	}

	if conf.MaxPriority < 1 {
		log.Log.Error().Int("max_priority", conf.MaxPriority).Msg("invalid max priority")
		return domain.ErrInvalidConfig
	}

	if conf.BufferSize < 1 {
		log.Log.Error().Int("buffer_size", conf.BufferSize).Msg("invalid buffer size")
		return domain.ErrInvalidConfig
	}

	if conf.BufferSize%conf.MaxPriority != 0 {
		log.Log.Error().Int("buffer_size", conf.BufferSize).Int("max_priority", conf.MaxPriority).Msg("max priority must be a factor of buffer size")
		return domain.ErrInvalidConfig
	}

	if conf.FlitSize < 1 {
		log.Log.Error().Int("flit_size", conf.FlitSize).Msg("invalid flit size")
		return domain.ErrInvalidConfig
	}

	if conf.BufferSize%conf.FlitSize != 0 {
		log.Log.Error().Int("buffer_size", conf.BufferSize).Int("flit_size", conf.FlitSize).Msg("flit size must be a factor of buffer size")
		return domain.ErrInvalidConfig
	}

	if conf.ProcessingDelay < 1 {
		log.Log.Error().Int("processing_delay", conf.ProcessingDelay).Msg("invalid processing delay")
		return domain.ErrInvalidConfig
	}

	return nil
}

func readYaml(fPath string) (domain.SimConfig, error) {
	bytes, err := os.ReadFile(fPath)
	if err != nil {
		log.Log.Error().Err(err).Str("path", fPath).Msg("error reading .yaml config file")
		return domain.SimConfig{}, err
	}
	log.Log.Debug().Msg("read .yaml config file")

	var config domain.SimConfig
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		log.Log.Error().Err(err).Str("path", fPath).Msg("error unmarshalling .yaml config file")
		return domain.SimConfig{}, err
	}
	log.Log.Debug().Msg("unmarshalled .yaml config file")

	return config, nil
}

func readJson(fPath string) (domain.SimConfig, error) {
	bytes, err := os.ReadFile(fPath)
	if err != nil {
		log.Log.Error().Err(err).Str("path", fPath).Msg("error reading .json config file")
		return domain.SimConfig{}, err
	}
	log.Log.Debug().Msg("read .json config file")

	var config domain.SimConfig
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Log.Error().Err(err).Str("path", fPath).Msg("error unmarshalling .yaml config file")
		return domain.SimConfig{}, err
	}
	log.Log.Debug().Msg("unmarshalled .json config file")

	return config, nil
}
