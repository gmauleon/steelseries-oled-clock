package clock

import (
	"time"

	"github.com/kardianos/service"
)

const (
	configFile = `C:\ProgramData\SteelSeries\SteelSeries Engine 3/coreProps.json`
)

// GameSenseClock represent a clock object
type GameSenseClock struct {
	DateFormat string
	TimeFormat string
	Address    string
	Ticker     *time.Ticker
	TickerDone chan bool
	Logger     service.Logger
}

// Binding represent a Steelseries GameSense binding request
type Binding struct {
	Game          string           `json:"game,omitempty"`
	Event         string           `json:"event,omitempty"`
	ValueOptional bool             `json:"value_optional,omitempty"`
	Handlers      []BindingHandler `json:"handlers,omitempty"`
}

// BindingHandler represent a Steelseries GameSense binding request
type BindingHandler struct {
	DeviceType string               `json:"device-type,omitempty"`
	Zone       string               `json:"zone,omitempty"`
	Mode       string               `json:"mode,omitempty"`
	Datas      []BindingHandlerData `json:"datas,omitempty"`
}

// BindingHandlerData represent a Steelseries GameSense binding request
type BindingHandlerData struct {
	IconID int                      `json:"icon-id,omitempty"`
	Lines  []BindingHandlerDataLine `json:"lines,omitempty"`
}

// BindingHandlerDataLine represent a Steelseries GameSense binding request
type BindingHandlerDataLine struct {
	HasText         bool   `json:"has-text,omitempty"`
	Bold            bool   `json:"bold,omitempty"`
	ContextFrameKey string `json:"context-frame-key,omitempty"`
}

// GameMetadata represent a Steelseries GameSense game metadata request
type GameMetadata struct {
	Game            string `json:"game,omitempty"`
	GameDisplayName string `json:"game_display_name,omitempty"`
	Developer       string `json:"developer,omitempty"`
}

// GameEvent represent a clock event request
type GameEvent struct {
	Game  string        `json:"game,omitempty"`
	Event string        `json:"event,omitempty"`
	Data  GameEventData `json:"data,omitempty"`
}

// GameEventData represent a clock event request
type GameEventData struct {
	Value int                `json:"value,omitempty"`
	Frame GameEventDataFrame `json:"frame,omitempty"`
}

// GameEventDataFrame represent a clock event request
type GameEventDataFrame struct {
	Date string `json:"date,omitempty"`
	Time string `json:"time,omitempty"`
}

// ServerDiscovery represent
type ServerDiscovery struct {
	Address string `json:"address,omitempty"`
}
