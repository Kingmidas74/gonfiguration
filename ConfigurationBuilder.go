package gonfiguration

import (
	"errors"
	"fmt"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type ConfigurationBuilder struct {
	data map[string]string
}

func (this *ConfigurationBuilder) init() {
	if this.data == nil {
		this.data = make(map[string]string)
	}
}

func (this *ConfigurationBuilder) flatMap(prefix string, value interface{}, flatmap map[string]string) {
	submap, ok := value.(map[interface{}]interface{})
	if ok {
		for k, v := range submap {
			this.flatMap(prefix+Colon+k.(string), v, flatmap)
		}
		return
	}
	stringlist, ok := value.([]interface{})
	if ok {
		this.flatMap(fmt.Sprintf("%s.size", prefix), len(stringlist), flatmap)
		for i, v := range stringlist {
			this.flatMap(fmt.Sprintf("%s%s%d", prefix, Colon, i), v, flatmap)
		}
		return
	}
	flatmap[prefix] = fmt.Sprintf("%v", value)
}

func (this *ConfigurationBuilder) unflatMap() (nested map[string]interface{}, err error) {
	nested = make(map[string]interface{})

	for k, v := range this.data {
		temp := this.uf(k, v).(map[string]interface{})
		err = mergo.Merge(&nested, temp, func(c *mergo.Config) { c.Overwrite = true })
		if err != nil {
			return
		}
	}

	return
}

func (this *ConfigurationBuilder) uf(k string, v interface{}) (n interface{}) {
	n = v

	keys := strings.Split(k, Colon)

	for i := len(keys) - 1; i >= 0; i-- {
		temp := make(map[string]interface{})
		temp[keys[i]] = n
		n = temp
	}

	return
}

func (this *ConfigurationBuilder) AddEnvironmentVariables() (*ConfigurationBuilder, error) {
	this.init()
	for _, env := range os.Environ() {
		temp := strings.Split(env, "=")
		if len(temp) < 2 {
			return nil, errors.New("Wrong env variable format: " + env)
		}
		this.data[temp[0]] = strings.Join(temp[1:], "=")
	}
	return this, nil
}

func (this *ConfigurationBuilder) AddYamlFile(filePath string) (*ConfigurationBuilder, error) {
	this.init()
	m := make(map[string]interface{})
	data, err := ioutil.ReadFile(filePath)
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	o := map[string]string{}
	for k, v := range m {
		this.flatMap(k, v, o)
	}

	for k, v := range o {
		this.data[k] = v
	}
	return this, nil
}

func (this *ConfigurationBuilder) Build() (Configuration, error) {
	nested, err := this.unflatMap()
	if err != nil {
		return Configuration{}, err
	}
	return Configuration{data: nested}, nil
}
