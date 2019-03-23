package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/didil/volusnap/pkg/models"
)

func init() {
	// self register with provider registry
	pRegistry.register("scaleway", newScalewayServiceFactory())
}

func newScalewayServiceFactory() *scalewayServiceFactory {
	return &scalewayServiceFactory{}
}

type scalewayServiceFactory struct{}

func (factory *scalewayServiceFactory) Build(token string) providerSvcer {
	return &scalewayService{
		token: token,
		rootURLs: map[string]string{
			"par1": "https://cp-par1.scaleway.com",
			"ams1": "https://cp-ams1.scaleway.com",
		},
	}
}

type scalewayService struct {
	token    string
	rootURLs map[string]string
}

func (svc *scalewayService) ListVolumes() ([]volume, error) {
	var volumes []volume

	for reg, rootURL := range svc.rootURLs {
		req, err := http.NewRequest(http.MethodGet, rootURL+"/volumes", nil)
		if err != nil {
			return nil, fmt.Errorf("Scaleway list volumes NewRequest err: %v", err)
		}

		req.Header.Set("X-Auth-Token", svc.token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "VoluSnap")

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Scaleway list volumes req err: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, fmt.Errorf("Scaleway list volumes %v : %v", resp.Status, string(body))
		}

		type scalewayVolume struct {
			ID           string  `json:"id,omitempty"`
			Name         string  `json:"name,omitempty"`
			Organization string  `json:"organization,omitempty"`
			Size         float64 `json:"size,omitempty"`
		}

		type volumesList struct {
			Volumes []scalewayVolume `json:"volumes,omitempty"`
		}

		var b volumesList

		err = json.NewDecoder(resp.Body).Decode(&b)
		if err != nil {
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, fmt.Errorf("Scaleway list volumes json decode err: %v , body: %v", err, body)
		}

		scalewayVolumes := b.Volumes
		for _, sVol := range scalewayVolumes {
			volumes = append(volumes, volume{
				ID:           sVol.ID,
				Name:         sVol.Name,
				Organization: sVol.Organization,
				Size:         sVol.Size / (math.Pow10(9)),
				Region:       reg,
			})
		}
	}

	return volumes, nil
}

type scalewayTakeSnapshotReq struct {
	VolumeID     string `json:"volume_id"`
	Name         string `json:"name"`
	Organization string `json:"organization"`
}

func (svc *scalewayService) TakeSnapshot(snapRule *models.SnapRule) (string, error) {
	var reqJSON bytes.Buffer

	json.NewEncoder(&reqJSON).Encode(&scalewayTakeSnapshotReq{
		VolumeID:     snapRule.VolumeID,
		Organization: snapRule.VolumeOrganization,
		Name:         "volusnap-" + snapRule.VolumeName + "-" + strconv.Itoa(int(time.Now().Unix())),
	})

	rootURL := svc.rootURLs[snapRule.VolumeRegion]
	if rootURL == "" {
		return "", fmt.Errorf("Scaleway rootURL not found for: %v", snapRule.VolumeRegion)
	}

	req, err := http.NewRequest(http.MethodPost, rootURL+"/snapshots", &reqJSON)
	if err != nil {
		return "", fmt.Errorf("Scaleway TakeSnapshot NewRequest err: %v", err)
	}

	req.Header.Set("X-Auth-Token", svc.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "VoluSnap")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Scaleway TakeSnapshot req err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("Scaleway TakeSnapshot %v : %v", resp.Status, string(body))
	}

	type scalewaySnapshot struct {
		ID    string `json:"id,omitempty"`
		State string `json:"state,omitempty"`
	}

	type actionResp struct {
		Snapshot scalewaySnapshot `json:"snapshot,omitempty"`
	}

	var a actionResp

	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("Scaleway TakeSnapshot json decode err: %v , body: %v", err, body)
	}

	providerSnapshotID := a.Snapshot.ID

	return providerSnapshotID, nil
}
