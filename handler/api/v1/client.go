package ccp
  
import (
	//"reflect"
	//"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"bytes"
        //"encoding/binary"
        "encoding/json"
        "fmt"
        "k8s.io/api/core/v1"
	//apiv1 "k8s.io/client-go/pkg/api/v1"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        //appv1 "k8s.io/api/apps/v1"
        "k8s.io/apimachinery/pkg/fields"
        "k8s.io/client-go/kubernetes"
        restclient "k8s.io/client-go/rest"
        "k8s.io/client-go/tools/cache"
        "k8s.io/client-go/tools/clientcmd"
        //"k8s.io/client-go/util/workqueue"
        "k8s.io/klog"
	"net/http"
        //"k8s.io/apimachinery/pkg/types"
        "os"
        //"time"
        "path/filepath"
        //"strconv"
        //"strings"
)

var (
        kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
        config     *restclient.Config
        err        error
        podcount = 0
	//endpointAPI = "http://10.221.42.109:30132/sdc/nitro/v1/config/app_event"
	endpointAPI = "http://localhost:8080/api/v1/pods"
)

//This is interface for Kubernetes API Server
type KubernetesAPIServer struct {
        Suffix string
        Client kubernetes.Interface
}

func Generate () string {
	uuid := uuid.New()
	s := uuid.String()
	return s
}

func CreateK8sApiserverClient() (*KubernetesAPIServer, error) {
        klog.Info("[INFO] Creating API Client")
        api := &KubernetesAPIServer{}
        config, err = clientcmd.BuildConfigFromFlags("", "")
        if err != nil {
                klog.Error("[WARNING] Citrix Node Controller Runs outside cluster")
                config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
                if err != nil {
                        klog.Error("[ERROR] Did not find valid kube config info")
                        klog.Fatal(err)
                }
        }

        client, err := kubernetes.NewForConfig(config)
        if err != nil {
                klog.Error("[ERROR] Failed to establish connection")
                klog.Fatal(err)
        }
        klog.Info("[INFO] Kubernetes Client is created")
        api.Client = client
        return api, nil
}
func ingressEventParseAndSendData(obj interface{}, eventType string) {
	objByte, err := json.Marshal(obj)
        if err != nil {
        	klog.Errorf("[ERROR] Failed to Marshal original object: %v", err)
        }
        var objJson v1beta1.Ingress
	if err = json.Unmarshal(objByte, &objJson); err != nil {
                klog.Errorf("[ERROR] Failed to unmarshal original object: %v", err)
        }
	message, err := json.MarshalIndent(objJson, "", "  ")
	parseAndSendData(string (message), objJson.ObjectMeta,  obj.(*v1beta1.Ingress).TypeMeta, "Ingress", eventType)

}
func IngressWatcher(api *KubernetesAPIServer) {
        ingressListWatcher := cache.NewListWatchFromClient(api.Client.ExtensionsV1beta1().RESTClient(), "ingresses", v1.NamespaceAll, fields.Everything())
        _, controller := cache.NewInformer(ingressListWatcher, &v1beta1.Ingress{}, 0, cache.ResourceEventHandlerFuncs{
                AddFunc: func(obj interface{}) {
			ingressEventParseAndSendData(obj, "ADDED")
                },
                UpdateFunc: func(obj interface{}, newobj interface{}) {
			ingressEventParseAndSendData(obj, "MODIFIED")
                },
                DeleteFunc: func(obj interface{}) {
			ingressEventParseAndSendData(obj, "DELETED")
                },
        },
        )
        stop := make(chan struct{})
        go controller.Run(stop)
        return
}
func endpointEventParseAndSendData(obj interface{}, eventType string) {
	objByte, err := json.Marshal(obj)
        if err != nil {
        	klog.Errorf("[ERROR] Failed to Marshal original object: %v", err)
        }
        var objJson v1.Endpoints
	if err = json.Unmarshal(objByte, &objJson); err != nil {
        	klog.Errorf("[ERROR] Failed to unmarshal original object: %v", err)
        }
	message, err := json.MarshalIndent(objJson, "", "  ")
	parseAndSendData(string (message), objJson.ObjectMeta, objJson.TypeMeta, "Endpoints", eventType)
}
func EndpointWatcher(api *KubernetesAPIServer) {
        endpointListWatcher := cache.NewListWatchFromClient(api.Client.Core().RESTClient(), "endpoints", v1.NamespaceAll, fields.Everything())
        _, controller := cache.NewInformer(endpointListWatcher, &v1.Endpoints{}, 0, cache.ResourceEventHandlerFuncs{
                AddFunc: func(obj interface{}) {
			endpointEventParseAndSendData(obj, "ADDED") 
                },
                UpdateFunc: func(obj interface{}, newobj interface{}) {
			endpointEventParseAndSendData(obj, "MODIFIED")
                },
                DeleteFunc: func(obj interface{}) {
			endpointEventParseAndSendData(obj, "DELETED")
                },
        },
        )
        stop := make(chan struct{})
        go controller.Run(stop)
        return
}


