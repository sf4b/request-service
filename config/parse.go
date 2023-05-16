package config

import (
	"encoding/json"
	"errors"
	"github.com/caarlos0/env/v6"
	"github.com/ghodss/yaml"
	"github.com/joho/godotenv"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

const (
	tagName        = "json"
	prepareTagName = "prepare"
	optValue       = "optional"
)

// Parse прочитать конфиг и выполнить Preparer на его полях.
func Parse(configFile, envFile string, target interface{}) error {
	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	var isYaml bool
	index := strings.LastIndex(configFile, ".")
	if index > 0 {
		if ext := configFile[index+1:]; ext == "yml" || ext == "yaml" {
			isYaml = true
		}
	}
	if isYaml {
		if err = yaml.Unmarshal(data, target); err != nil {
			return err
		}
	} else {
		if err = json.Unmarshal(data, target); err != nil {
			return err
		}
	}
	if err = enrichWithEnv(envFile, target); err != nil {
		return err
	}
	if err = Prepare(target); err != nil {
		return err
	}
	return nil
}

func enrichWithEnv(envFile string, target interface{}) error {
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			return err
		}
	}
	if err := env.Parse(target); err != nil {
		return err
	}
	return nil
}

// Preparer имплементируется структурами, которые формируют конфиги сервисов.
type Preparer interface {
	// Prepare валидирует структуру конфига.
	Prepare() error
}

// Prepare рекурсивно вызывает методы валидации у структуры
// конфига и её составных частей.
func Prepare(src interface{}) error {
	if src == nil {
		return nil
	}

	v := reflect.ValueOf(src)

	pr, ok := src.(Preparer)
	if ok {
		err := pr.Prepare()
		if err != nil {
			return err //nolint:wrapcheck
		}
	}
	return traverse(v, true)
}

func traverse(v reflect.Value, parentTraversed bool) (err error) {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		if !v.IsNil() && v.CanInterface() {
			if err := tryPrepareInterface(v.Interface()); err != nil {
				return err
			}
			if err := traverse(v.Elem(), true); err != nil {
				return err
			}
		}
	case reflect.Struct:
		if !parentTraversed && v.CanInterface() {
			if err := tryPrepareInterface(v.Interface()); err != nil {
				return err
			}
		}
		for j := 0; j < v.NumField(); j++ {
			optTag := v.Type().Field(j).Tag.Get(prepareTagName)
			if optTag == optValue && v.Field(j).IsNil() {
				continue
			}

			err := traverse(v.Field(j), false)
			if err != nil {
				tagValue := v.Type().Field(j).Tag.Get(tagName)
				return errors.New("wrong tag: " + tagValue + " ,cause:" + err.Error())
			}
		}
	default:
		if v.CanInterface() {
			return tryPrepareInterface(v.Interface())
		}
	}
	return nil
}

func tryPrepareInterface(v interface{}) (err error) {
	pr, ok := v.(Preparer)
	if ok {
		err = pr.Prepare()
	}
	return
}
