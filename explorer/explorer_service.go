package explorer

import (
	"fmt"
	"strings"
)

type ServiceExplorer struct {
	Items              []*MenuItem
	PreviousItem       *MenuItem
	NamespaceToExplore string
	ServiceToExplore   string
	PreviousExplorer   Explorable
}

func (n *ServiceExplorer) List() error {
	n.Items = []*MenuItem{}

	n.Items = AddEdit(n.PreviousItem, n.Items)

	n.Items = append(n.Items, &MenuItem{
		kind: logsLabel,
		name: "Tail logs",
	})

	n.Items = append(n.Items, &MenuItem{
		kind: execLabel,
		name: "Shell exec (sh) into Pod",
	})

	n.Items = append(n.Items, &MenuItem{
		kind: describeLabel,
		name: "Describe the Pod",
	})

	n.Items = append(n.Items, &MenuItem{
		kind: logAndDescribeLabel,
		name: "Describe the Pod and then tail the logs",
	})

	// Logs
	// Exec
	// Describe
	// Log and Describe

	n.Items = AddGoBack(n.Items)
	n.Items = AddExit(n.Items)
	return nil
}

func (n *ServiceExplorer) RunPrompt() (string, error) {
	var strs []string
	for _, item := range n.Items {
		strs = append(strs, item.GetReadable())
	}
	selection := transXY.Prompt("Select Pod resources", strs)
	checkExitItem(selection)
	return selection, nil
}

func (n *ServiceExplorer) Execute(selection string) error {
	item := NewMenuItemFromReadable(selection)
	switch item.GetKind() {
	case logAndDescribeLabel:
		Exec("kubectl", []string{"describe", "service", "--namespace", n.NamespaceToExplore, n.ServiceToExplore})
		Exec("kubectl", []string{"logs", n.ServiceToExplore, "-n", n.NamespaceToExplore, "-f"})
		return Explore(n)
	case describeLabel:
		Exec("kubectl", []string{"describe", "service", "--namespace", n.NamespaceToExplore, n.ServiceToExplore})
		return Explore(n)
	case execLabel:
		Exec("kubectl", []string{"exec", "-it", "--namespace", n.NamespaceToExplore, n.ServiceToExplore, "sh"})
		return Explore(n)
	case logsLabel:
		Exec("kubectl", []string{"logs", n.ServiceToExplore, "-n", n.NamespaceToExplore, "-f"})
		return Explore(n)
	case editLabel:
		Exec("kubectl", []string{"edit", "pods", n.ServiceToExplore, "-n", n.NamespaceToExplore})
		return Explore(n)
	case actionLabel:
		if strings.Contains(item.GetName(), "../") {
			return Explore(n.PreviousExplorer)
		}
		return fmt.Errorf("unknown action selection: %s", selection)
	case exitLabel:
		return Exit()
	default:
		return fmt.Errorf("unable to parse selection: %s", selection)
	}
}

func (n *ServiceExplorer) Kind() string {
	return "service"
}
