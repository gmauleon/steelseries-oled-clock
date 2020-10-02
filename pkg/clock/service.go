package clock

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/kardianos/service"
)

// NewGameSenseClockService create a GameSenseClock service
func NewGameSenseClockService() *GameSenseClockService {

	svcConfig := &service.Config{
		Name:        "SteelSeriesOLEDClock",
		DisplayName: "SteelSeries OLED Clock",
		Description: "This service communicate with GameSense to display the current date and time.",
		Arguments: []string{
			"run",
		},
	}

	gsc := &GameSenseClockService{
		DateFormat: "2006-01-02",
		TimeFormat: "15:04",
	}

	svc, err := service.New(gsc, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := svc.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	gsc.Service = svc
	gsc.Logger = logger
	return gsc
}

// RunService run the underlying service
func (c *GameSenseClockService) RunService() {

	if err := c.Service.Run(); err != nil {
		c.Logger.Error(err)
	}
}

// InstallService the clock in the engine app panel
func (c *GameSenseClockService) InstallService() {
	apiAddress, err := c.discoverGameSenseAPI()
	if err != nil {
		c.Logger.Error(err)
		return
	}

	if err := c.Service.Install(); err != nil {
		c.Logger.Error(err)
	} else {
		c.Logger.Info("service installed")
	}

	bindingRequest := Binding{
		Game:          "CLOCK",
		Event:         "TIME",
		ValueOptional: true,
		Handlers: []BindingHandler{
			{
				DeviceType: "screened",
				Zone:       "one",
				Mode:       "screen",
				Datas: []BindingHandlerData{
					{
						IconID: 15,
						Lines: []BindingHandlerDataLine{
							{
								HasText:         true,
								Bold:            true,
								ContextFrameKey: "date",
							},
							{
								HasText:         true,
								Bold:            true,
								ContextFrameKey: "time",
							},
						},
					},
				},
			},
		},
	}

	if err := c.post(apiAddress, "bind_game_event", bindingRequest); err != nil {
		c.Logger.Error(err)
	} else {
		c.Logger.Info("clock added to GameSense")
	}

	gameMedataRequest := GameMetadata{
		Game:            "CLOCK",
		GameDisplayName: "Clock",
		Developer:       "Gael Mauleon",
	}

	if err := c.post(apiAddress, "game_metadata", gameMedataRequest); err != nil {
		c.Logger.Error(err)
	} else {
		c.Logger.Info("clock configured")
	}
}

// UninstallService the clock in the engine app panel
func (c *GameSenseClockService) UninstallService() {
	apiAddress, err := c.discoverGameSenseAPI()
	if err != nil {
		c.Logger.Error(err)
		return
	}

	gameRemovalRequest := GameMetadata{
		Game: "CLOCK",
	}
	if err := c.post(apiAddress, "remove_game", gameRemovalRequest); err != nil {
		c.Logger.Error(err)
	} else {
		c.Logger.Info("clock removed from GameSense")
	}

	if err := c.Service.Stop(); err != nil {
		c.Logger.Error(err)
	} else {
		c.Logger.Info("service stopped")
	}

	if err := c.Service.Uninstall(); err != nil {
		c.Logger.Error(err)
	} else {
		c.Logger.Info("service uninstalled")
	}
}

// Start statisfy the Start service interface
func (c *GameSenseClockService) Start(service service.Service) error {

	c.Ticker = time.NewTicker(5 * time.Second)
	c.TickerDone = make(chan bool)

	go func() {
		var err error
		var apiAddress string

		for {

			select {
			case <-c.TickerDone:
				return
			case t := <-c.Ticker.C:

				if apiAddress == "" {
					apiAddress, err = c.discoverGameSenseAPI()
					if err != nil {
						c.Logger.Error(err)
						continue
					}
				}

				gameEvent := GameEvent{
					Game:  "CLOCK",
					Event: "TIME",
					Data: GameEventData{
						Value: 1,
						Frame: GameEventDataFrame{
							Date: t.Format("2006-01-02"),
							Time: t.Format("15:04"),
						},
					},
				}

				if err = c.post(apiAddress, "game_event", gameEvent); err != nil {
					c.Logger.Error(err)
				}
			}
		}
	}()

	c.Logger.Info("clock started")
	return nil
}

// Stop satisfy the service interface
func (c *GameSenseClockService) Stop(service service.Service) error {
	c.Ticker.Stop()
	c.TickerDone <- true
	c.Logger.Info("clock stopped")

	return nil
}

func (c *GameSenseClockService) post(apiAddress string, path string, data interface{}) error {
	jsonPayload, _ := json.Marshal(data)
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(fmt.Sprintf("http://%s/%s", apiAddress, path), "application/json", bytes.NewBuffer(jsonPayload))

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed http request: %s - %s", resp.Status, string(bodyBytes))
	}

	return nil
}

func (c *GameSenseClockService) discoverGameSenseAPI() (string, error) {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return "", fmt.Errorf("failed reading GameSense configuration file: %s", err)
	}

	serverDiscovery := ServerDiscovery{}
	err = json.Unmarshal(file, &serverDiscovery)
	if err != nil {
		return "", fmt.Errorf("failed unmarshaling GameSense configuration file: %s", err)
	}

	conn, err := net.Dial("tcp", serverDiscovery.Address)
	if err != nil {
		return "", errors.New("GameSense API address is invalid")
	}
	conn.Close()

	c.Logger.Infof("discovered GameSense API address: %s", serverDiscovery.Address)
	return serverDiscovery.Address, nil
}
