package kustomize

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type UpgradeAction struct {
	Steps []UpgradeStep `yaml:"upgrade"`
}

// UpgradeStep represents the structure of an Upgrade step
type UpgradeStep struct {
	UpgradeArguments `yaml:"kustomize"`
}

// UpgradeArguments represent the arguments available to the Upgrade step
type UpgradeArguments struct {
	Step `yaml:",inline"`

	Namespace     string            `yaml:"namespace"`
	Name          string            `yaml:"name"`
	Kustomization string            `yaml:"kustomization"`
	Version       string            `yaml:"version"`
	Set           map[string]string `yaml:"set"`
	Values        []string          `yaml:"values"`
	Wait          bool              `yaml:"wait"`
	ResetValues   bool              `yaml:"resetValues"`
	ReuseValues   bool              `yaml:"reuseValues"`
}

// Upgrade issues a kustomize upgrade command for a release using the provided UpgradeArguments
func (m *Mixin) Upgrade() error {
	payload, err := m.getPayloadData()
	if err != nil {
		return err
	}

	var action UpgradeAction
	err = yaml.Unmarshal(payload, &action)
	if err != nil {
		return err
	}
	if len(action.Steps) != 1 {
		return errors.Errorf("expected a single step, but got %d", len(action.Steps))
	}
	step := action.Steps[0]

	cmd := m.NewCommand("kustomize", "upgrade", step.Name, step.Kustomization)

	if step.Namespace != "" {
		cmd.Args = append(cmd.Args, "--namespace", step.Namespace)
	}

	if step.Version != "" {
		cmd.Args = append(cmd.Args, "--version", step.Version)
	}

	if step.ResetValues {
		cmd.Args = append(cmd.Args, "--reset-values")
	}

	if step.ReuseValues {
		cmd.Args = append(cmd.Args, "--reuse-values")
	}

	if step.Wait {
		cmd.Args = append(cmd.Args, "--wait")
	}

	for _, v := range step.Values {
		cmd.Args = append(cmd.Args, "--values", v)
	}

	// sort the set consistently
	setKeys := make([]string, 0, len(step.Set))
	for k := range step.Set {
		setKeys = append(setKeys, k)
	}
	sort.Strings(setKeys)

	for _, k := range setKeys {
		cmd.Args = append(cmd.Args, "--set", fmt.Sprintf("%s=%s", k, step.Set[k]))
	}

	cmd.Stdout = m.Out
	cmd.Stderr = m.Err

	prettyCmd := fmt.Sprintf("%s %s", cmd.Path, strings.Join(cmd.Args, " "))
	_, err = fmt.Fprintln(m.Out, prettyCmd)
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("could not execute command, %s: %s", prettyCmd, err)
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}

	/*
		for _, output := range step.Outputs {

					//err = m.Context.WriteMixinOutputToFile(output.Name, val)
			if err != nil {
				return errors.Wrapf(err, "unable to write output '%s'", output.Name)
			}
		}

	*/
	return nil
}
