package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func init() {
	// self register with provider registry
	pRegistry.register("digital_ocean", newDigitalOceanServiceFactory())
}

func newDigitalOceanServiceFactory() *digitalOceanServiceFactory {
	return &digitalOceanServiceFactory{}
}

type digitalOceanServiceFactory struct{}

func (factory *digitalOceanServiceFactory) Build(token string) providerSvcer {
	return &digitalOceanService{token: token, rootURL: "https://api.digitalocean.com/v2"}
}

type digitalOceanService struct {
	token   string
	rootURL string
}

func (do *digitalOceanService) ListVolumes() ([]Volume, error) {
	req, err := http.NewRequest(http.MethodGet, do.rootURL+"/droplets", nil)
	if err != nil {
		return nil, fmt.Errorf("DO list droplets NewRequest err: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+do.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "VoluSnap")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("DO list droplets req err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("DO list droplets %v : %v", resp.Status, string(body))
	}

	type droplet struct {
		ID   float64 `json:"id,omitempty"`
		Name string  `json:"name,omitempty"`
		Disk float64 `json:"disk,omitempty"`
	}

	type dropletsList struct {
		Droplets []droplet `json:"droplets,omitempty"`
	}

	var b dropletsList

	err = json.NewDecoder(resp.Body).Decode(&b)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("DO list droplets json decode err: %v , body: %v", err, body)
	}

	var volumes []Volume
	droplets := b.Droplets
	for _, d := range droplets {
		volumes = append(volumes, Volume{
			ID:   strconv.Itoa(int(d.ID)),
			Name: d.Name,
			Size: d.Disk,
		})
	}

	return volumes, nil
}
