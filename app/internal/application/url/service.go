package urlapp

import (
	"github.com/sanctumlabs/curtz/app/internal/domain/url"
	"github.com/sanctumlabs/curtz/app/internal/ports"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Service exposes the Url use cases.
type Service struct {
	url                url.UrlDatastore
	shortCodeGenerator ports.ShortCodeGenerator
	tracer             trace.Tracer
	logPrefix          string
}

// NewService wires the Url use cases to their ports.
func NewService(url url.UrlDatastore, shortCodeGenerator ports.ShortCodeGenerator) *Service {
	return &Service{
		url:                url,
		shortCodeGenerator: shortCodeGenerator,
		tracer:             otel.Tracer(instrumentationName),
		logPrefix:          "UrlService",
	}
}
