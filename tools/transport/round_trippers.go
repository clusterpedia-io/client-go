/*
Copyright 2021 clusterpedia Authors

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

package transport

import (
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/clusterpedia-io/client-go/constants"
)

var (
	// regex matches "/apis/clusterpedia.io/{version}/{path}"
	clusterpediaAPIsRegex = regexp.MustCompile(`^(/apis/clusterpedia\.io/v\w*.)/(.*)`)
)

// ClusternetTransport is a transport to redirect requests to clusternet-hub
type ClusterpediaTransport struct {
	// relative paths may omit leading slash
	path    string
	cluster string

	rt http.RoundTripper
}

func (t *ClusterpediaTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.normalizeLocation(req.URL)
	return t.rt.RoundTrip(req)
}

// normalizeLocation format the request URL to Clusternet shadow GVKs
func (t *ClusterpediaTransport) normalizeLocation(location *url.URL) {
	curPath := location.Path
	// Trim returns a slice of the string s with all leading and trailing Unicode code points contained in cutset removed.
	// so we use Replace here
	reqPath := strings.Replace(curPath, t.path, "", 1)

	// we don't normalize request for Group clusterpedia.io
	if clusterpediaAPIsRegex.MatchString(reqPath) {
		return
	}

	paths := []string{constants.ClusterPediaAPIPath}
	if len(t.cluster) > 0 {
		paths = append(paths, "clusters", t.cluster)
	}
	location.Path = path.Join(append(paths, reqPath)...)
}

func NewTransportForCluster(host, cluster string, rt http.RoundTripper) *ClusterpediaTransport {
	// host must be a host string, a host:port pair, or a URL to the base of the apiserver.
	// If a URL is given then the (optional) Path of that URL represents a prefix that must
	// be appended to all request URIs used to access the apiserver. This allows a frontend
	// proxy to easily relocate all of the apiserver endpoints.
	return &ClusterpediaTransport{
		path:    urlMustParse(host).Path,
		cluster: cluster,
		rt:      rt,
	}
}

func NewTransport(host string, rt http.RoundTripper) *ClusterpediaTransport {
	return NewTransportForCluster(host, "", rt)
}

func urlMustParse(path string) *url.URL {
	location, err := url.Parse(strings.TrimRight(path, "/"))
	if err != nil {
		panic(err)
	}
	return location
}
