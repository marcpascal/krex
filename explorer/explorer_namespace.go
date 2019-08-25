package explorer

import (
	"fmt"

	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	statefulSetLabel = "StatefulSet"
	deploymentLabel  = "Deployment"
	daemonSetLabel   = "DaemonSet"
	replicasetLabel  = "Replicaset"
	debugLabel       = "Debug"
	actionLabel      = "Action"
)

type NamespaceExplorer struct {
	Items              []*MenuItem
	PreviousItem       *MenuItem
	NamespaceToExplore string
	PreviousExplorer   Explorable
}

func (n *NamespaceExplorer) List() error {
	n.Items = []*MenuItem{}
	n.Items = AddEdit(n.PreviousItem, n.Items)

	// StatefulSet
	ss, err := k8sclient.AppsV1().StatefulSets(n.NamespaceToExplore).List(v1.ListOptions{})
	if err != nil {
		return err
	}
	for _, item := range ss.Items {
		m := &MenuItem{}
		m.SetName(item.Name)
		m.SetKind(statefulSetLabel)
		n.Items = append(n.Items, m)
	}

	// Deployment
	ds, err := k8sclient.AppsV1().Deployments(n.NamespaceToExplore).List(v1.ListOptions{})
	if err != nil {
		return err
	}
	for _, item := range ds.Items {
		m := &MenuItem{}
		m.SetName(item.Name)
		m.SetKind(deploymentLabel)
		n.Items = append(n.Items, m)
	}

	// DaemonSet
	dss, err := k8sclient.AppsV1().DaemonSets(n.NamespaceToExplore).List(v1.ListOptions{})
	if err != nil {
		return err
	}
	for _, item := range dss.Items {
		m := &MenuItem{}
		m.SetName(item.Name)
		m.SetKind(daemonSetLabel)
		n.Items = append(n.Items, m)
	}

	// Replicaset
	dreplicatset, err := k8sclient.AppsV1().ReplicaSets(n.NamespaceToExplore).List(v1.ListOptions{})
	if err != nil {
		return err
	}
	for _, item := range dreplicatset.Items {
		m := &MenuItem{}
		m.SetName(item.Name)
		m.SetKind(replicasetLabel)
		n.Items = append(n.Items, m)
	}

	m := &MenuItem{}
	m.SetKind(debugLabel)
	m.SetName(fmt.Sprintf("Run a debugging pod in the Namespace and shell exec [%s]", options.ShellExecImage))
	n.Items = append(n.Items, m)
	n.Items = AddGoBack(n.Items)
	n.Items = AddExit(n.Items)
	return nil
}

func (n *NamespaceExplorer) RunPrompt() (string, error) {
	var strs []string
	for _, item := range n.Items {
		strs = append(strs, item.GetReadable())
	}
	selection := transXY.Prompt("Select Namespace resource", strs)
	checkExitItem(selection)
	return selection, nil
}

func (n *NamespaceExplorer) Execute(selection string) error {
	item := NewMenuItemFromReadable(selection)
	switch item.GetKind() {
	case daemonSetLabel:
		daemonSetExplorer := &DaemonSetExplorer{
			PreviousItem:       item,
			DaemonSetToExplore: item.name,
			NamespaceToExplore: n.NamespaceToExplore,
			PreviousExplorer:   n,
		}
		return Explore(daemonSetExplorer)
	case deploymentLabel:
		deploymentExplorer := &DeploymentExplorer{
			PreviousItem:        item,
			DeploymentToExplore: item.name,
			NamespaceToExplore:  n.NamespaceToExplore,
			PreviousExplorer:    n,
		}
		return Explore(deploymentExplorer)
	case statefulSetLabel:
		statefulSetExplorer := &StatefulSetExplorer{
			PreviousItem:         item,
			StatefulSetToExplore: item.name,
			NamespaceToExplore:   n.NamespaceToExplore,
			PreviousExplorer:     n,
		}
		return Explore(statefulSetExplorer)
	case actionLabel:
		if strings.Contains(item.GetName(), "../") {
			if n.PreviousExplorer != nil {
				return Explore(n.PreviousExplorer)
			}
		}
		return fmt.Errorf("Invalid action: %s", selection)
	case debugLabel:
		// Deploy a pod and exec into it
		Exec("kubectl", []string{"run", "-n", n.NamespaceToExplore, "-i", "--tty", "krex-debug-pod", "--image", options.ShellExecImage, "--", "sh"})
		return Explore(n)
	case editLabel:
		Exec("kubectl", []string{"edit", "namespace", n.NamespaceToExplore})
		return Explore(n)
	case exitLabel:
		return Exit()
	default:
		return fmt.Errorf("unable to parse selection: %s", selection)
	}
}

func (n *NamespaceExplorer) Kind() string {
	return "namespace"
}
