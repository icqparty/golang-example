package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        next.ServeHTTP(w, r)
        
        log.Printf(
            "%s\t%s\t%s\t%s",
            r.Method,
            r.URL.Path,
            r.RemoteAddr,
            time.Since(start),
        )
    }
}


func handlerHome(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	}

func handlerHello(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		fmt.Fprintf(w, "Hello, %s!", name)
	}

func main() {
	var configPath string
	flag.StringVar(&configPath, "f", "", "путь к файлу конфигурации")
	flag.Parse()

	// Проверяем наличие пути к файлу конфигурации
	if len(configPath) == 0 {
		log.Fatalf( "Не указан файл конфигурации.")
		os.Exit(1)
	}

	// Читаем конфигурационный файл
	config, err := readConfig(configPath)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла конфигурации: %v :", err)
		os.Exit(1)
	}


    http.HandleFunc("/", LoggerMiddleware(handlerHome))
    http.HandleFunc("/hello", LoggerMiddleware(handlerHello))

	log.Printf("Сервер запущен и слушает порт %s ...", config.Port)
	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		log.Fatalf("Ошибка сервера: %v\n", err)
		os.Exit(1)
	}
    

}
