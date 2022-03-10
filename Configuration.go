package gonfiguration

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"strings"
)

type Configuration struct {
	data map[string]interface{}
}

func (this *Configuration) GetSection(path string) (map[string]interface{}, error) {
	result := this.data
	for _, p := range strings.Split(path, Colon) {
		if result[p] == nil {
			return nil, errors.New("Key " + path + " doesn't exist")
		}
		if reflect.ValueOf(result[p]).Kind() == reflect.Map {
			result = result[p].(map[string]interface{})
		} else {
			return nil, errors.New(path + " there is not a section!")
		}
	}
	return result, nil
}

func (this *Configuration) GetValue(path string) (interface{}, error) {
	result := this.data
	for _, p := range strings.Split(path, Colon) {
		if result[p] == nil {
			return nil, errors.New("Key " + path + " doesn't exist")
		}
		if reflect.ValueOf(result[p]).Kind() == reflect.Map {
			result = result[p].(map[string]interface{})
		} else {
			return result[p], nil
		}
	}
	return nil, nil
}

func (this *Configuration) Bind(path string, result interface{}) error {
	temp := this.data
	for _, p := range strings.Split(path, Colon) {
		if temp[p] == nil {
			return errors.New("Key " + path + " doesn't exist")
		}
		if reflect.ValueOf(temp[p]).Kind() == reflect.Map {
			temp = temp[p].(map[string]interface{})
		} else {
			return errors.New(path + " there is not a section!")
		}
	}
	err := mapstructure.Decode(temp, &result)
	if err != nil {
		return err
	}
	return nil
}
