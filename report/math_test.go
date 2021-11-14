package report

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func Test_latencies(t *testing.T) {
	type args struct {
		latencies []float64
	}
	tests := []struct {
		name string
		args args
		want []LatencyDistribution
	}{
		// TODO: Add test cases.
		{
			name: "case1",

			args: args{
				latencies: []float64{1, 1, 1, 2, 3, 4, 5, 7, 7, 7}, //[]float64
			},
			want: []LatencyDistribution{{10, 1 * time.Second}, {25, 1 * time.Second}, {50, 3 * time.Second}, {75, 7 * time.Second}, {90, 7 * time.Second}, {95, 7 * time.Second}, {99, 7 * time.Second}}, //[]LatencyDistribution,
		},
	}
	for _, tt := range tests {
		tt := tt
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		t.Run(tt.name, func(t *testing.T) {
			if got := latencies(tt.args.latencies); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("latencies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_histogram(t *testing.T) {
	type args struct {
		latencies []float64
		slowest   float64
		fastest   float64
	}
	tests := []struct {
		name string
		args args
		want []Bucket
	}{
		// TODO: Add test cases.
		{
			name: "case1",

			args: args{
				latencies: []float64{1, 1, 1, 2, 3, 4, 5, 7, 7, 7}, //[]float64
				slowest:   7.0,                                     //float64
				fastest:   1.0,                                     //float64
			},
			want: []Bucket{{7, 10, 1}, {6.4, 0, 0}, {5.8, 0, 0}, {5.2, 0, 0}, {4.6, 0, 0}, {4, 0, 0}, {3.4000000000000004, 0, 0}, {2.8000000000000003, 0, 0}, {2.2, 0, 0}, {1.6, 0, 0}, {1, 0, 0}}, //[]Bucket,
		},
	}
	for _, tt := range tests {
		tt := tt
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		t.Run(tt.name, func(t *testing.T) {
			if got := histogram(tt.args.latencies, tt.args.slowest, tt.args.fastest); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("histogram() = %v, want %v", got, tt.want)
			}
		})
	}
}
