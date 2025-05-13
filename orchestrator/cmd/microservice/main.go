package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rafaelmascaro/Weather-By-CEP-With-Tracing/orchestrator/configs"
	"github.com/rafaelmascaro/Weather-By-CEP-With-Tracing/orchestrator/internal/adapters/api"
	"github.com/rafaelmascaro/Weather-By-CEP-With-Tracing/orchestrator/internal/infra/web"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
)

func initProvider(serviceName, collectorURL string, zipkinURL string) (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	conn, err := grpc.NewClient(collectorURL, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	zipkinExporter, err := zipkin.New(zipkinURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create Zipkin exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(zipkinExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	shutdown, err := initProvider(
		configs.OtelServiceName,
		configs.OtelExporterOtlpEndpoint,
		configs.OtelExporterZipkinUrl,
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	tracer := otel.Tracer("microservice-tracer")

	locationClient := api.NewLocationClient(
		configs.LocationClientUrl,
		tracer,
		configs.LocationSpanNameOTEL,
	)

	weatherClient := api.NewWeatherClient(
		configs.WeatherClientUrl,
		configs.WeatherClientKey,
		tracer,
		configs.WeatherSpanNameOTEL,
	)

	server := web.NewServer(locationClient, weatherClient, tracer, configs.RequestNameOTEL)
	router := server.CreateServer()

	go func() {
		fmt.Println("Starting web server on port ", configs.WebServerPort)
		if err := http.ListenAndServe(configs.WebServerPort, router); err != nil {
			panic(err)
		}
	}()

	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed...")
	case <-ctx.Done():
		log.Println("Shutting down due to other reason...")
	}

	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
}
