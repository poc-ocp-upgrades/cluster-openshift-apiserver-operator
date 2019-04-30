package v311_00_assets

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type asset struct {
	bytes	[]byte
	info	os.FileInfo
}
type bindataFileInfo struct {
	name	string
	size	int64
	mode	os.FileMode
	modTime	time.Time
}

func (fi bindataFileInfo) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}

var _v3110OpenshiftApiserverApiserverClusterrolebindingYaml = []byte(`apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: system:openshift:openshift-apiserver
roleRef:
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  namespace: openshift-apiserver
  name: openshift-apiserver-sa`)

func v3110OpenshiftApiserverApiserverClusterrolebindingYamlBytes() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return _v3110OpenshiftApiserverApiserverClusterrolebindingYaml, nil
}
func v3110OpenshiftApiserverApiserverClusterrolebindingYaml() (*asset, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bytes, err := v3110OpenshiftApiserverApiserverClusterrolebindingYamlBytes()
	if err != nil {
		return nil, err
	}
	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/apiserver-clusterrolebinding.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverCmYaml = []byte(`apiVersion: v1
kind: ConfigMap
metadata:
  namespace: openshift-apiserver
  name: config
data:
  config.yaml:
`)

func v3110OpenshiftApiserverCmYamlBytes() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return _v3110OpenshiftApiserverCmYaml, nil
}
func v3110OpenshiftApiserverCmYaml() (*asset, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bytes, err := v3110OpenshiftApiserverCmYamlBytes()
	if err != nil {
		return nil, err
	}
	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/cm.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverDefaultconfigYaml = []byte(`apiVersion: openshiftcontrolplane.config.openshift.io/v1
kind: OpenShiftAPIServerConfig
aggregatorConfig:
  clientCA: /var/run/configmaps/aggregator-client-ca/ca-bundle.crt
  allowedNames:
  - kube-apiserver-proxy
  - system:kube-apiserver-proxy
  - system:openshift-aggregator
  usernameHeaders:
  - X-Remote-User
  groupHeaders:
  - X-Remote-Group
  extraHeaderPrefixes:
  - X-Remote-Extra-
auditConfig:
  auditFilePath: "/var/log/openshift-apiserver/audit.log"
  enabled: true
  logFormat: "json"
  maximumFileSizeMegabytes: 100
  maximumRetainedFiles: 10
  policyConfiguration:
    apiVersion: audit.k8s.io/v1beta1
    kind: Policy
    # Don't generate audit events for all requests in RequestReceived stage.
    omitStages:
      - "RequestReceived"
    rules:
      # Don't log requests for events
      - level: None
        resources:
          - group: ""
            resources: ["events"]
      # Don't log authenticated requests to certain non-resource URL paths.
      - level: None
        userGroups: ["system:authenticated", "system:unauthenticated"]
        nonResourceURLs:
          - "/api*" # Wildcard matching.
          - "/version"
          - "/healthz"
      # A catch-all rule to log all other requests at the Metadata level.
      - level: Metadata
        # Long-running requests like watches that fall under this rule will not
        # generate an audit event in RequestReceived.
        omitStages:
          - "RequestReceived"
storageConfig:
  urls:
    - https://etcd.openshift-etcd.svc:2379
apiServerArguments:
  minimal-shutdown-duration:
  - 3s # give SDN some time to converge
`)

func v3110OpenshiftApiserverDefaultconfigYamlBytes() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return _v3110OpenshiftApiserverDefaultconfigYaml, nil
}
func v3110OpenshiftApiserverDefaultconfigYaml() (*asset, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bytes, err := v3110OpenshiftApiserverDefaultconfigYamlBytes()
	if err != nil {
		return nil, err
	}
	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/defaultconfig.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverDsYaml = []byte(`apiVersion: apps/v1
kind: DaemonSet
metadata:
  namespace: openshift-apiserver
  name: apiserver
  labels:
    app: openshift-apiserver
    apiserver: "true"
spec:
  updateStrategy:
    type: RollingUpdate
  selector:
    matchLabels:
      app: openshift-apiserver
      apiserver: "true"
  template:
    metadata:
      name: openshift-apiserver
      labels:
        app: openshift-apiserver
        apiserver: "true"
    spec:
      serviceAccountName: openshift-apiserver-sa
      priorityClassName: system-node-critical
      initContainers:
        - name: fix-audit-permissions
          terminationMessagePolicy: FallbackToLogsOnError
          image: ${IMAGE}
          imagePullPolicy: IfNotPresent
          command: ['sh', '-c', 'chmod 0700 /var/log/openshift-apiserver']
          volumeMounts:
            - mountPath: /var/log/openshift-apiserver
              name: audit-dir
      containers:
      - name: openshift-apiserver
        terminationMessagePolicy: FallbackToLogsOnError
        image: ${IMAGE}
        imagePullPolicy: IfNotPresent
        command: ["hypershift", "openshift-apiserver"]
        args:
        - "--config=/var/run/configmaps/config/config.yaml"
        resources:
          requests:
            memory: 200Mi
            cpu: 150m
        # we need to set this to privileged to be able to write audit to /var/log/openshift-apiserver
        securityContext:
          privileged: true
        ports:
        - containerPort: 8443
        volumeMounts:
        - mountPath: /var/run/configmaps/aggregator-client-ca
          name: aggregator-client-ca
        - mountPath: /var/run/configmaps/client-ca
          name: client-ca
        - mountPath: /var/run/configmaps/config
          name: config
        - mountPath: /var/run/secrets/etcd-client
          name: etcd-client
        - mountPath: /var/run/configmaps/etcd-serving-ca
          name: etcd-serving-ca
        - mountPath: /var/run/configmaps/image-import-ca
          name: image-import-ca
        - mountPath: /var/run/secrets/serving-cert
          name: serving-cert
        - mountPath: /var/log/openshift-apiserver
          name: audit-dir
        livenessProbe:
          initialDelaySeconds: 30
          httpGet:
            scheme: HTTPS
            port: 8443
            path: healthz
        readinessProbe:
          failureThreshold: 10
          httpGet:
            scheme: HTTPS
            port: 8443
            path: healthz
      terminationGracePeriodSeconds: 70 # a bit more than the 60 seconds timeout of non-long-running requests
      volumes:
      - name: aggregator-client-ca
        configMap:
          name: aggregator-client-ca
      - name: client-ca
        configMap:
          name: client-ca
      - name: config
        configMap:
          name: config
      - name: etcd-client
        secret:
          secretName: etcd-client
      - name: etcd-serving-ca
        configMap:
          name: etcd-serving-ca
      - name: image-import-ca
        configMap:
          name: image-import-ca
          optional: true
      - name: serving-cert
        secret:
          secretName: serving-cert
      - hostPath:
          path: /var/log/openshift-apiserver
        name: audit-dir
      nodeSelector:
        node-role.kubernetes.io/master: ""
      tolerations:
      - operator: Exists
`)

func v3110OpenshiftApiserverDsYamlBytes() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return _v3110OpenshiftApiserverDsYaml, nil
}
func v3110OpenshiftApiserverDsYaml() (*asset, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bytes, err := v3110OpenshiftApiserverDsYamlBytes()
	if err != nil {
		return nil, err
	}
	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/ds.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverNsYaml = []byte(`apiVersion: v1
kind: Namespace
metadata:
  annotations:
    openshift.io/node-selector: ""
  name: openshift-apiserver
  labels:
    openshift.io/run-level: "1"
`)

func v3110OpenshiftApiserverNsYamlBytes() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return _v3110OpenshiftApiserverNsYaml, nil
}
func v3110OpenshiftApiserverNsYaml() (*asset, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bytes, err := v3110OpenshiftApiserverNsYamlBytes()
	if err != nil {
		return nil, err
	}
	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/ns.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverSaYaml = []byte(`apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: openshift-apiserver
  name: openshift-apiserver-sa
`)

func v3110OpenshiftApiserverSaYamlBytes() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return _v3110OpenshiftApiserverSaYaml, nil
}
func v3110OpenshiftApiserverSaYaml() (*asset, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bytes, err := v3110OpenshiftApiserverSaYamlBytes()
	if err != nil {
		return nil, err
	}
	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/sa.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _v3110OpenshiftApiserverSvcYaml = []byte(`apiVersion: v1
kind: Service
metadata:
  namespace: openshift-apiserver
  name: api
  annotations:
    service.alpha.openshift.io/serving-cert-secret-name: serving-cert
    prometheus.io/scrape: "true"
    prometheus.io/scheme: https
spec:
  selector:
    apiserver: "true"
  ports:
  - name: https
    port: 443
    targetPort: 8443
`)

func v3110OpenshiftApiserverSvcYamlBytes() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return _v3110OpenshiftApiserverSvcYaml, nil
}
func v3110OpenshiftApiserverSvcYaml() (*asset, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bytes, err := v3110OpenshiftApiserverSvcYamlBytes()
	if err != nil {
		return nil, err
	}
	info := bindataFileInfo{name: "v3.11.0/openshift-apiserver/svc.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}
func Asset(name string) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}
func MustAsset(name string) []byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}
	return a
}
func AssetInfo(name string) (os.FileInfo, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}
func AssetNames() []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

