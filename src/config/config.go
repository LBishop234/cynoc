package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"main/log"
	"main/src/domain"

	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidConfig          = errors.New("invalid config")
	ErrInvalidCycleLimit      = errors.New("invalid cycle limit")
	ErrInvalidMaxPriority     = errors.New("invalid max priority")
	ErrInvalidBufferSize      = errors.New("invalid buffer size")
	ErrInvalidProcessingDelay = errors.New("invalid processing delay")
	ErrInvalidLinkBandwidth   = errors.New("invalid link bandwidth")
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
		log.Log.Error().Err(ErrInvalidConfig).Str("ext", filepath.Ext(fPath)).Msg("invalid config file extension must be yaml or json")
		return domain.SimConfig{}, ErrInvalidConfig
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
		err := errors.Join(ErrInvalidConfig, ErrInvalidCycleLimit)
		log.Log.Error().Err(err).Int("cycle_limit", conf.CycleLimit).Msg("cycle limit must be greater than 0")
		return err
	}

	if conf.MaxPriority < 1 {
		err := errors.Join(ErrInvalidConfig, ErrInvalidMaxPriority)
		log.Log.Error().Err(err).Int("max_priority", conf.MaxPriority).Msg("max priority must be greater than 0")
		return err
	}

	if conf.BufferSize < 1 {
		err := errors.Join(ErrInvalidConfig, ErrInvalidBufferSize)
		log.Log.Error().Err(err).Int("buffer_size", conf.BufferSize).Msg("buffer size must be greater than 0")
		return err
	}

	if conf.BufferSize%conf.MaxPriority != 0 {
		err := errors.Join(ErrInvalidConfig, ErrInvalidBufferSize)
		log.Log.Error().Err(err).Int("buffer_size", conf.BufferSize).Int("max_priority", conf.MaxPriority).Msg("buffer size must be a multiple of max priority")
		return err
	}

	if conf.BufferSize < 1 {
		err := errors.Join(ErrInvalidConfig, ErrInvalidBufferSize)
		log.Log.Error().Err(err).Int("buffer_size", conf.BufferSize).Msg("buffer size must be greater than 0")
		return err
	}

	if conf.ProcessingDelay < 1 {
		err := errors.Join(ErrInvalidConfig, ErrInvalidProcessingDelay)
		log.Log.Error().Err(err).Int("processing_delay", conf.ProcessingDelay).Msg("processing delay must be greater than 0")
		return err
	}

	if conf.LinkBandwidth < 1 {
		err := errors.Join(ErrInvalidConfig, ErrInvalidLinkBandwidth)
		log.Log.Error().Err(err).Int("link_bandwidth", conf.LinkBandwidth).Msg("link bandwidth must be greater than 0")
		return err
	}

	if conf.LinkBandwidth > (conf.BufferSize / conf.MaxPriority) {
		err := errors.Join(ErrInvalidConfig, ErrInvalidLinkBandwidth)
		log.Log.Error().Err(err).Int("link_bandwidth", conf.LinkBandwidth).Int("buffer_size", conf.BufferSize).Int("max_priority", conf.MaxPriority).Msg("link bandwidth must be less than or equal to virtual channel size")
		return err
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
		log.Log.Error().Err(err).Str("path", fPath).Msg("error unmarshalling .json config file")
		return domain.SimConfig{}, err
	}
	log.Log.Debug().Msg("unmarshalled .json config file")

	return config, nil
}
