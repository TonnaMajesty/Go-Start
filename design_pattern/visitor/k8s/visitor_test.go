package k8s

import (
	"testing"
)

func TestVisitor(t *testing.T) {
	info := Info{}
	var v DecoratedVisitor = &info
	v = LogVisitor{v}
	v = NameVisitor{v}
	v = OtherThingsVisitor{v}

	// loadfile 是真正的访问者
	loadFile := func(info *Info, err error) error {
		info.Name = "Hao Chen"
		info.Namespace = "MegaEase"
		info.OtherThings = "We are running as remote team."
		return nil
	}
	v.Accept(loadFile)
}