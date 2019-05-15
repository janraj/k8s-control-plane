package main
  
import (
        "log"
	clt "k8s-control-plane/handler/api/v1"
)

func InitCitrixControlPlane() error {
        log.Println("Initializing CCP")
        return nil
}
func StartCitrixControlPlane() {
	log.Println("StartCitrixControlPlane")
	api, err := clt.CreateK8sApiserverClient()
        if err != nil {
                log.Println("[ERROR] K8s Client Error", err)
        }
	clt.EndpointWatcher(api)
	clt.ServiceWatcher(api)
	clt.PodWatcher(api)
	clt.IngressWatcher(api)
	select {}
}
func main() {
        InitCitrixControlPlane()
        StartCitrixControlPlane()
}
