// Command otelcol-static-observer-demo runs a minimal OpenTelemetry Collector
// that demonstrates the static_observer extension with receiver_creator.
//
// The static_observer fires a single synthetic endpoint of type "static" on
// startup. The receiver_creator matches it and instantiates two hostmetrics
// receiver instances — each with distinct resource_attributes (service.name,
// deployment.environment) — without any changes to the hostmetrics receiver.
//
// Both pipelines fan into the debug exporter (verbosity: detailed) so you can
// see the per-instance resource attributes printed to stdout.
package main

import (
	"log"
	"os"

	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/staticobserver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/receivercreator"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	envprovider "go.opentelemetry.io/collector/confmap/provider/envprovider"
	fileprovider "go.opentelemetry.io/collector/confmap/provider/fileprovider"
	yamlprovider "go.opentelemetry.io/collector/confmap/provider/yamlprovider"
	"go.opentelemetry.io/collector/exporter/debugexporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/receiver"
	otelconftelemetry "go.opentelemetry.io/collector/service/telemetry/otelconftelemetry"
)

func main() {
	info := component.BuildInfo{
		Command:     "otelcol-static-observer-demo",
		Description: "Demo collector showing static_observer + receiver_creator feature",
		Version:     "0.0.1",
	}

	set := otelcol.CollectorSettings{
		BuildInfo: info,
		Factories: components,
		ConfigProviderSettings: otelcol.ConfigProviderSettings{
			ResolverSettings: confmap.ResolverSettings{
				ProviderFactories: []confmap.ProviderFactory{
					envprovider.NewFactory(),
					fileprovider.NewFactory(),
					yamlprovider.NewFactory(),
				},
				DefaultScheme: "file",
			},
		},
	}

	cmd := otelcol.NewCommand(set)
	if err := cmd.Execute(); err != nil {
		log.Printf("collector error: %v", err)
		os.Exit(1)
	}
}

func components() (otelcol.Factories, error) {
	var f otelcol.Factories
	var err error

	f.Extensions, err = otelcol.MakeFactoryMap[extension.Factory](
		staticobserver.NewFactory(),
	)
	if err != nil {
		return f, err
	}

	f.Receivers, err = otelcol.MakeFactoryMap[receiver.Factory](
		receivercreator.NewFactory(),
		hostmetricsreceiver.NewFactory(),
	)
	if err != nil {
		return f, err
	}

	f.Exporters, err = otelcol.MakeFactoryMap(
		debugexporter.NewFactory(),
	)
	if err != nil {
		return f, err
	}

	f.Telemetry = otelconftelemetry.NewFactory()
	return f, nil
}
