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

func readConfig(filepath string) (*ConfigStruct, error) {
    ymlFile, err := os.ReadFile(filepath)
    if err != nil {
        return nil, fmt.Errorf("не удалось прочитать файл конфигурации: %w", err)
    }

    var config ConfigStruct
    err = yaml.Unmarshal(ymlFile, &config)
    if err != nil {
        return nil, fmt.Errorf("не удалось распарсить YAML-файл: %w", err)
    }

    return &config, nil
}

func main() {
	var configPath string
    flag.StringVar(&configPath, "f", "", "путь к файлу конфигурации")
    flag.Parse()

    // Проверяем наличие пути к файлу конфигурации
    if len(configPath) == 0 {
        fmt.Fprintln(os.Stderr, "Ошибка: не указан файл конфигурации.")
        os.Exit(1)
    }

    // Читаем конфигурационный файл
    config, err := readConfig(configPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Ошибка при чтении файла конфигурации: %v\n", err)
        os.Exit(1)
    }

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