func serviceEventParseAndSendData(obj interface{}, eventType string) {
	objByte, err := json.Marshal(obj)
        if err != nil {
        	klog.Errorf("[ERROR] Failed to Marshal original object: %v", err)
        }
        var objJson v1.Service
	if err = json.Unmarshal(objByte, &objJson); err != nil {
        	klog.Errorf("[ERROR] Failed to unmarshal original object: %v", err)
        }
	message, err := json.MarshalIndent(objJson, "", "  ")
	parseAndSendData(string (message), objJson.ObjectMeta, objJson.TypeMeta, "Service", eventType)
}

func ServiceWatcher(api *KubernetesAPIServer) {
        serviceListWatcher := cache.NewListWatchFromClient(api.Client.Core().RESTClient(), string(v1.ResourceServices), v1.NamespaceAll, fields.Everything())
        _, controller := cache.NewInformer(serviceListWatcher, &v1.Service{}, 0, cache.ResourceEventHandlerFuncs{
                AddFunc: func(obj interface{}) {
			serviceEventParseAndSendData(obj, "ADDED")
                },
                UpdateFunc: func(obj interface{}, newobj interface{}) {
			serviceEventParseAndSendData(obj, "MODIFIED")
                },
                DeleteFunc: func(obj interface{}) {
			serviceEventParseAndSendData(obj, "DELETED")
                },
        },
        )
        stop := make(chan struct{})
        go controller.Run(stop)
        return
}


func podEventParseAndSendData(obj interface{}, eventType string){
	objByte, err := json.Marshal(obj)
        if err != nil {
        	klog.Errorf("[ERROR] Failed to Marshal original object: %v", err)
        }
        var objJson v1.Pod
	if err = json.Unmarshal(objByte, &objJson); err != nil {
                klog.Errorf("[ERROR] Failed to unmarshal original object: %v", err)
        }
	message, err := json.MarshalIndent(objJson, "", "  ")
	parseAndSendData(string (message), objJson.ObjectMeta, objJson.TypeMeta, "Pod", eventType)
}

func PodWatcher(api *KubernetesAPIServer) {
        PodListWatcher := cache.NewListWatchFromClient(api.Client.Core().RESTClient(), string(v1.ResourcePods), v1.NamespaceAll, fields.Everything())
        _, controller := cache.NewInformer(PodListWatcher, &v1.Pod{}, 0, cache.ResourceEventHandlerFuncs{
                AddFunc: func(obj interface{}) {
			podEventParseAndSendData(obj, "ADDED")
                },
                UpdateFunc: func(obj interface{}, newobj interface{}) {
			podEventParseAndSendData(obj, "MODIFIED")
                },
                DeleteFunc: func(obj interface{}) {
			podEventParseAndSendData(obj, "DELETED")
                },
        },
        )
        stop := make(chan struct{})
        go controller.Run(stop)
        return
}

func parseAndSendData(obj string, metaData  metav1.ObjectMeta, metaHeader metav1.TypeMeta, kind string, objtype string) {
	resp := make(map[string]string)
	resp["app_event_id"] = Generate() 
	resp["resource_type"] = kind 
	resp["resource_name"] = metaData.Name
	resp["resource_generations"] = string(metaData.Generation)
	resp["resource_id"] = string(metaData.UID)
	resp["app_environment_id"] = "joan-test"
	resp["app_environment_type"] = "Kubernetes"
	resp["type"] = objtype
	resp["trig_time"] = ""
	resp["server_group_id"] = metaData.Namespace
	resp["resource_version"] = metaData.ResourceVersion
	resp["message"] = obj
	respJson, err := json.Marshal(resp)
        if err != nil {
        	klog.Errorf("[ERROR] Failed to Marshal original object: %v", err)
        	fmt.Printf("[ERROR] Failed to Marshal original object: %v", err)
	}
	sendData(bytes.NewBuffer(respJson))
}

func sendData(data *bytes.Buffer){
	result, err := http.Post(endpointAPI, "application/json", data)
	if (err != nil){
		klog.Info("[INFO]", result, err)
		fmt.Printf("[INFO]", result, err)
	}
}
