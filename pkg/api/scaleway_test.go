package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_scalewayService_ListVolumes(t *testing.T) {
	token := "my-token"
	factory := newScalewayServiceFactory()
	scalewaySvc := factory.Build(token).(*scalewayService)

	volumes := []Volume{
		Volume{ID: "f929fe39-63f8-4be8-a80e-1e9c8ae22a76", Name: "volume-0-1", Size: 10, Region: "ams1"},
		Volume{ID: "0facb6b5-b117-441a-81c1-f28b1d723779", Name: "volume-0-2", Size: 20, Region: "ams1"},
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.Path, "/volumes")

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"volumes": [
				{
					"export_uri": null,
					"id": "f929fe39-63f8-4be8-a80e-1e9c8ae22a76",
					"name": "volume-0-1",
					"organization": "000a115d-2852-4b0a-9ce8-47f1134ba95a",
					"server": null,
					"size": 10000000000,
					"volume_type": "l_ssd"
				},
				{
					"export_uri": null,
					"id": "0facb6b5-b117-441a-81c1-f28b1d723779",
					"name": "volume-0-2",
					"organization": "000a115d-2852-4b0a-9ce8-47f1134ba95a",
					"server": null,
					"size": 20000000000,
					"volume_type": "l_ssd"
				}
			]
		}`))
	}))
	defer s.Close()

	scalewaySvc.rootURLs = map[string]string{"ams1": s.URL}

	myVolumes, err := scalewaySvc.ListVolumes()
	assert.NoError(t, err)

	assert.ElementsMatch(t, myVolumes, volumes)
}
