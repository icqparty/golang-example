package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigStruct struct {
	Port string `yaml:"port"`
}

func readConfig(filepath string) ConfigStruct {
	ymlFile, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	var config ConfigStruct
	err = yaml.Unmarshal(ymlFile, &config)
	if err != nil {
		panic(err)
	}
	return config
}

func main() {
	var configPath string
	configDefaultPath:= "config.yaml"
	flag.StringVar(&configPath, "f", "", "путь к файлу конфигурации")
	flag.Parse()

	if len(configPath) > 0 {
		log.Printf("Файл конфигурации: %s\n", configPath)
	} else {
		configPath = configDefaultPath
		log.Println("Файл конфигурации не указан берем поу молчанию config.yml.")
	}

	config := readConfig(configPath)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		fmt.Fprintf(w, "Hello, %s!", name)
	})
	log.Printf("Server is listening on port %s ...", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
