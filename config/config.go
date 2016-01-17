package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jeffail/gabs"
)

var Config *gabs.Container

func LoadConfig() {
	filename := os.ExpandEnv("$HOME/.civo.json")
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			createNewJSONConfig(filename)
		}
	}

	contents, err := ioutil.ReadFile(filename)
	Config, err = gabs.ParseJSON(contents)
	if err != nil {
		fmt.Println("Error ", err.Error())
	}
}

func createNewJSONConfig(filename string) {
	newConfig := gabs.New()
	newConfig.SetP(false, "meta.admin")
	newConfig.SetP("https://api.civo.com", "meta.url")
	newConfig.SetP("1", "meta.version")
	ioutil.WriteFile(filename, []byte(newConfig.String()), 0600)
}

func save() {
	filename := os.ExpandEnv("$HOME/.civo.json")
	ioutil.WriteFile(filename, []byte(Config.String()), 0600)
}

func getBool(path string) bool {
	value, _ := Config.Path(path).Data().(bool)
	return value
}

func getString(path string) string {
	value, _ := Config.Path(path).Data().(string)
	return value
}

func Admin() bool {
	if Config == nil {
		LoadConfig()
	}
	return getBool("meta.admin")
}

func URL() string {
	return getString("meta.url")
}

func CurrentToken() string {
	if currentTokenKey := getString("meta.current_token"); currentTokenKey != "" {
		return getString(fmt.Sprintf("tokens.%s", currentTokenKey))
	}
	tokens, _ := Config.S("tokens").ChildrenMap()
	for name, token := range tokens {
		Config.SetP(name, "meta.current_token")
		save()
		return token.Data().(string)
	}
	fmt.Println("You haven't got a token saved, ask your provider for one and save it using 'civo tokens:save'")
	os.Exit(-1)
	return ""
}
