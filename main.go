/*
Copyright 2022 The Koordinator Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"github.com/koordinator-sh/koordinator/cmd/koordlet/options"
	"github.com/koordinator-sh/koordinator/pkg/features"
	agent "github.com/koordinator-sh/koordinator/pkg/koordlet"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_audit"
	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_config"
	metricsutil "github.com/koordinator-sh/koordinator/pkg/util/metrics"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/component-base/logs"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"time"

	"github.com/koordinator-sh/koordinator/pkg/koordlet/over_metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog/v2"
)

func main() {

	os.Args = append(os.Args,
		"--runtime-hooks-addr=/etc/runtime/hookserver.d/koordlet.sock", "--v=9",
	)

	cfg := over_config.NewConfiguration()
	cfg.InitFlags(flag.CommandLine)
	logs.AddGoFlags(flag.CommandLine)
	flag.Parse()

	go wait.Forever(klog.Flush, 5*time.Second)
	defer klog.Flush()

	if *options.EnablePprof {
		go func() {
			klog.V(4).Infof("Starting pprof on %v", *options.PprofAddr)
			if err := http.ListenAndServe(*options.PprofAddr, nil); err != nil {
				klog.Errorf("Unable to start pprof on %v, error: %v", *options.PprofAddr, err)
			}
		}()
	}

	if err := features.DefaultMutableKoordletFeatureGate.SetFromMap(cfg.FeatureGates); err != nil {
		klog.Fatalf("Unable to setup feature-gates: %v", err)
	}

	stopCtx := signals.SetupSignalHandler()

	// setup the default auditor
	if features.DefaultKoordletFeatureGate.Enabled(features.AuditEvents) {
		over_audit.SetupDefaultAuditor(cfg.AuditConf, stopCtx.Done())
	}

	// Get a config to talk to the apiserver
	klog.Info("Setting up kubeconfig for koordlet")
	err := cfg.InitKubeConfigForKoordlet(*options.KubeAPIQPS, *options.KubeAPIBurst)
	if err != nil {
		klog.Fatalf("Unable to setup kubeconfig: %v", err)
	}

	d, err := agent.NewDaemon(cfg)
	if err != nil {
		klog.Fatalf("Unable to setup koordlet daemon: %v", err)
	}

	// Expose the Prometheus http endpoint
	go installHTTPHandler()

	// Start the Cmd
	klog.Info("Starting the koordlet daemon")
	d.Run(stopCtx.Done())
}

func installHTTPHandler() {
	klog.Infof("Starting prometheus server on %v", *options.ServerAddr)
	mux := http.NewServeMux()
	mux.Handle(over_metrics.ExternalHTTPPath, promhttp.HandlerFor(over_metrics.ExternalRegistry, promhttp.HandlerOpts{}))
	mux.Handle(over_metrics.InternalHTTPPath, promhttp.HandlerFor(over_metrics.InternalRegistry, promhttp.HandlerOpts{}))
	// merge internal and external
	mux.Handle(over_metrics.DefaultHTTPPath, promhttp.HandlerFor(
		metricsutil.MergedGatherFunc(over_metrics.InternalRegistry, over_metrics.ExternalRegistry), promhttp.HandlerOpts{}))
	if features.DefaultKoordletFeatureGate.Enabled(features.AuditEventsHTTPHandler) {
		mux.HandleFunc("/events", over_audit.HttpHandler())
	}
	// http.HandleFunc("/healthz", d.HealthzHandler())
	klog.Fatalf("Prometheus monitoring failed: %v", http.ListenAndServe(*options.ServerAddr, mux))
}
