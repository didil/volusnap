package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/didil/volusnap/pkg/models"
	"github.com/stretchr/testify/assert"
)

func Test_scalewayService_ListVolumes(t *testing.T) {
	token := "my-token"
	factory := newScalewayServiceFactory()
	scalewaySvc := factory.Build(token).(*scalewayService)

	volumes := []volume{
		volume{ID: "f929fe39-63f8-4be8-a80e-1e9c8ae22a76", Name: "volume-0-1", Size: 10, Region: "ams1", Organization: "000a115d-2852-4b0a-9ce8-47f1134ba95a"},
		volume{ID: "0facb6b5-b117-441a-81c1-f28b1d723779", Name: "volume-0-2", Size: 20, Region: "ams1", Organization: "000a115d-2852-4b0a-9ce8-47f1134ba95a"},
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.Path, "/volumes")
		assert.Equal(t, token, r.Header.Get("X-Auth-Token"))

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

func Test_scalewayService_TakeSnapshot(t *testing.T) {
	token := "my-token"
	factory := newScalewayServiceFactory()
	scalewaySvc := factory.Build(token).(*scalewayService)

	snapRule := &models.SnapRule{
		VolumeID:           "vol-3",
		VolumeOrganization: "org-11",
		VolumeName:         "myvolu",
		VolumeRegion:       "ams1",
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.Path, "/snapshots")
		assert.Equal(t, token, r.Header.Get("X-Auth-Token"))

		var reqJSON scalewayTakeSnapshotReq
		err := json.NewDecoder(r.Body).Decode(&reqJSON)
		assert.NoError(t, err)
		assert.Equal(t, reqJSON.VolumeID, snapRule.VolumeID)
		assert.Equal(t, reqJSON.Organization, snapRule.VolumeOrganization)
		assert.Contains(t, reqJSON.Name, "volusnap-"+snapRule.VolumeName+"-")

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"snapshot": {
			  "base_volume": {
				"id": "701a8946-ff9d-4579-95e3-1c2c2d0f892d",
				"name": "vol simple snapshot"
			  },
			  "creation_date": "2014-05-22T12:10:05.596769+00:00",
			  "id": "f0361e7b-cbe4-4882-a999-945192b7171b",
			  "name": "snapshot-0-1",
			  "organization": "000a115d-2852-4b0a-9ce8-47f1134ba95a",
			  "size": 10000000000,
			  "state": "snapshotting",
			  "volume_type": "l_ssd"
			}
		  }`))
	}))
	defer s.Close()

	scalewaySvc.rootURLs = map[string]string{"ams1": s.URL}

	providerSnapshotID, err := scalewaySvc.TakeSnapshot(snapRule)
	assert.NoError(t, err)

	assert.Equal(t, providerSnapshotID, "f0361e7b-cbe4-4882-a999-945192b7171b")
}
