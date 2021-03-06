package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/didil/volusnap/pkg/models"
)

func init() {
	// self register with provider registry
	pRegistry.register("digital_ocean", newDigitalOceanServiceFactory())
}

func newDigitalOceanServiceFactory() *digitalOceanServiceFactory {
	return &digitalOceanServiceFactory{}
}

type digitalOceanServiceFactory struct{}

func (factory *digitalOceanServiceFactory) Build(token string) ProviderSvcer {
	return &digitalOceanService{token: token, rootURL: "https://api.digitalocean.com/v2"}
}

type digitalOceanService struct {
	token   string
	rootURL string
}

func (svc *digitalOceanService) ListVolumes() ([]Volume, error) {
	req, err := http.NewRequest(http.MethodGet, svc.rootURL+"/droplets", nil)
	if err != nil {
		return nil, fmt.Errorf("DO list droplets NewRequest err: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+svc.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "VoluSnap")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("DO list droplets req err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("DO list droplets %v : %v", resp.Status, string(body))
	}

	type dropletRegion struct {
		Slug string `json:"slug,omitempty"`
	}

	type droplet struct {
		ID     float64       `json:"id,omitempty"`
		Name   string        `json:"name,omitempty"`
		Disk   float64       `json:"disk,omitempty"`
		Region dropletRegion `json:"region,omitempty"`
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
			ID:     strconv.Itoa(int(d.ID)),
			Name:   d.Name,
			Size:   d.Disk,
			Region: d.Region.Slug,
		})
	}

	return volumes, nil
}

type doTakeSnapshotReq struct {
	Type string `json:"type,omitempty"`
}

func (svc *digitalOceanService) TakeSnapshot(snapRule *models.SnapRule) (string, error) {
	var r bytes.Buffer
	json.NewEncoder(&r).Encode(&doTakeSnapshotReq{Type: "snapshot"})

	req, err := http.NewRequest(http.MethodPost, svc.rootURL+"/droplets/"+snapRule.VolumeID+"/actions", &r)
	if err != nil {
		return "", fmt.Errorf("DO TakeSnapshot NewRequest err: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+svc.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "VoluSnap")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("DO TakeSnapshot req err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("DO TakeSnapshot %v : %v", resp.Status, string(body))
	}

	type action struct {
		ID     float64 `json:"id,omitempty"`
		Status string  `json:"status,omitempty"`
	}

	type actionResp struct {
		Action action `json:"action,omitempty"`
	}

	var a actionResp

	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("DO TakeSnapshot json decode err: %v , body: %v", err, body)
	}

	providerSnapshotID := strconv.Itoa(int(a.Action.ID))

	return providerSnapshotID, nil
}
