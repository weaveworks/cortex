package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/weaveworks/cortex/configs"
	"github.com/weaveworks/cortex/configs/api"
)

const (
	rulesEndpoint        = "/api/prom/rules"
	rulesPrivateEndpoint = "/private/api/prom/rules"

	alertConfigEndpoint        = "/api/prom/alerts"
	alertConfigPrivateEndpoint = "/private/api/prom/alerts"
)

var (
	rulesClient       = configurable{rulesEndpoint, rulesPrivateEndpoint}
	alertConfigClient = configurable{alertConfigEndpoint, alertConfigPrivateEndpoint}

	allClients = []configurable{rulesClient, alertConfigClient}
)

// The root page returns 200 OK.
func Test_Root_OK(t *testing.T) {
	setup(t)
	defer cleanup(t)

	w := request(t, "GET", "/", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

type configurable struct {
	Endpoint        string
	PrivateEndpoint string
}

// post a config
func (c configurable) post(t *testing.T, userID string, config configs.Config) configs.ConfigView {
	w := requestAsUser(t, userID, "POST", c.Endpoint, jsonObject(config).Reader(t))
	require.Equal(t, http.StatusNoContent, w.Code)
	return c.get(t, userID)
}

// get a config
func (c configurable) get(t *testing.T, userID string) configs.ConfigView {
	w := requestAsUser(t, userID, "GET", c.Endpoint, nil)
	return parseConfigView(t, w.Body.Bytes())
}

// configs returns 401 to requests without authentication.
func Test_GetConfig_Anonymous(t *testing.T) {
	setup(t)
	defer cleanup(t)

	for _, c := range allClients {
		w := request(t, "GET", c.Endpoint, nil)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	}
}

// configs returns 404 if no config has been created yet.
func Test_GetConfig_NotFound(t *testing.T) {
	setup(t)
	defer cleanup(t)

	userID := makeUserID()
	for _, c := range allClients {
		w := requestAsUser(t, userID, "GET", c.Endpoint, nil)
		assert.Equal(t, http.StatusNotFound, w.Code)
	}
}

// configs returns 401 to requests without authentication.
func Test_PostConfig_Anonymous(t *testing.T) {
	setup(t)
	defer cleanup(t)

	for _, c := range allClients {
		w := request(t, "POST", c.Endpoint, nil)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	}
}

// Posting to a configuration sets it so that you can get it again.
func Test_PostConfig_CreatesConfig(t *testing.T) {
	setup(t)
	defer cleanup(t)

	userID := makeUserID()
	config := makeConfig()
	content := jsonObject(config)
	for _, c := range allClients {
		{
			w := requestAsUser(t, userID, "POST", c.Endpoint, content.Reader(t))
			assert.Equal(t, http.StatusNoContent, w.Code)
		}
		{
			w := requestAsUser(t, userID, "GET", c.Endpoint, nil)
			assert.Equal(t, config, parseConfigView(t, w.Body.Bytes()).Config)
		}
	}
}

// Posting to a configuration sets it so that you can get it again.
func Test_PostConfig_UpdatesConfig(t *testing.T) {
	setup(t)
	defer cleanup(t)

	userID := makeUserID()
	for _, c := range allClients {
		view1 := c.post(t, userID, makeConfig())
		config2 := makeConfig()
		view2 := c.post(t, userID, config2)
		assert.True(t, view2.ID > view1.ID, "%v > %v", view2.ID, view1.ID)
		assert.Equal(t, config2, view2.Config)
	}
}

// Different users can have different configurations.
func Test_PostConfig_MultipleUsers(t *testing.T) {
	setup(t)
	defer cleanup(t)

	userID1 := makeUserID()
	userID2 := makeUserID()
	for _, c := range allClients {
		config1 := c.post(t, userID1, makeConfig())
		config2 := c.post(t, userID2, makeConfig())
		foundConfig1 := c.get(t, userID1)
		assert.Equal(t, config1, foundConfig1)
		foundConfig2 := c.get(t, userID2)
		assert.Equal(t, config2, foundConfig2)
		assert.True(t, config2.ID > config1.ID, "%v > %v", config2.ID, config1.ID)
	}
}

// GetAllConfigs returns an empty list of configs if there aren't any.
func Test_GetAllConfigs_Empty(t *testing.T) {
	setup(t)
	defer cleanup(t)

	for _, c := range allClients {
		w := request(t, "GET", c.PrivateEndpoint, nil)
		assert.Equal(t, http.StatusOK, w.Code)
		var found api.ConfigsView
		err := json.Unmarshal(w.Body.Bytes(), &found)
		assert.NoError(t, err, "Could not unmarshal JSON")
		assert.Equal(t, api.ConfigsView{Configs: map[string]configs.ConfigView{}}, found)
	}
}

// GetAllConfigs returns all created configs.
func Test_GetAllConfigs(t *testing.T) {
	setup(t)
	defer cleanup(t)

	userID := makeUserID()
	config := makeConfig()

	for _, c := range allClients {
		view := c.post(t, userID, config)
		w := request(t, "GET", c.PrivateEndpoint, nil)
		assert.Equal(t, http.StatusOK, w.Code)
		var found api.ConfigsView
		err := json.Unmarshal(w.Body.Bytes(), &found)
		assert.NoError(t, err, "Could not unmarshal JSON")
		assert.Equal(t, api.ConfigsView{Configs: map[string]configs.ConfigView{
			userID: view,
		}}, found)
	}
}

// GetAllConfigs returns the *newest* versions of all created configs.
func Test_GetAllConfigs_Newest(t *testing.T) {
	setup(t)
	defer cleanup(t)

	userID := makeUserID()

	for _, c := range allClients {
		c.post(t, userID, makeConfig())
		c.post(t, userID, makeConfig())
		lastCreated := c.post(t, userID, makeConfig())

		w := request(t, "GET", c.PrivateEndpoint, nil)
		assert.Equal(t, http.StatusOK, w.Code)
		var found api.ConfigsView
		err := json.Unmarshal(w.Body.Bytes(), &found)
		assert.NoError(t, err, "Could not unmarshal JSON")
		assert.Equal(t, api.ConfigsView{Configs: map[string]configs.ConfigView{
			userID: lastCreated,
		}}, found)
	}
}

func Test_GetConfigs_IncludesNewerConfigsAndExcludesOlder(t *testing.T) {
	setup(t)
	defer cleanup(t)

	for _, c := range allClients {
		c.post(t, makeUserID(), makeConfig())
		config2 := c.post(t, makeUserID(), makeConfig())
		userID3 := makeUserID()
		config3 := c.post(t, userID3, makeConfig())

		w := request(t, "GET", fmt.Sprintf("%s?since=%d", c.PrivateEndpoint, config2.ID), nil)
		assert.Equal(t, http.StatusOK, w.Code)
		var found api.ConfigsView
		err := json.Unmarshal(w.Body.Bytes(), &found)
		assert.NoError(t, err, "Could not unmarshal JSON")
		assert.Equal(t, api.ConfigsView{Configs: map[string]configs.ConfigView{
			userID3: config3,
		}}, found)
	}
}
