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
            Name: "http_requests_total",
            Help: "Total number of HTTP requests.",
        },
        []string{"endpoint"},
    )

var  requestLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request latency in seconds.",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}, // Группы задержки
	},
	[]string{"endpoint"},
)


func init() {
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestLatency)
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


func LoggerMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
        next.ServeHTTP(w, r)

		duration := float64(time.Since(start).Seconds())
        requestLatency.WithLabelValues(r.URL.Path).Observe(duration)
        requestCounter.WithLabelValues(r.URL.Path).Inc()
		requestCounter.WithLabelValues(r.URL.Path).Inc()

		log.Printf("%s %s %v\n", r.Method, r.URL.Path, duration)
    })
}

func handlerHome(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World!"))
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

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlerHome)
	mux.HandleFunc("/hello", handlerHello)
	mux.Handle("/metrics", promhttp.Handler())

	handler := LoggerMiddleware(mux)

    server := &http.Server{
        Addr:           ":"+config.Port,
        Handler:        handler,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }

	log.Printf("Сервер запущен и слушает порт %s ...", config.Port)

	if config.SSL.Enabled {
		err = server.ListenAndServeTLS(config.SSL.Cert, config.SSL.Key)
	} else {
		err = server.ListenAndServe()
	}

	if err != nil {
		log.Fatalf("Ошибка сервера: %v\n", err)
		os.Exit(1)
	}

}
