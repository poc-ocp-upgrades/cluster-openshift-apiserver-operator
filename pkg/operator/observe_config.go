package operator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/golang/glog"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/diff"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/client-go/util/workqueue"

	operatorconfigclientv1alpha1 "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/clientset/versioned/typed/openshiftapiserver/v1alpha1"
	operatorconfiginformerv1alpha1 "github.com/openshift/cluster-openshift-apiserver-operator/pkg/generated/informers/externalversions/openshiftapiserver/v1alpha1"
)

type observeConfigFunc func(kubernetes.Interface, *rest.Config, map[string]interface{}) (map[string]interface{}, error)

type ConfigObserver struct {
	operatorConfigClient operatorconfigclientv1alpha1.OpenshiftapiserverV1alpha1Interface

	kubeClient   kubernetes.Interface
	clientConfig *rest.Config

	// queue only ever has one item, but it has nice error handling backoff/retry semantics
	queue workqueue.RateLimitingInterface

	rateLimiter flowcontrol.RateLimiter
	observers   []observeConfigFunc
}

func NewConfigObserver(
	operatorConfigInformer operatorconfiginformerv1alpha1.OpenShiftAPIServerOperatorConfigInformer,
	kubeInformersForKubeApiserverNamespace kubeinformers.SharedInformerFactory,
	kubeInformersForKubeSystemNamespace kubeinformers.SharedInformerFactory,
	operatorConfigClient operatorconfigclientv1alpha1.OpenshiftapiserverV1alpha1Interface,
	kubeClient kubernetes.Interface,
	clientConfig *rest.Config,
) *ConfigObserver {
	c := &ConfigObserver{
		operatorConfigClient: operatorConfigClient,
		kubeClient:           kubeClient,
		clientConfig:         clientConfig,

		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ConfigObserver"),

		rateLimiter: flowcontrol.NewTokenBucketRateLimiter(0.05 /*3 per minute*/, 4),
		observers: []observeConfigFunc{
			observeKubeAPIServerPublicInfo,
			observeEtcdEndpoints,
		},
	}

	operatorConfigInformer.Informer().AddEventHandler(c.eventHandler())
	kubeInformersForKubeApiserverNamespace.Core().V1().ConfigMaps().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForKubeApiserverNamespace.Core().V1().ServiceAccounts().Informer().AddEventHandler(c.eventHandler())
	kubeInformersForKubeSystemNamespace.Core().V1().ConfigMaps().Informer().AddEventHandler(c.eventHandler())

	return c
}

func observeKubeAPIServerPublicInfo(kubeClient kubernetes.Interface, config *rest.Config, observedConfig map[string]interface{}) (map[string]interface{}, error) {
	kubeAPIServerPublicInfo, err := kubeClient.CoreV1().ConfigMaps(kubeAPIServerNamespaceName).Get("public-info", metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if kubeAPIServerPublicInfo != nil {
		if val, ok := kubeAPIServerPublicInfo.Data["projectConfig.defaultNodeSelector"]; ok {
			unstructured.SetNestedField(observedConfig, val, "projectConfig", "defaultNodeSelector")
		}
	}

	return observedConfig, nil
}

// observeEtcdEndpoints reads the etcd endpoints from the endpoints object and then manually pull out the hostnames to
// get the etcd urls for our config. Setting them observed config causes the normal reconciliation loop to run
func observeEtcdEndpoints(kubeClient kubernetes.Interface, clientConfig *rest.Config, observedConfig map[string]interface{}) (map[string]interface{}, error) {
	etcdURLs := []string{}
	etcdEndpoints, err := kubeClient.CoreV1().Endpoints(etcdNamespaceName).Get("etcd", metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if etcdEndpoints != nil {
		for _, subset := range etcdEndpoints.Subsets {
			for _, address := range subset.Addresses {
				etcdURLs = append(etcdURLs, "https://"+address.Hostname+"."+etcdEndpoints.Annotations["alpha.installer.openshift.io/dns-suffix"]+":2379")
			}
		}
	}
	if len(etcdURLs) > 0 {
		unstructured.SetNestedStringSlice(observedConfig, etcdURLs, "storageConfig", "urls")
	}

	return observedConfig, nil
}

// sync reacts to a change in prereqs by finding information that is required to match another value in the cluster. This
// must be information that is logically "owned" by another component.
func (c ConfigObserver) sync() error {
	var err error
	observedConfig := map[string]interface{}{}

	for _, observer := range c.observers {
		observedConfig, err = observer(c.kubeClient, c.clientConfig, observedConfig)
		if err != nil {
			return err
		}
	}

	operatorConfig, err := c.operatorConfigClient.OpenShiftAPIServerOperatorConfigs().Get("instance", metav1.GetOptions{})
	if err != nil {
		return err
	}

	// don't worry about errors
	currentConfig := map[string]interface{}{}
	json.NewDecoder(bytes.NewBuffer(operatorConfig.Spec.ObservedConfig.Raw)).Decode(&currentConfig)
	if reflect.DeepEqual(currentConfig, observedConfig) {
		return nil
	}

	glog.Info("writing updated observedConfig: %v", diff.ObjectDiff(operatorConfig.Spec.ObservedConfig.Object, observedConfig))
	operatorConfig.Spec.ObservedConfig = runtime.RawExtension{Object: &unstructured.Unstructured{Object: observedConfig}}
	if _, err := c.operatorConfigClient.OpenShiftAPIServerOperatorConfigs().Update(operatorConfig); err != nil {
		return err
	}

	return nil
}

func (c *ConfigObserver) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	glog.Infof("Starting ConfigObserver")
	defer glog.Infof("Shutting down ConfigObserver")

	// doesn't matter what workers say, only start one.
	go wait.Until(c.runWorker, time.Second, stopCh)

	<-stopCh
}

func (c *ConfigObserver) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *ConfigObserver) processNextWorkItem() bool {
	dsKey, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(dsKey)

	// before we call sync, we want to wait for token.  We do this to avoid hot looping.
	c.rateLimiter.Accept()

	err := c.sync()
	if err == nil {
		c.queue.Forget(dsKey)
		return true
	}

	utilruntime.HandleError(fmt.Errorf("%v failed with : %v", dsKey, err))
	c.queue.AddRateLimited(dsKey)

	return true
}

// eventHandler queues the operator to check spec and status
func (c *ConfigObserver) eventHandler() cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { c.queue.Add(workQueueKey) },
		UpdateFunc: func(old, new interface{}) { c.queue.Add(workQueueKey) },
		DeleteFunc: func(obj interface{}) { c.queue.Add(workQueueKey) },
	}
}

func (c *ConfigObserver) namespaceEventHandler() cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ns, ok := obj.(*corev1.Namespace)
			if !ok {
				c.queue.Add(workQueueKey)
			}
			if ns.Name == targetNamespaceName {
				c.queue.Add(workQueueKey)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			ns, ok := old.(*corev1.Namespace)
			if !ok {
				c.queue.Add(workQueueKey)
			}
			if ns.Name == targetNamespaceName {
				c.queue.Add(workQueueKey)
			}
		},
		DeleteFunc: func(obj interface{}) {
			ns, ok := obj.(*corev1.Namespace)
			if !ok {
				tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
					return
				}
				ns, ok = tombstone.Obj.(*corev1.Namespace)
				if !ok {
					utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a Namespace %#v", obj))
					return
				}
			}
			if ns.Name == targetNamespaceName {
				c.queue.Add(workQueueKey)
			}
		},
	}
}
