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
		log.Log.Error().Err(domain.ErrInvalidConfig).Int("cycle_limit", conf.CycleLimit).Msg("cycle limit must be greater than 0")
		return domain.ErrInvalidConfig
	}

	if conf.MaxPriority < 1 {
		log.Log.Error().Err(domain.ErrInvalidConfig).Int("max_priority", conf.MaxPriority).Msg("max priority must be greater than 0")
		return domain.ErrInvalidConfig
	}

	if conf.BufferSize < 1 {
		log.Log.Error().Err(domain.ErrInvalidConfig).Int("buffer_size", conf.BufferSize).Msg("buffer size must be greater than 0")
		return domain.ErrInvalidConfig
	}

	if conf.BufferSize%conf.MaxPriority != 0 {
		log.Log.Error().Err(domain.ErrInvalidConfig).Int("buffer_size", conf.BufferSize).Int("max_priority", conf.MaxPriority).Msg("buffer size must be a multiple of max priority")
		return domain.ErrInvalidConfig
	}

	if conf.FlitSize < 1 {
		log.Log.Error().Err(domain.ErrInvalidConfig).Int("flit_size", conf.FlitSize).Msg("flit size must be greater than 0")
		return domain.ErrInvalidConfig
	}

	if conf.BufferSize%conf.FlitSize != 0 {
		log.Log.Error().Err(domain.ErrInvalidConfig).Int("buffer_size", conf.BufferSize).Int("flit_size", conf.FlitSize).Msg("buffer size must be a multiple of flit size")
		return domain.ErrInvalidConfig
	}

	if conf.ProcessingDelay < 1 {
		log.Log.Error().Err(domain.ErrInvalidConfig).Int("processing_delay", conf.ProcessingDelay).Msg("processing delay must be greater than 0")
		return domain.ErrInvalidConfig
	}

	if conf.LinkBandwidth < 1 {
		log.Log.Error().Err(domain.ErrInvalidConfig).Int("link_bandwidth", conf.LinkBandwidth).Msg("link bandwidth must be greater than 0")
		return domain.ErrInvalidConfig
	}

	if conf.LinkBandwidth%conf.FlitSize != 0 {
		log.Log.Error().Err(domain.ErrInvalidConfig).Int("link_bandwidth", conf.LinkBandwidth).Int("flit_size", conf.FlitSize).Msg("link bandwidth must be a multiple of flit size")
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
