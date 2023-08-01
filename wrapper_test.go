package monitoring

import (
	"testing"
)

func TestMetric_getName(t *testing.T) {
	type fields struct {
		Namespace   string
		Subsystem   string
		Name        string
		ConstLabels map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "case 1",
			fields: fields{},
			want:   "",
		},
		{
			name: "case 2",
			fields: fields{
				Namespace: "namespace",
			},
			want: "namespace",
		},
		{
			name: "case 3",
			fields: fields{
				Namespace: "namespace",
				Subsystem: "subsystem",
			},
			want: "namespace_subsystem",
		},
		{
			name: "case 4",
			fields: fields{
				Namespace: "namespace",
				Subsystem: "subsystem",
				Name:      "name",
			},
			want: "namespace_subsystem_name",
		},
		{
			name: "case 5",
			fields: fields{
				Namespace: "namespace",
				Subsystem: "subsystem",
				Name:      "name",
				ConstLabels: map[string]string{
					"label1": "value1",
				},
			},
			want: `namespace_subsystem_name{label1="value1"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metric{
				Namespace:   tt.fields.Namespace,
				Subsystem:   tt.fields.Subsystem,
				Name:        tt.fields.Name,
				ConstLabels: tt.fields.ConstLabels,
			}
			if got := m.getName(); got != tt.want {
				t.Errorf("getName() = %v, want %v", got, tt.want)
			}
		})
	}
}
