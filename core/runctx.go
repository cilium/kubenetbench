package core

import (
	"fmt"
	"log"
	"os"
	"time"
)

// RunCtx is the context for a benchmark run
type RunCtx struct {
	id          string    // id identifies the run
	dir         string    // directory to store results/etc.
	cliAffinity string    // client affinity
	quiet       bool      // supress output
	cleanup     bool      // perform cleanup: remove k8s entitites (pods, policies, etc.)
	benchmark   Benchmark // underlying benchmark interface
}

// NewRunCtx creates a new RunCtx
func NewRunCtx(
	rid string,
	ridDirBase string,
	cliAffinity string,
	quiet bool,
	cleanup bool,
	benchmark Benchmark,
) *RunCtx {
	datestr := time.Now().Format("20060102-150405")
	rundir := fmt.Sprintf("%s/%s-%s", ridDirBase, rid, datestr)
	err := os.Mkdir(rundir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	return &RunCtx{
		id:          rid,
		dir:         rundir,
		cliAffinity: cliAffinity,
		quiet:       quiet,
		cleanup:     cleanup,
		benchmark:   benchmark,
	}
}
