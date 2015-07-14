// Copyright 2015 Google Inc. All Rights Reserved.
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

package collector

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyConfig(t *testing.T) {
	assert := assert.New(t)

	emptyConfig := `
        {
                "endpoint" : "http://localhost:8000/nginx_status",
                "metrics_config"  : [
                ]
        }
        `

	//Create a temporary config file 'temp.json' with invalid json format
	assert.NoError(ioutil.WriteFile("temp.json", []byte(emptyConfig), 0777))

	_, err := NewCollector("tempCollector", "temp.json")
	assert.Error(err)

	assert.NoError(os.Remove("temp.json"))
}

func TestConfigWithErrors(t *testing.T) {
	assert := assert.New(t)

	//Syntax error: Missed '"' after activeConnections
	invalid := `
	{
		"endpoint" : "http://localhost:8000/nginx_status",
		"metrics_config"  : [
			{
				 "name" : "activeConnections,  
		  		 "metric_type" : "gauge",
		 	 	 "units" : "integer",
		  		 "polling_frequency" : 10,
		    		 "regex" : "Active connections: ([0-9]+)"			
			}
		]
	}
	`

	//Create a temporary config file 'temp.json' with invalid json format
	assert.NoError(ioutil.WriteFile("temp.json", []byte(invalid), 0777))

	_, err := NewCollector("tempCollector", "temp.json")
	assert.Error(err)

	assert.NoError(os.Remove("temp.json"))
}

func TestConfigWithRegexErrors(t *testing.T) {
	assert := assert.New(t)

	//Error: Missed operand for '+' in activeConnections regex
	invalid := `
        {
                "endpoint" : "host:port/nginx_status",
                "metrics_config"  : [
                        {
                                 "name" : "activeConnections",
                                 "metric_type" : "gauge",
                                 "units" : "integer",
                                 "polling_frequency" : 10,
                                 "regex" : "Active connections: (+)"
                        },
                        {
                                 "name" : "reading",
                                 "metric_type" : "gauge",
                                 "units" : "integer",
                                 "polling_frequency" : 10,
                                 "regex" : "Reading: ([0-9]+) .*"
                        }
                ]
        }
        `

	//Create a temporary config file 'temp.json'
	assert.NoError(ioutil.WriteFile("temp.json", []byte(invalid), 0777))

	_, err := NewCollector("tempCollector", "temp.json")
	assert.Error(err)

	assert.NoError(os.Remove("temp.json"))
}

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	//Create an nginx collector using the config file 'sample_config.json'
	collector, err := NewCollector("nginx", "config/sample_config.json")
	assert.NoError(err)
	assert.Equal(collector.name, "nginx")
	assert.Equal(collector.configFile.Endpoint, "http://localhost:8000/nginx_status")
	assert.Equal(collector.configFile.MetricsConfig[0].Name, "activeConnections")
}