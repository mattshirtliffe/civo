package api

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/absolutedevops/civo/config"
	"github.com/google/go-querystring/query"
	"github.com/jeffail/gabs"
)

type InstanceParams struct {
	Name        string `url:"hostname"`
	Size        string `url:"size"`
	Region      string `url:"region"`
	SSHKeyID    string `url:"ssh_key_id"`
	Template    string `url:"template"`
	InitialUser string `url:"initial_user"`
	NetworkID   string `url:"network_id"`
	FirewallID  string `url:"firewall_id"`
	Tags        string `url:"tags"`
	PublicIP    bool   `url:"public_ip"`
}

func InstancesList(tags string) (json *gabs.Container, err error) {
	if tags != "" {
		return makeJSONCall(config.URL()+"/v2/instances?per_page=10000000&tags="+tags, HTTPGet, "")
	} else {
		return makeJSONCall(config.URL()+"/v2/instances?per_page=10000000", HTTPGet, "")
	}
}

func InstanceCreate(params InstanceParams) (json *gabs.Container, err error) {
	v, _ := query.Values(params)
	return makeJSONCall(config.URL()+"/v2/instances", HTTPPost, v.Encode())
}

func InstanceReboot(id string, hard bool) (json *gabs.Container, err error) {
	if hard {
		return makeJSONCall(config.URL()+"/v2/instances/"+id+"/hard_reboots", HTTPPost, "")
	} else {
		return makeJSONCall(config.URL()+"/v2/instances/"+id+"/soft_reboots", HTTPPost, "")
	}
}

func InstanceDestroy(id string) (json *gabs.Container, err error) {
	return makeJSONCall(config.URL()+"/v2/instances/"+id, HTTPDelete, "")
}

func InstanceRestore(id, snapshot string) (json *gabs.Container, err error) {
	return makeJSONCall(config.URL()+"/v2/instances/"+id+"/restore", HTTPPut, "snapshot="+snapshot)
}

func InstanceRebuild(id string) (json *gabs.Container, err error) {
	return makeJSONCall(config.URL()+"/v2/instances/"+id+"/rebuild", HTTPPut, "")
}

func InstanceFirewall(id, firewall string) (json *gabs.Container, err error) {
	return makeJSONCall(config.URL()+"/v2/instances/"+id+"/firewall", HTTPPut, "firewall_id="+firewall)
}

func InstanceTags(id, tags string) (json *gabs.Container, err error) {
	return makeJSONCall(config.URL()+"/v2/instances/"+id+"/tags", HTTPPut, "tags="+tags)
}

func InstanceUpgrade(id, size string) (json *gabs.Container, err error) {
	return makeJSONCall(config.URL()+"/v2/instances/"+id+"/resize", HTTPPut, "size="+size)
}

// Utility functions ---------------------------------------------------------------------------------------------------

func InstanceFind(search string) string {
	instances, err := InstancesList("")
	if err != nil {
		fmt.Println("DEBUG: Returning early because err is", err)
		return ""
	}
	items, _ := instances.S("items").Children()
	for _, child := range items {
		id := child.S("id").Data().(string)
		name := child.S("hostname").Data().(string)
		if strings.Contains(id, search) {
			return id
		}
		if strings.Contains(name, search) {
			return id
		}
	}
	return ""
}

var ADJECTIVES = []string{
	"autumn", "hidden", "bitter", "misty", "silent", "empty", "dry", "dark",
	"summer", "icy", "delicate", "quiet", "white", "cool", "spring", "winter",
	"patient", "twilight", "dawn", "crimson", "wispy", "weathered", "blue",
	"billowing", "broken", "cold", "damp", "falling", "frosty", "green",
	"long", "late", "lingering", "bold", "little", "morning", "muddy", "old",
	"red", "rough", "still", "small", "sparkling", "throbbing", "shy",
	"wandering", "withered", "wild", "black", "young", "holy", "solitary",
	"fragrant", "aged", "snowy", "proud", "floral", "restless", "divine",
	"polished", "ancient", "purple", "lively", "nameless", "lucky", "odd", "tiny",
	"free", "dry", "yellow", "orange", "gentle", "tight", "super", "royal", "broad",
	"steep", "flat", "square", "round", "mute", "noisy", "hushy", "raspy", "soft",
	"shrill", "rapid", "sweet", "curly", "calm", "jolly", "fancy", "plain", "shinny",
}

var NOUNS = []string{
	"waterfall", "river", "breeze", "moon", "rain", "wind", "sea", "morning",
	"snow", "lake", "sunset", "pine", "shadow", "leaf", "dawn", "glitter",
	"forest", "hill", "cloud", "meadow", "sun", "glade", "bird", "brook",
	"butterfly", "bush", "dew", "dust", "field", "fire", "flower", "firefly",
	"feather", "grass", "haze", "mountain", "night", "pond", "darkness",
	"snowflake", "silence", "sound", "sky", "shape", "surf", "thunder",
	"violet", "water", "wildflower", "wave", "water", "resonance", "sun",
	"wood", "dream", "cherry", "tree", "fog", "frost", "voice", "paper",
	"frog", "smoke", "star", "atom", "band", "bar", "base", "block", "boat",
	"term", "credit", "art", "fashion", "truth", "disk", "math", "unit", "cell",
	"scene", "heart", "recipe", "union", "limit", "bread", "toast", "bonus",
	"lab", "mud", "mode", "poetry", "tooth", "hall", "king", "queen", "lion", "tiger",
	"penguin", "kiwi", "cake", "mouse", "rice", "coke", "hola", "salad", "hat",
}

func InstanceSuggestHostname() string {
	rand.Seed(time.Now().UTC().UnixNano())
	return ADJECTIVES[rand.Intn(len(ADJECTIVES))] + "-" + NOUNS[rand.Intn(len(NOUNS))] + ".example.com"
}
