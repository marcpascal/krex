package explorer

import (
	"fmt"
	"strings"
)

const ()

type DeploymentExplorer struct {
	Items               []*MenuItem
	PreviousItem        *MenuItem
	NamespaceToExplore  string
	DeploymentToExplore string
	PreviousExplorer    Explorable
}

func (n *DeploymentExplorer) List() error {
	n.Items = []*MenuItem{}
	n.Items = AddEdit(n.PreviousItem, n.Items)
	m := &MenuItem{}
	m.SetKind(podsLabel)
	m.SetName("Get Pods")
	n.Items = append(n.Items, m)
	o := &MenuItem{}
	o.SetKind(servicesLabel)
	o.SetName("Get Services")
	n.Items = append(n.Items, o)
	n.Items = AddGoBack(n.Items)
	n.Items = AddExit(n.Items)
	return nil
}

func (n *DeploymentExplorer) RunPrompt() (string, error) {
	var strs []string
	for _, item := range n.Items {
		strs = append(strs, item.GetReadable())
	}
	selection := transXY.Prompt("Select Deployment resources", strs)
	checkExitItem(selection)
	return selection, nil
}

func (n *DeploymentExplorer) Execute(selection string) error {
	item := NewMenuItemFromReadable(selection)
	switch item.GetKind() {
	case podsLabel:
		podsExplorer := &PodsExplorer{
			PreviousItem: item,
			Filters: map[string]string{
				"k8s-app": n.DeploymentToExplore,
			},
			NamespaceToExplore:   n.NamespaceToExplore,
			PreviousExplorer:     n,
			PreviousResourceName: n.DeploymentToExplore,
		}
		return Explore(podsExplorer)
	case servicesLabel:
		servicesExplorer := &ServicesExplorer{
			PreviousItem: item,
			Filters: map[string]string{
				"k8s-app": n.DeploymentToExplore,
			},
			NamespaceToExplore:   n.NamespaceToExplore,
			PreviousExplorer:     n,
			PreviousResourceName: n.DeploymentToExplore,
		}
		return Explore(servicesExplorer)
	case editLabel:
		Exec("kubectl", []string{"edit", "deployment", n.DeploymentToExplore, "-n", n.NamespaceToExplore})
		return Explore(n)
	case actionLabel:
		if strings.Contains(item.GetName(), "../") {
			return Explore(n.PreviousExplorer)
		}
		return fmt.Errorf("unknown action selection: %s, in explore_deployment.go", selection)
	case exitLabel:
		return Exit()
	default:
		return fmt.Errorf("unable to parse selection: %s, in explore_deployment.go", selection)
	}
}

func (n *DeploymentExplorer) Kind() string {
	return "deployment"
}