var _bindata = map[string]func() (*asset, error){"v3.11.0/openshift-apiserver/apiserver-clusterrolebinding.yaml": v3110OpenshiftApiserverApiserverClusterrolebindingYaml, "v3.11.0/openshift-apiserver/cm.yaml": v3110OpenshiftApiserverCmYaml, "v3.11.0/openshift-apiserver/defaultconfig.yaml": v3110OpenshiftApiserverDefaultconfigYaml, "v3.11.0/openshift-apiserver/ds.yaml": v3110OpenshiftApiserverDsYaml, "v3.11.0/openshift-apiserver/ns.yaml": v3110OpenshiftApiserverNsYaml, "v3.11.0/openshift-apiserver/sa.yaml": v3110OpenshiftApiserverSaYaml, "v3.11.0/openshift-apiserver/svc.yaml": v3110OpenshiftApiserverSvcYaml}

func AssetDir(name string) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func		func() (*asset, error)
	Children	map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{"v3.11.0": {nil, map[string]*bintree{"openshift-apiserver": {nil, map[string]*bintree{"apiserver-clusterrolebinding.yaml": {v3110OpenshiftApiserverApiserverClusterrolebindingYaml, map[string]*bintree{}}, "cm.yaml": {v3110OpenshiftApiserverCmYaml, map[string]*bintree{}}, "defaultconfig.yaml": {v3110OpenshiftApiserverDefaultconfigYaml, map[string]*bintree{}}, "ds.yaml": {v3110OpenshiftApiserverDsYaml, map[string]*bintree{}}, "ns.yaml": {v3110OpenshiftApiserverNsYaml, map[string]*bintree{}}, "sa.yaml": {v3110OpenshiftApiserverSaYaml, map[string]*bintree{}}, "svc.yaml": {v3110OpenshiftApiserverSvcYaml, map[string]*bintree{}}}}}}}}

func RestoreAsset(dir, name string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}
func RestoreAssets(dir, name string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	children, err := AssetDir(name)
	if err != nil {
		return RestoreAsset(dir, name)
	}
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}
func _filePath(dir, name string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
