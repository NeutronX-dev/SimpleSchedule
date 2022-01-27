package ConfigLoader

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type Config struct {
	Path string
	// Preformance
	CheckInterval int `json:"check_interval"`
	// Misc
	Instructive bool   `json:"instructive"`
	AudioPath   string `json:"audio_path"`
}

func (cfg *Config) CheckIntervalDuration() time.Duration {
	return time.Duration(cfg.CheckInterval)
}

func (cfg *Config) Save() error {
	data, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(cfg.Path, data, 0777)

	if err != nil {
		return err
	}

	return nil
}

func ReadConfig(path string) (*Config, error) {
	var res *Config = &Config{Path: path, CheckInterval: 1000, AudioPath: "./assets/notification.mp3"}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		dat, err := ioutil.ReadFile(path)
		if err != nil {
			return res, err
		}

		err = json.Unmarshal(dat, &res)
		if err != nil {
			return res, err
		}
		return res, nil
	} else {
		err := CreateConfig(path, res)
		if err != nil {
			return res, err
		} else {
			return res, nil
		}
	}
}

func CreateConfig(path string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, data, 0777)

	if err != nil {
		return err
	}

	return nil
}
