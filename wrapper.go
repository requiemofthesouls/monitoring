package monitoring

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

func New() Wrapper {
	return &wrapper{
		collector: metrics.NewSet(),
	}
}

type (
	Wrapper interface {
		Counter(*Metric, int)
		CounterFloat(*Metric, float64)
		Inc(*Metric)
		Dec(*Metric)
		ExecutionTime(*Metric, func() error) error
		Val(*Metric, uint64)
		MetricsHandler() http.Handler
	}

	Metric struct {
		Namespace   string
		Subsystem   string
		Name        string
		ConstLabels map[string]string
	}

	wrapper struct {
		collector *metrics.Set
	}
)

func (w *wrapper) MetricsHandler() http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		w.collector.WritePrometheus(resp)
		metrics.WritePrometheus(resp, true)
	})
}

func (w *wrapper) Counter(metric *Metric, n int) {
	counter := w.collector.GetOrCreateCounter(metric.getName())
	counter.Add(n)
}

func (w *wrapper) CounterFloat(metric *Metric, f float64) {
	counter := w.collector.GetOrCreateFloatCounter(metric.getName())
	counter.Add(f)
}

func (w *wrapper) Inc(metric *Metric) {
	counter := w.collector.GetOrCreateCounter(metric.getName())
	counter.Inc()
}

func (w *wrapper) Dec(metric *Metric) {
	counter := w.collector.GetOrCreateCounter(metric.getName())
	counter.Dec()
}

func (w *wrapper) ExecutionTime(metric *Metric, f func() error) error {
	startTime := time.Now()
	err := f()
	w.collector.GetOrCreateSummary(metric.getName()).UpdateDuration(startTime)

	return err
}

func (w *wrapper) Val(metric *Metric, n uint64) {
	counter := w.collector.GetOrCreateCounter(metric.getName())
	counter.Set(n)
}

func (m *Metric) getName() string {
	strs := make([]string, 0, 3)
	if m.Namespace != "" {
		strs = append(strs, m.Namespace)
	}

	if m.Subsystem != "" {
		strs = append(strs, m.Subsystem)
	}

	if m.Name != "" {
		strs = append(strs, m.Name)
	}

	name := strings.Join(strs, "_")
	labels := marshalLabels(m.ConstLabels)

	metricWithLabels := fmt.Sprintf(`%s%s`, name, labels)

	return metricWithLabels
}

func marshalLabels(rawLabels map[string]string) string {
	if len(rawLabels) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("{")
	i := 0
	for label, value := range rawLabels {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(label)
		sb.WriteString(`=`)
		sb.WriteString(`"`)
		sb.WriteString(value)
		sb.WriteString(`"`)
		i++
	}
	sb.WriteString("}")
	return sb.String()
}
