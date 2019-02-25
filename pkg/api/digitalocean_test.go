package api

import (
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
		Volume{ID: "3164444", Name: "example.com", Size: 25},
		Volume{ID: "95874511", Name: "my-other-droplet", Size: 50},
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
				"size_slug": "s-1vcpu-1gb"
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
				"size_slug": "s-1vcpu-1gb"
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
