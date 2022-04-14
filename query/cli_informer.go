package query

import (
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"time"
)

var (
	RestCliConfig      *rest.Config
	ClientSet          *kubernetes.Clientset
	InfoFact           informers.SharedInformerFactory
	ReloadInfoFactChan chan bool
)

func InitClient() {
	// get config flag from generic cli options (like: kubectl)
	cfgFlags := genericclioptions.NewConfigFlags(true)
	config, err := cfgFlags.ToRawKubeConfigLoader().ClientConfig()
	if err != nil {
		log.Fatalln(err)
	}
	RestCliConfig = config
	ClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}
}

type GenericHandler struct {
}

func (this *GenericHandler) OnAdd(obj interface{}) {
}

func (this *GenericHandler) OnUpdate(oldObj, newObj interface{}) {
}

func (this *GenericHandler) OnDelete(obj interface{}) {
}

func InitInformerCache() {
	InfoFact = informers.NewSharedInformerFactory(ClientSet, time.Second*600)
	// pods
	InfoFact.Core().V1().Pods().Informer().AddEventHandler(&GenericHandler{})
	// deployments
	InfoFact.Apps().V1().Deployments().Informer().AddEventHandler(&GenericHandler{})
	// services
	InfoFact.Core().V1().Services().Informer().AddEventHandler(&GenericHandler{})
	// namespace
	InfoFact.Core().V1().Namespaces().Informer().AddEventHandler(&GenericHandler{})
	// daemonset
	InfoFact.Apps().V1().DaemonSets().Informer().AddEventHandler(&GenericHandler{})
	// statefulset
	InfoFact.Apps().V1().StatefulSets().Informer().AddEventHandler(&GenericHandler{})
	// jobs
	InfoFact.Batch().V1().Jobs().Informer().AddEventHandler(&GenericHandler{})
	// replicaset
	InfoFact.Apps().V1().ReplicaSets().Informer().AddEventHandler(&GenericHandler{})
	// events
	InfoFact.Core().V1().Events().Informer().AddEventHandler(&GenericHandler{})
	// configmaps
	InfoFact.Core().V1().ConfigMaps().Informer().AddEventHandler(&GenericHandler{})
	// nodes
	InfoFact.Core().V1().Nodes().Informer().AddEventHandler(&GenericHandler{})

	// start
	InfoFact.Start(wait.NeverStop)
	InfoFact.WaitForCacheSync(wait.NeverStop)
}
