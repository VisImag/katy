package pods

import (
	"errors"
	"log"
	"strings"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

var kubeClient *kubernetes.Clientset
var podsGetter v1.PodsGetter

func PushClientSet(client *kubernetes.Clientset) {
	kubeClient = client
	podsGetter = kubeClient.CoreV1()
}

func getPods(namespace string) v1.PodInterface {
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}
	pods := podsGetter.Pods(namespace)
	return pods
}

// Get number of pods present in a namespace
func GetNumOfPods(namespace string) (int32, error) {
	pods := getPods(namespace)
	list, err := pods.List(metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return int32(len(list.Items)), nil
}

// Check if the given pod is running
func CheckIfRunning(namespace string, podName string) (bool, error) {
	pods := getPods(namespace)
	pod, err := pods.Get(podName, metav1.GetOptions{})
	if pod == nil {
		return false, errors.New("Pod not present")
	}
	if err != nil {
		log.Println(err)
		return false, err
	}
	if pod.Status.Phase == apiv1.PodRunning {
		return true, nil
	}
	return false, nil
}

// Check if the given pod is ready
func CheckIfReady(namespace string, podName string) (bool, error) {
	pods := getPods(namespace)
	pod, err := pods.Get(podName, metav1.GetOptions{})
	if pod == nil {
		return false, errors.New("Pod not present")
	}
	if err != nil {
		log.Println(err)
		return false, err
	}
	var flag bool = false
	conditions := pod.Status.Conditions
	for _, c := range conditions {
		if c.Type == apiv1.PodReady && c.Status == apiv1.ConditionTrue {
			flag = true
			break
		}
	}
	if flag {
		return true, nil
	} else {
		return false, nil
	}
}

// Get phase of a given pod
func GetPodPhase(namespace string, podName string) (string, error) {
	pods := getPods(namespace)
	pod, err := pods.Get(podName, metav1.GetOptions{})
	if pod == nil {
		return "", errors.New("Pod not present")
	}
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(pod.Status.Phase), nil
}

func getCurrentConditionDetails(conditions []apiv1.PodCondition) (string, string) {
	var time metav1.Time = metav1.Time{}
	var condition string
	var reason *strings.Builder = new(strings.Builder)
	for _, c := range conditions {
		if time.IsZero() && c.Status == apiv1.ConditionTrue {
			time = c.LastTransitionTime
			condition = string(c.Type)
			reason.WriteString(c.Reason)
			reason.WriteString("-")
			reason.WriteString(c.Message)
		}
		if time.Before(&c.LastTransitionTime) {
			time = c.LastTransitionTime
			condition = string(c.Type)
			reason.Reset()
			reason.WriteString(c.Reason)
			reason.WriteString("-")
			reason.WriteString(c.Message)
		}
	}
	return condition, reason.String()
}

// Get status/condition of a give pod
func GetPodStatus(namespace string, podName string) (string, error) {
	pods := getPods(namespace)
	pod, err := pods.Get(podName, metav1.GetOptions{})
	if pod == nil {
		return "", errors.New("Pod not present")
	}
	if err != nil {
		log.Println(err)
		return "", err
	}
	condition, _ := getCurrentConditionDetails(pod.Status.Conditions)
	return condition, nil
}

// Get reason for the pod condition
func GetConditionReason(namespace string, podName string) (string, error) {
	pods := getPods(namespace)
	pod, err := pods.Get(podName, metav1.GetOptions{})
	if pod == nil {
		return "", errors.New("Pod not present")
	}
	if err != nil {
		log.Println(err)
		return "", err
	}
	_, reason := getCurrentConditionDetails(pod.Status.Conditions)
	return reason, nil
}

// Get timestamp when the pod was started
func GetPodStartTime(namespace string, podName string) (string, error) {
	pods := getPods(namespace)
	pod, err := pods.Get(podName, metav1.GetOptions{})
	if pod == nil {
		return "", errors.New("Pod not present")
	}
	if err != nil {
		log.Println(err)
		return "", err
	}
	return pod.Status.StartTime.String(), nil
}
