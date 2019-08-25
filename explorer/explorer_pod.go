package explorer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	logsLabel           = "Logs"
	execLabel           = "Exec"
	describeLabel       = "Describe"
	logAndDescribeLabel = "Logs and Describe"
)

type PodExplorer struct {
	Items              []*MenuItem
	PreviousItem       *MenuItem
	NamespaceToExplore string
	PodToExplore       string
	PreviousExplorer   Explorable
}

func (n *PodExplorer) List() error {
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

func (n *PodExplorer) RunPrompt() (string, error) {
	var strs []string
	for _, item := range n.Items {
		strs = append(strs, item.GetReadable())
	}
	selection := transXY.Prompt("Select Pod resources", strs)
	checkExitItem(selection)
	return selection, nil
}

func (n *PodExplorer) Execute(selection string) error {
	item := NewMenuItemFromReadable(selection)
	switch item.GetKind() {
	case logAndDescribeLabel:
		cmd1 := exec.Command("/usr/bin/xterm", "-geometry", "80x40", "-hold", "-e", "/usr/bin/kubectl describe pod --namespace  "+n.NamespaceToExplore+" "+n.PodToExplore)
		cmd1.Stdin = os.Stdin
		cmd1.Stdout = os.Stdout
		cmd1.Stderr = os.Stderr
		cmd1.Start()

		cmd2 := exec.Command("/usr/bin/xterm", "-geometry", "80x40", "-hold", "-e", "/usr/bin/kubectl logs "+n.PodToExplore+" -n "+n.NamespaceToExplore+" -f")
		cmd2.Stdin = os.Stdin
		cmd2.Stdout = os.Stdout
		cmd2.Stderr = os.Stderr
		cmd2.Start()
		return Explore(n)
	case describeLabel:
		cmd := exec.Command("/usr/bin/xterm", "-geometry", "80x40", "-hold", "-e", "/usr/bin/kubectl describe pod --namespace  "+n.NamespaceToExplore+" "+n.PodToExplore)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
		return Explore(n)
	case execLabel:
		Exec("kubectl", []string{"exec", "-it", "--namespace", n.NamespaceToExplore, n.PodToExplore, "sh"})
		return Explore(n)
	case logsLabel:
		cmd2 := exec.Command("/usr/bin/xterm", "-geometry", "80x40", "-hold", "-e", "/usr/bin/kubectl logs "+n.PodToExplore+" -n "+n.NamespaceToExplore+" -f")
		cmd2.Stdin = os.Stdin
		cmd2.Stdout = os.Stdout
		cmd2.Stderr = os.Stderr
		cmd2.Start()
		return Explore(n)
	case editLabel:
		Exec("kubectl", []string{"edit", "pods", n.PodToExplore, "-n", n.NamespaceToExplore})
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

func (n *PodExplorer) Kind() string {
	return "pod"
}
