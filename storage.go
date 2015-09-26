package main

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

type channelConfig struct {
	Enabled bool
	Config  map[string]string
}

type pluginStorage struct {
	FilePath      string
	ChannelConfig map[string]*channelConfig
}

func (storage *pluginStorage) Start(channel int) error {
	storage.ChannelConfig[strconv.Itoa(channel)] = &channelConfig{
		Enabled: false,
		Config: map[string]string{
			"start": "Hi! I'm glad to be a part of your group. Look at /bothelp@GroupWelcomeBot for help.",
			"new":   "Hello, USER_NAME!",
			"left":  "Good bye, USER_NAME :(",
			"rules": "No rules have been set yet!",
		},
	}

	return storage.Save()
}

func (storage *pluginStorage) Enable(channel int) error {
	storage.ChannelConfig[strconv.Itoa(channel)].Enabled = true

	return storage.Save()
}

func (storage *pluginStorage) Disable(channel int) error {
	storage.ChannelConfig[strconv.Itoa(channel)].Enabled = false

	return storage.Save()
}

func (storage *pluginStorage) IsEnabled(channel int) bool {
	return storage.ChannelConfig[strconv.Itoa(channel)].Enabled
}

func (storage *pluginStorage) Get(channel int) map[string]string {
	return storage.GetByString(strconv.Itoa(channel))
}

func (storage *pluginStorage) GetByString(channel string) map[string]string {
	return storage.ChannelConfig[channel].Config
}

func (storage *pluginStorage) GetChannelConfig(channel string) *channelConfig {
	return storage.ChannelConfig[channel]
}

func (storage *pluginStorage) Set(channel int, key, value string) error {
	storage.ChannelConfig[strconv.Itoa(channel)].Config[key] = value

	return storage.Save()
}

func (storage *pluginStorage) Create() error {
	storage.ChannelConfig = make(map[string]*channelConfig)

	return storage.Save()
}

func (storage *pluginStorage) Load() error {
	f, err := ioutil.ReadFile(storage.FilePath)
	if err != nil {
		return err
	}

	var data map[string]*channelConfig
	err = json.Unmarshal(f, &data)
	if err != nil {
		return err
	}

	storage.ChannelConfig = data

	return nil
}

func (storage *pluginStorage) Save() error {
	j, err := json.MarshalIndent(storage.ChannelConfig, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(storage.FilePath, j, 0644)

	return err
}
