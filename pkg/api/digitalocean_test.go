package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_digitalOceanService_ListVolumes(t *testing.T) {
	token := "my-token"
	factory := newDigitalOceanServiceFactory()
	doSvc := factory.Build(token).(*digitalOceanService)

	volumes := []Volume{
		Volume{ID: "3164444", Name: "example.com", Size: 25, Region: "nyc3"},
		Volume{ID: "95874511", Name: "my-other-droplet", Size: 50, Region: "nyc1"},
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.Path, "/droplets")

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"droplets": [
			  {
				"id": 3164444,
				"name": "example.com",
				"memory": 1024,
				"vcpus": 1,
				"disk": 25,
				"locked": false,
				"status": "active",
				"volume_ids": [				],
				"size": {				},
				"size_slug": "s-1vcpu-1gb",
				"region": {
					"name": "New York 3",
					"slug": "nyc3",
					"sizes": [
			
					],
					"features": [
					  "virtio",
					  "private_networking",
					  "backups",
					  "ipv6",
					  "metadata"
					],
					"available": null
				  }
			  },
			  {
				"id": 95874511,
				"name": "my-other-droplet",
				"memory": 2048,
				"vcpus": 1,
				"disk": 50,
				"locked": false,
				"status": "active",
				"volume_ids": [				],
				"size": {				},
				"size_slug": "s-1vcpu-1gb",
				"region": {
					"name": "New York 1",
					"slug": "nyc1",
					"sizes": [
			
					],
					"features": [
					  "virtio",
					  "private_networking",
					  "backups",
					  "ipv6",
					  "metadata"
					],
					"available": null
				  }
			  }
			]
		  }`))
	}))
	defer s.Close()

	doSvc.rootURL = s.URL

	myVolumes, err := doSvc.ListVolumes()
	assert.NoError(t, err)

	assert.ElementsMatch(t, myVolumes, volumes)
}

func Test_digitalOceanService_TakeSnapshot(t *testing.T) {
	token := "my-token"
	factory := newDigitalOceanServiceFactory()
	doSvc := factory.Build(token).(*digitalOceanService)

	volumeID := "vol-3"

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.Path, "/droplets/"+volumeID+"/actions")

		var reqJSON doTakeSnapshotReq
		err := json.NewDecoder(r.Body).Decode(&reqJSON)
		assert.NoError(t, err)
		assert.Equal(t, reqJSON.Type, "snapshot")

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"action": {
				"id": 36805022,
				"status": "in-progress",
				"type": "snapshot",
				"started_at": "2014-11-14T16:34:39Z",
				"completed_at": null,
				"resource_id": 3164450,
				"resource_type": "droplet",
				"region": {
					"name": "New York 3",
					"slug": "nyc3",
					"sizes": [
						"s-1vcpu-3gb",
						"m-1vcpu-8gb",
						"s-3vcpu-1gb",
						"s-1vcpu-2gb",
						"s-2vcpu-2gb",
						"s-2vcpu-4gb",
						"s-4vcpu-8gb",
						"s-6vcpu-16gb",
						"s-8vcpu-32gb",
						"s-12vcpu-48gb",
						"s-16vcpu-64gb",
						"s-20vcpu-96gb",
						"s-1vcpu-1gb",
						"c-1vcpu-2gb",
						"s-24vcpu-128gb"
					],
					"features": [
						"private_networking",
						"backups",
						"ipv6",
						"metadata",
						"server_id",
						"install_agent",
						"storage",
						"image_transfer"
					],
					"available": true
				},
				"region_slug": "nyc3"
			}
		}`))
	}))
	defer s.Close()

	doSvc.rootURL = s.URL

	providerSnapshotID, err := doSvc.TakeSnapshot(volumeID)
	assert.NoError(t, err)

	assert.Equal(t, providerSnapshotID, "36805022")
}
