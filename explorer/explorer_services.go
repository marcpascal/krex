package explorer

import (
	"fmt"
	"strings"

	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	serviceLabel = "Service"
)

type ServicesExplorer struct {
	Items                []*MenuItem
	PreviousItem         *MenuItem
	NamespaceToExplore   string
	Filters              map[string]string
	PreviousExplorer     Explorable
	PreviousResourceName string
}

func (n *ServicesExplorer) List() error {
	n.Items = []*MenuItem{}

	// List Services
	services, err := k8sclient.Core().Services(n.NamespaceToExplore).List(v1.ListOptions{})
	if err != nil {
		return err
	}

	var servicelist []apiv1.Service
	switch n.PreviousExplorer.Kind() {
	case "deployment":
		// get replicasets
		replicaSetMap := make(map[string]bool)
		replicaSets, err := k8sclient.Extensions().ReplicaSets(n.NamespaceToExplore).List(v1.ListOptions{})
		if err != nil {
			return err
		}

		for _, replicaset := range replicaSets.Items {
			for _, owner := range replicaset.GetOwnerReferences() {
				if strings.ToLower(owner.Name) == n.PreviousResourceName {
					replicaSetMap[replicaset.GetName()] = true
					break
				}
			}
		}

		// get services
		servicelist = getOwnerServices(services, replicaSetMap, "replicaset")

	case "statefulset":
		servicelist = getOwnerServices(services, map[string]bool{n.PreviousResourceName: true}, "statefulset")

	case "daemonset":
		servicelist = getOwnerServices(services, map[string]bool{n.PreviousResourceName: true}, "daemonset")
	}

	for _, item := range servicelist {
		m := &MenuItem{}
		m.SetName(item.Name)
		m.SetKind(serviceLabel)
		n.Items = append(n.Items, m)
	}
	n.Items = AddGoBack(n.Items)
	n.Items = AddExit(n.Items)
	return nil
}

// getOwnerServices is a functio that gets pods from map of owners
func getOwnerServices(services *apiv1.ServiceList, owners map[string]bool, resourceKindToMatch string) []apiv1.Service {
	var servicelist []apiv1.Service

	for _, service := range services.Items {
		for _, owner := range service.GetOwnerReferences() {
			if _, ok := owners[owner.Name]; ok && strings.ToLower(owner.Kind) == resourceKindToMatch {
				servicelist = append(servicelist, service)
				break
			}
		}
	}

	return servicelist
}

func (n *ServicesExplorer) RunPrompt() (string, error) {
	var strs []string
	for _, item := range n.Items {
		strs = append(strs, item.GetReadable())
	}
	selection := transXY.Prompt("Select Service resource", strs)
	checkExitItem(selection)
	return selection, nil
}

func (n *ServicesExplorer) Execute(selection string) error {
	item := NewMenuItemFromReadable(selection)
	switch item.GetKind() {
	case serviceLabel:
		serviceExplorer := &ServiceExplorer{
			PreviousItem:       item,
			ServiceToExplore:   item.GetName(),
			NamespaceToExplore: n.NamespaceToExplore,
			PreviousExplorer:   n,
		}
		return Explore(serviceExplorer)
	case actionLabel:
		if strings.Contains(item.GetName(), "../") {
			return Explore(n.PreviousExplorer)
		}
		return fmt.Errorf("unknown action selection: %s, in explore_services.go", selection)
	case exitLabel:
		return Exit()
	default:
		return fmt.Errorf("unable to parse selection: %s, in explore_services.go", selection)
	}
}

func (n *ServicesExplorer) Kind() string {
	return "services"
}
