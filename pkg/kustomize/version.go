package kustomize

import (
	"fmt"

	"github.com/donmstewart/porter-kustomize/pkg"
)

func (m *Mixin) PrintVersion() {
	fmt.Fprintf(m.Out, "kustomize mixin %s (%s)\n", pkg.Version, pkg.Commit)
}