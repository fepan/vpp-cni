// Copyright 2017 CNI authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This is a sample chained plugin that supports multiple CNI versions. It
// parses prevResult according to the cniVersion
package main

import (
	"encoding/json"
	"fmt"
	//"net"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/cni/pkg/version"
)

// PluginConf is whatever you expect your configuration json to be. This is whatever
// is passed in on stdin. Your plugin may wish to expose its functionality via
// runtime args, see CONVENTIONS.md in the CNI spec.
type PluginConf struct {
	types.NetConf // You may wish to not nest this type
	RuntimeConfig *struct {
		SampleConfig map[string]interface{} `json:"sample"`
	} `json:"runtimeConfig"`

	// This is the previous result, when called in the context of a chained
	// plugin. Because this plugin supports multiple versions, we'll have to
	// parse this in two passes. If your plugin is not chained, this can be
	// removed (though you may wish to error if a non-chainable plugin is
	// chained.
	// If you need to modify the result before returning it, you will need
	// to actually convert it to a concrete versioned struct.
	// RawPrevResult *map[string]interface{} `json:"prevResult"`
	// PrevResult    *current.Result         `json:"-"`

	// Add plugin-specifc flags here
	MyAwesomeFlag     bool   `json:"myAwesomeFlag"`
	AnotherAwesomeArg string `json:"anotherAwesomeArg"`
}

// parseConfig parses the supplied configuration (and prevResult) from stdin.
func parseConfig(stdin []byte) (*PluginConf, error) {
	conf := PluginConf{}

	if err := json.Unmarshal(stdin, &conf); err != nil {
		return nil, fmt.Errorf("failed to parse network configuration: %v", err)
	}

	if conf.AnotherAwesomeArg == "" {
		return nil, fmt.Errorf("anotherAwesomeArg must be specified")
	}

	return &conf, nil
}

// cmdAdd is called for ADD requests
func cmdAdd(args *skel.CmdArgs) error {
	conf, err := parseConfig(args.StdinData)
	fmt.Printf("cmdAdd has been called\n")
        if err != nil {
		return err
	}

	pass_json_result := `{
                              "ip4": {
                                      "ip": "10.15.20.2/24",
                                      "gateway": "10.15.20.1",
                                      "routes": [
                                                 {
                                                  "dst": "0.0.0.0/0"
                                                 },
                                                 {
                                                  "dst": "1.1.1.1/32",
                                                  "gw": "10.15.20.1"
                                                 }
                                                ]
                                     },
                               "dns": {}
			     }`
       rawIn := json.RawMessage(pass_json_result)
       resultBytes, err := rawIn.MarshalJSON()
       fmt.Println(string(resultBytes))
       if err != nil {
	 fmt.Printf("Cannot convert JSON to Bytes") 
         return err
       }
       res, err := version.NewResult(conf.CNIVersion, resultBytes)
       if err != nil {
	    fmt.Printf("Cannot create NewResult")
	    return err	
       }
       result, err := current.NewResultFromResult(res)
       if err != nil {
	   fmt.Printf("Cannot create NewResult from Result")
           return err
       } 
       return types.PrintResult(result, conf.CNIVersion)
}

// cmdDel is called for DELETE requests
func cmdDel(args *skel.CmdArgs) error {
	conf, err := parseConfig(args.StdinData)
	if err != nil {
		return err
	}
	_ = conf

	// Do your delete here

	return nil
}

func main() {
	skel.PluginMain(cmdAdd, cmdDel, version.PluginSupports("", "0.1.0", "0.2.0", version.Current()))
}
