package query

import (
	"encoding/json"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/yaml"
)

func ParserLabels(labels_ []map[string]string) labels.Selector {
	var l labels.Selector
	if len(labels_) != 0 {
		l = labels.SelectorFromSet(labels_[0])
	} else {
		l = labels.Everything()
	}

	return l

}

func FetchAllNamespace() []*v1.Namespace {
	namespaces, err := InfoFact.Core().V1().Namespaces().Lister().List(labels.Everything())
	WrapError(err)
	return namespaces
}

func FetchPods(ns string, labels_ ...map[string]string) []*v1.Pod {
	l := ParserLabels(labels_)
	if ns == AllNamespace {
		pods, err := InfoFact.Core().V1().Pods().Lister().List(l)
		WrapError(err)
		return pods
	}
	pods, err := InfoFact.Core().V1().Pods().Lister().Pods(ns).List(l)
	WrapError(err)
	return pods
}

func FetchPodWithName(name, ns string) *v1.Pod {
	pod, err := InfoFact.Core().V1().Pods().Lister().Pods(ns).Get(name)
	WrapError(err)
	return pod
}

func FetchDeployments(ns string, labels_ ...map[string]string) []*appsv1.Deployment {
	l := ParserLabels(labels_)
	if ns == AllNamespace {
		deploys, err := InfoFact.Apps().V1().Deployments().Lister().List(l)
		WrapError(err)
		return deploys
	}
	deploys, err := InfoFact.Apps().V1().Deployments().Lister().Deployments(ns).List(l)
	WrapError(err)
	return deploys
}

func FetchDeploymentWithName(name, ns string) *appsv1.Deployment {
	deploy, err := InfoFact.Apps().V1().Deployments().Lister().Deployments(ns).Get(name)
	WrapError(err)
	return deploy
}

func FetchServices(ns string, labels_ ...map[string]string) []*v1.Service {
	l := ParserLabels(labels_)
	if ns == AllNamespace {
		svcs, err := InfoFact.Core().V1().Services().Lister().List(l)
		WrapError(err)
		return svcs
	}
	svcs, err := InfoFact.Core().V1().Services().Lister().Services(ns).List(l)
	WrapError(err)
	return svcs
}

func FetchServiceWithName(name, ns string) *v1.Service {
	obj, err := InfoFact.Core().V1().Services().Lister().Services(ns).Get(name)
	WrapError(err)
	return obj
}

func FetchReplicaSets(ns string, labels_ ...map[string]string) []*appsv1.ReplicaSet {
	l := ParserLabels(labels_)
	if ns == AllNamespace {
		objs, err := InfoFact.Apps().V1().ReplicaSets().Lister().List(l)
		WrapError(err)
		return objs
	}
	objs, err := InfoFact.Apps().V1().ReplicaSets().Lister().ReplicaSets(ns).List(l)
	WrapError(err)
	return objs
}

func FetchReplicaSetWithName(name, ns string) *appsv1.ReplicaSet {
	obj, err := InfoFact.Apps().V1().ReplicaSets().Lister().ReplicaSets(ns).Get(name)
	WrapError(err)
	return obj
}

func FetchDaemonSets(ns string, labels_ ...map[string]string) []*appsv1.DaemonSet {
	l := ParserLabels(labels_)
	if ns == AllNamespace {
		objs, err := InfoFact.Apps().V1().DaemonSets().Lister().List(l)
		WrapError(err)
		return objs
	}
	objs, err := InfoFact.Apps().V1().DaemonSets().Lister().DaemonSets(ns).List(l)
	WrapError(err)
	return objs
}

func FetchDaemonSetWithName(name, ns string) *appsv1.DaemonSet {
	obj, err := InfoFact.Apps().V1().DaemonSets().Lister().DaemonSets(ns).Get(name)
	WrapError(err)
	return obj
}

func FetchStatefulSets(ns string, labels_ ...map[string]string) []*appsv1.StatefulSet {
	l := ParserLabels(labels_)
	if ns == AllNamespace {
		objs, err := InfoFact.Apps().V1().StatefulSets().Lister().List(l)
		WrapError(err)
		return objs
	}
	objs, err := InfoFact.Apps().V1().StatefulSets().Lister().StatefulSets(ns).List(l)
	WrapError(err)
	return objs
}

func FetchStatefulSetWithName(name, ns string) *appsv1.StatefulSet {
	obj, err := InfoFact.Apps().V1().StatefulSets().Lister().StatefulSets(ns).Get(name)
	WrapError(err)
	return obj
}

func FetchJobs(ns string, labels_ ...map[string]string) []*batchv1.Job {
	l := ParserLabels(labels_)
	if ns == AllNamespace {
		objs, err := InfoFact.Batch().V1().Jobs().Lister().List(l)
		WrapError(err)
		return objs
	}
	objs, err := InfoFact.Batch().V1().Jobs().Lister().Jobs(ns).List(l)
	WrapError(err)
	return objs
}

func FetchJobWithName(name, ns string) *batchv1.Job {
	obj, err := InfoFact.Batch().V1().Jobs().Lister().Jobs(ns).Get(name)
	WrapError(err)
	return obj
}

func FetchEvents(uid types.UID, labels_ ...map[string]string) []*v1.Event {
	l := ParserLabels(labels_)
	var res []*v1.Event
	events, err := InfoFact.Core().V1().Events().Lister().List(l)
	Logger("Events: ", events)
	WrapError(err)
	for _, e := range events {
		if e.InvolvedObject.UID == uid {
			res = append(res, e)
		}
	}
	return res
}

func FetchConfigMaps(ns string, labels_ ...map[string]string) []*v1.ConfigMap {
	l := ParserLabels(labels_)
	if ns == AllNamespace {
		objs, err := InfoFact.Core().V1().ConfigMaps().Lister().List(l)
		WrapError(err)
		return objs
	}
	objs, err := InfoFact.Core().V1().ConfigMaps().Lister().ConfigMaps(ns).List(l)
	WrapError(err)
	return objs
}

func FetchConfigMapWithName(name, ns string) *v1.ConfigMap {
	obj, err := InfoFact.Core().V1().ConfigMaps().Lister().ConfigMaps(ns).Get(name)
	WrapError(err)
	return obj
}

func RuntimeObject2Yaml(obj runtime.Object) string {
	if obj == nil {
		return ""
	}
	b, err := yaml.Marshal(obj)
	WrapError(err)
	return string(b)
}

func RuntimeObject2Json(obj runtime.Object) string {
	if obj == nil {
		return ""
	}
	b, err := json.MarshalIndent(obj, "", "  ")
	WrapError(err)
	return string(b)
}
