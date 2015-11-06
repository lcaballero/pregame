package build

import (
	"github.com/spf13/viper"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"log"
)

type Yaml []byte


func (yaml Yaml) Load() *Value {
	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewBuffer(yaml))
	if err != nil {
		log.Fatal(err)
	}
	return NewValue(v)
}

func ReadYaml(file string) Yaml {
	yaml, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return yaml
}
