package token_bucket

import "github.com/prometheus/client_golang/prometheus"

var functionPodLabels = []string{"function_name"}

var funcAlives = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "fission_function_alive_num",
		Help: "A binary value indicating is the function_name, function_uid alive num",
	},
	functionPodLabels,
)

func init() {
	prometheus.MustRegister(funcAlives)
}

func SetFuncAliveNumInc(funcname string) {
	funcAlives.WithLabelValues(funcname).Inc()
}

func SetFuncAliveNumDec(funcname string) {
	funcAlives.WithLabelValues(funcname).Dec()
}
