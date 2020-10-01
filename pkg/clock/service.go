package clock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kardianos/service"
)

// NewGameSenseClock create a GameSenseClock
func NewGameSenseClock() *GameSenseClock {

	return &GameSenseClock{
		Address:    "",
		DateFormat: "2006-01-02",
		TimeFormat: "15:04",
	}
}

// Init discover the game sense endpoint
func (c *GameSenseClock) Init(logger service.Logger) {
	c.Logger = logger
	for {
		if c.Address == "" {
			file, err := ioutil.ReadFile(configFile)
			if err != nil {
				c.Logger.Errorf("failed reading configuration file: %s", err)
				continue
			}

			serverDiscovery := ServerDiscovery{}
			err = json.Unmarshal(file, &serverDiscovery)
			if err != nil {
				c.Logger.Errorf("failed unmarshaling configuration file: %s", err)
				continue
			}

			c.Address = serverDiscovery.Address
			c.Logger.Infof("discovered address: %s", c.Address)

			return
		}

		c.Logger.Infof("will retry in 1 seconds")
		time.Sleep(1 * time.Second)
	}
}

// Register the clock in the engine app panel
func (c *GameSenseClock) Register() error {

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

	err := c.post("bind_game_event", bindingRequest)
	if err != nil {
		return err
	}

	gameMedataRequest := GameMetadata{
		Game:            "CLOCK",
		GameDisplayName: "Clock",
		Developer:       "Gael Mauleon",
	}

	return c.post("game_metadata", gameMedataRequest)

}

// Unregister the clock in the engine app panel
func (c *GameSenseClock) Unregister() error {
	gameRemovalRequest := GameMetadata{
		Game: "CLOCK",
	}
	return c.post("remove_game", gameRemovalRequest)
}

// Start the clock service
func (c *GameSenseClock) Start(service service.Service) error {

	c.Ticker = time.NewTicker(5 * time.Second)
	c.TickerDone = make(chan bool)

	go func() {
		for {

			select {
			case <-c.TickerDone:
				return
			case t := <-c.Ticker.C:
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

				c.post("game_event", gameEvent)
			}
		}
	}()

	c.Logger.Info("clock started")
	return nil
}

// Stop the clock service
func (c *GameSenseClock) Stop(service service.Service) error {
	c.Ticker.Stop()
	c.TickerDone <- true
	c.Logger.Info("clock stopped")

	return nil
}

func (c *GameSenseClock) post(path string, data interface{}) error {
	jsonPayload, _ := json.Marshal(data)
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(fmt.Sprintf("http://%s/%s", c.Address, path), "application/json", bytes.NewBuffer(jsonPayload))

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.Logger.Errorf("failed to read http response body: %s", err)
			return err
		}
		c.Logger.Errorf("failed http request: %s - %s", resp.Status, string(bodyBytes))
	}

	return nil
}
