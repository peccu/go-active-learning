package command

import (
	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning/lib/annotation"
	"github.com/syou6162/go-active-learning/lib/diagnosis"
)

var Commands = []cli.Command{
	annotation.CommandAnnotate,
	diagnosis.CommandDiagnose,
}
