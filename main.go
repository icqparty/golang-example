package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

var requestCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total1",
		Help: "Количество запросов",
	}, []string{"code"})

func init() {
	prometheus.MustRegister(requestCounter)
}

type ConfigStruct struct {
	Port string `yaml:"port"`
	SSL  struct {
		Enabled bool   `yaml:"enabled"`
		Cert    string `yaml:"cert"`
		Key     string `yaml:"key"`
	} `yaml:"ssl"`
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World!"))
	requestCounter.WithLabelValues("200").Inc()
	start := time.Now()
	requestCounter.WithLabelValues("time").Add(time.Since(start).Seconds())
}

func handlerHello(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello ," + name + "!"))
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "f", "", "путь к файлу конфигурации")
	flag.Parse()

	// Проверяем наличие пути к файлу конфигурации
	if len(configPath) == 0 {
		log.Fatalf("Не указан файл конфигурации.")
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
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Сервер запущен и слушает порт %s ...", config.Port)

	if config.SSL.Enabled {
		err = http.ListenAndServeTLS(":"+config.Port, config.SSL.Cert, config.SSL.Key, nil)
	} else {
		err = http.ListenAndServe(":"+config.Port, nil)
	}

	if err != nil {
		log.Fatalf("Ошибка сервера: %v\n", err)
		os.Exit(1)
	}

}
