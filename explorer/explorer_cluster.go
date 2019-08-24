package explorer

import (
	"fmt"
	"os"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	namespaceLabel        = "Namespace"
	startNewTerminal      = "Terminal"
	startMidightCommander = "Midnight Commander"
)

type ClusterExplorer struct {
	Items []*MenuItem
}

func (n *ClusterExplorer) List() error {
	n.Items = []*MenuItem{}

	ns, err := k8sclient.CoreV1().Namespaces().List(v1.ListOptions{})
	if err != nil {
		return err
	}
	for _, item := range ns.Items {
		m := &MenuItem{}
		m.SetKind(namespaceLabel)
		m.SetName(item.Name)
		n.Items = append(n.Items, m)
	}
	n.Items = AddExit(n.Items)

	n.Items = append(n.Items, &MenuItem{
		kind: startNewTerminal,
		name: "Start a new detached console",
	})

	n.Items = append(n.Items, &MenuItem{
		kind: startMidightCommander,
		name: "Start Midnight Commander for directory operations",
	})

	// TODO add CRDs
	return nil
}

func (n *ClusterExplorer) RunPrompt() (string, error) {
	var strs []string
	for _, item := range n.Items {
		strs = append(strs, item.GetReadable())
	}
	selection := transXY.Prompt("Select cluster resource", strs)
	checkExitItem(selection)
	return selection, nil
}

var attr = os.ProcAttr{
	Dir: "/bin",
	Env: os.Environ(),
	Files: []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
	},
}

func (n *ClusterExplorer) Execute(selection string) error {
	item := NewMenuItemFromReadable(selection)
	switch item.GetKind() {
	case namespaceLabel:
		namespaceExplorer := &NamespaceExplorer{
			PreviousItem:       item,
			NamespaceToExplore: item.name,
			PreviousExplorer:   n,
		}
		return Explore(namespaceExplorer)
	case startMidightCommander:
		Exec("/usr/bin/mc", []string{""})
		return Explore(n)
	case startNewTerminal:
		os.StartProcess("/usr/bin/xterm", []string{""}, &attr)
		return Explore(n)
	case exitLabel:
		return Exit()
	default:
		return fmt.Errorf("unable to parse selection: %s", selection)
	}
}

func (n *ClusterExplorer) Kind() string {
	return "cluster"
}
