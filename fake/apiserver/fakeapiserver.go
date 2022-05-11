package apiserver

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"

	"github.com/clusterpedia-io/client-go/fake/apiserver/registry/clusterpedia/resources"
	"github.com/clusterpedia-io/client-go/fake/kubeapiserver"
	"github.com/clusterpedia-io/client-go/fake/storage"
	"github.com/clusterpedia-io/client-go/fake/utils/filters"
	internal "github.com/clusterpedia-io/api/clusterpedia"
	"github.com/clusterpedia-io/api/clusterpedia/install"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/apiserver/pkg/authorization/authorizerfactory"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/healthz"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

var (
	// Scheme defines methods for serializing and deserializing API objects.
	Scheme = runtime.NewScheme()
	// Codecs provides methods for retrieving codecs and serializers for specific
	// versions and content types.
	Codecs = serializer.NewCodecFactory(Scheme)

	// ParameterCodec handles versioning of objects that are converted to query parameters.
	ParameterCodec = runtime.NewParameterCodec(Scheme)
)

func init() {
	install.Install(Scheme)

	// we need to add the options to empty v1
	// TODO fix the server code to avoid this
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})
	_ = metainternal.AddToScheme(Scheme)

	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	Scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)
}

func NewFakeApiserver(storageFactory storage.StorageFactory) (*httptest.Server, error) {
	genericConfig := genericapiserver.NewRecommendedConfig(Codecs)
	genericConfig.SecureServing = &genericapiserver.SecureServingInfo{Listener: fakeLocalhost443Listener{}}
	genericConfig.Authorization.Authorizer = authorizerfactory.NewAlwaysAllowAuthorizer()
	genericConfig.LoopbackClientConfig = &restclient.Config{
		ContentConfig: restclient.ContentConfig{NegotiatedSerializer: Codecs},
	}
	genericConfig.Version = &version.Info{
		Major: "1",
		Minor: "0",
	}
	completedConfig := genericConfig.Complete()

	// init apiGroupResources
	initialAPIGroupResources := []*restmapper.APIGroupResources{}
	err := json.Unmarshal([]byte(apiGroupResources), &initialAPIGroupResources)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshal api group resources json file. ")
	}

	resourceServerConfig := kubeapiserver.NewDefaultConfig()
	resourceServerConfig.GenericConfig.ExternalAddress = completedConfig.ExternalAddress
	resourceServerConfig.GenericConfig.LoopbackClientConfig = completedConfig.LoopbackClientConfig
	resourceServerConfig.ExtraConfig = kubeapiserver.ExtraConfig{
		StorageFactory:           storageFactory,
		InitialAPIGroupResources: initialAPIGroupResources,
	}
	kubeResourceAPIServer, err := resourceServerConfig.Complete().New(genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	handlerChainFunc := genericConfig.BuildHandlerChainFunc
	genericConfig.BuildHandlerChainFunc = func(apiHandler http.Handler, c *genericapiserver.Config) http.Handler {
		handler := handlerChainFunc(apiHandler, c)
		handler = filters.WithRequestQuery(handler)
		return handler
	}

	genericServer, err := completedConfig.New("clusterpedia", hooksDelegate{kubeResourceAPIServer})
	if err != nil {
		return nil, err
	}

	v1beta1storage := map[string]rest.Storage{}
	v1beta1storage["resources"] = resources.NewREST(kubeResourceAPIServer.Handler)

	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(internal.GroupName, Scheme, ParameterCodec, Codecs)
	apiGroupInfo.VersionedResourcesStorageMap["v1beta1"] = v1beta1storage
	if err := genericServer.InstallAPIGroup(&apiGroupInfo); err != nil {
		return nil, err
	}

	return httptest.NewServer(genericServer.Handler), nil
}

type fakeLocalhost443Listener struct{}

func (fakeLocalhost443Listener) Accept() (net.Conn, error) {
	return nil, nil
}

func (fakeLocalhost443Listener) Close() error {
	return nil
}

func (fakeLocalhost443Listener) Addr() net.Addr {
	return &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 443,
	}
}

type hooksDelegate struct {
	genericapiserver.DelegationTarget
}

func (s hooksDelegate) UnprotectedHandler() http.Handler {
	return nil
}

func (s hooksDelegate) HealthzChecks() []healthz.HealthChecker {
	return []healthz.HealthChecker{}
}

func (s hooksDelegate) ListedPaths() []string {
	return []string{}
}

func (s hooksDelegate) NextDelegate() genericapiserver.DelegationTarget {
	return nil
}
