package main

import (
	"fmt"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func main() {

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	restConfig, _ := clientcmd.BuildConfigFromFlags("", "/Users/acejilam/.kube/172.20.53.21.config")
	restConfig.TLSClientConfig.Insecure = true
	restConfig.TLSClientConfig.CAData = nil
	restConfig.TLSClientConfig.CAFile = ""
	if restConfig != nil && rest.IsConfigTransportTLS(*restConfig) {
		transport, _ := rest.TransportFor(restConfig)
		client.Transport = transport
	}

	path := "/configz"
	configzURL := url.URL{
		Scheme: "https",
		Host:   net.JoinHostPort("172.20.53.21", strconv.Itoa(10250)),
		Path:   path,
	}
	rsp, _ := client.Get(configzURL.String())
	fmt.Println(rsp)
}
