package trace

import (
	"flag"
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	//FIXME: some use dashes, and some use underscores. What's the convention going forward?
	// kubeconfigPath is a string that gives the location of a valid kubeconfig file
	kubeconfigPath = flag.String("provisioner_k8s_kubeconfig", "", "Path to a valid kubeconfig file.")

	// configContext is a string that can be used to override the default context
	configContext = flag.String("provisioner_k8s_context", "", "The kubeconfig context to use, overrides the 'current-context' from the config")

	// configNamespace is a string that can be used to override the default namespace for objects
	configNamespace = flag.String("provisioner_k8s_namespace", "", "The kubernetes namespace to use for all objects. Default comes from the context or in-cluster config")


	// kubeconfigPath is a string that gives the location of a valid kubeconfig file
	 = flag.String("provisioner_k8s_", "", "Path to a valid kubeconfig file.")
)

func newDatadogTracer(serviceName string) (tracingService, io.Closer, error) {
	if *dataDogHost == "" || *dataDogPort == "" {
		return nil, nil, fmt.Errorf("need host and port to datadog agent to use datadog tracing")
	}

	t := opentracer.New(
		ddtracer.WithAgentAddr(*dataDogHost+":"+*dataDogPort),
		ddtracer.WithServiceName(serviceName),
		ddtracer.WithDebugMode(true),
		ddtracer.WithSampler(ddtracer.NewRateSampler(*samplingRate)),
	)

	opentracing.SetGlobalTracer(t)

	return openTracingService{Tracer: &datadogTracer{actual: t}}, &ddCloser{}, nil
}

var _ io.Closer = (*ddCloser)(nil)

type ddCloser struct{}

func (ddCloser) Close() error {
	ddtracer.Stop()
	return nil
}

func init() {
	tracingBackendFactories["opentracing-datadog"] = newDatadogTracer
}

var _ tracer = (*datadogTracer)(nil)

type datadogTracer struct {
	actual opentracing.Tracer
}

func (dt *datadogTracer) GetOpenTracingTracer() opentracing.Tracer {
	return dt.actual
}
