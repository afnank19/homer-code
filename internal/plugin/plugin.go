package plugin

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Plugin struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Command     string `json:"command"`
}

const PLUGINS_PATH = "homer-code/plugins"

func LoadPlugins() []Plugin {
	var plugins = []Plugin{}

	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	pluginsDir := filepath.Join(configDir, PLUGINS_PATH)

	dir, err := os.Open(pluginsDir)
	if err != nil {
		panic(err)
	}
	defer dir.Close()

	// List contents
	files, err := dir.Readdir(-1)
	if err != nil {
		panic(err)
	}

	fmt.Println("Plugins folder:", pluginsDir)
	for _, f := range files {
		fmt.Println("Reading", pluginsDir+"/"+f.Name())
		file, err := os.Open(pluginsDir + "/" + f.Name())
		if err != nil {
			panic(err)
		}
		defer file.Close()

		byteValue, _ := io.ReadAll(file)

		var pluginJson Plugin
		if err := json.Unmarshal(byteValue, &pluginJson); err != nil {
			panic(err)
		}

		fmt.Println("APPENDING: ", pluginJson)
		plugins = append(plugins, pluginJson)
	}

	return plugins
}
