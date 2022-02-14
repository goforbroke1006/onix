package cmd

import (
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/goforbroke1006/onix/internal/component/stub/prometheus"
)

func NewStubPrometheusCmd() *cobra.Command {
	return &cobra.Command{
		Use: "prometheus",
		Run: func(cmd *cobra.Command, args []string) {
			httpAddr := viper.GetString("server.http.stub.prometheus")

			// stub for https://prometheus.io/docs/prometheus/latest/querying/api/#instant-queries
			http.HandleFunc("/api/v1/query", func(w http.ResponseWriter, req *http.Request) {

			})

			// stub for https://prometheus.io/docs/prometheus/latest/querying/api/#range-queries
			http.HandleFunc("/api/v1/query_range", func(w http.ResponseWriter, req *http.Request) {
				var (
					query   = req.URL.Query()["query"][0]
					start   = req.URL.Query()["start"][0]
					end     = req.URL.Query()["end"][0]
					step    = req.URL.Query()["step"][0]
					timeout = req.URL.Query()["timeout"][0]
				)
				_, _, _, _, _ = query, start, end, step, timeout

				if start == "1643894976" /*&& end == "1643981340"*/ { // release 1.19.0
					if strings.Contains(query, "api_requests_processing_duration_seconds_bucket") &&
						strings.Contains(query, "ONE") {
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte(prometheus.ResponseQueryRangeRelease1_19_0_20220208_20220209_by15m_DO))
						return
					}
					if strings.Contains(query, "api_requests_processing_duration_seconds_bucket") &&
						strings.Contains(query, "TWO") {
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte(prometheus.ResponseQueryRangeRelease1_19_0_20220208_20220209_by15m_FX))
						return
					}
				}

				if start == "1642877700" /*&& end == "1643981340"*/ { // release 1.17.0
					if strings.Contains(query, "api_requests_processing_duration_seconds_bucket") &&
						strings.Contains(query, "ONE") {
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte(prometheus.ResponseQueryRangeRelease1_17_0_20220122_20220123_by15m_DO))
						return
					}
					if strings.Contains(query, "api_requests_processing_duration_seconds_bucket") &&
						strings.Contains(query, "TWO") {
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte(prometheus.ResponseQueryRangeRelease1_17_0_20220122_20220123_by15m_FX))
						return
					}
				}

				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(prometheus.ResponseFailed))
				return

			})
			log.Fatal(http.ListenAndServe(httpAddr, nil))
		},
	}
}
