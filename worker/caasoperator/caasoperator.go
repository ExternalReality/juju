// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package caasoperator

import (
	"path/filepath"

	"github.com/juju/errors"
	"github.com/juju/loggo"
	"github.com/juju/utils/clock"
	"gopkg.in/juju/names.v2"
	"gopkg.in/juju/worker.v1"

	agenttools "github.com/juju/juju/agent/tools"
	jworker "github.com/juju/juju/worker"
	"github.com/juju/juju/worker/caasoperator/commands"
	"github.com/juju/juju/worker/catacomb"
)

var logger = loggo.GetLogger("juju.worker.caasoperator")

// A CaasOperatorExecutionObserver gets the appropriate methods called when a hook
// is executed and either succeeds or fails.  Missing hooks don't get reported
// in this way.
type CaasOperatorExecutionObserver interface {
	HookCompleted(hookName string)
	HookFailed(hookName string)
}

// caasOperator implements the capabilities of the caasoperator agent. It is not intended to
// implement the actual *behaviour* of the caasoperator agent; that responsibility is
// delegated to Mode values, which are expected to react to events and direct
// the caasoperator's responses to them.
type caasOperator struct {
	catacomb catacomb.Catacomb
	config   Config
}

// Config hold the configuration for a caasoperator worker.
type Config struct {
	// Application holds the name of the application that
	// this CAAS operator manages.
	Application string

	// DataDir holds the path to the Juju "data directory",
	// i.e. "/var/lib/juju" (by default). The CAAS operator
	// expects to find the jujud binary at <data-dir>/tools/jujud.
	DataDir string

	// Clock holds the clock to be used by the CAAS operator
	// for time-related operations.
	Clock clock.Clock

	// StatusSetter is an interface used for setting the
	// application status.
	StatusSetter StatusSetter
}

func (config Config) Validate() error {
	if !names.IsValidApplication(config.Application) {
		return errors.NotValidf("application name %q", config.Application)
	}
	if config.DataDir == "" {
		return errors.NotValidf("missing DataDir")
	}
	if config.Clock == nil {
		return errors.NotValidf("missing Clock")
	}
	if config.StatusSetter == nil {
		return errors.NotValidf("missing StatusSetter")
	}
	return nil
}

// NewWorker creates a new worker which will install and operate a
// CaaS-based application, by executing hooks and operations in
// response to application state changes.
func NewWorker(config Config) (worker.Worker, error) {
	if err := config.Validate(); err != nil {
		return nil, errors.Trace(err)
	}

	op := &caasOperator{
		config: config,
	}
	if err := op.init(); err != nil {
		if err == jworker.ErrTerminateAgent {
			return nil, err
		}
		return nil, errors.Annotatef(err,
			"failed to initialize caasoperator for %q",
			op.config.Application,
		)
	}

	if err := catacomb.Invoke(catacomb.Plan{
		Site: &op.catacomb,
		Work: op.loop,
	}); err != nil {
		return nil, errors.Trace(err)
	}
	return op, nil
}

func (op *caasOperator) loop() (err error) {
	for {
		select {
		case <-op.catacomb.Dying():
			return op.catacomb.ErrDying()
		}
	}
}

func (op *caasOperator) agentBinaryDir() string {
	return filepath.Join(op.config.DataDir, "tools")
}

func (op *caasOperator) init() (err error) {
	agentBinaryDir := op.agentBinaryDir()
	logger.Debugf("creating caas operator symlinks in %v", agentBinaryDir)
	if err := agenttools.EnsureSymlinks(
		agentBinaryDir,
		agentBinaryDir,
		commands.CommandNames(),
	); err != nil {
		return err
	}
	return nil
}

// Kill is part of the worker.Worker interface.
func (op *caasOperator) Kill() {
	op.catacomb.Kill(nil)
}

// Wait is part of the worker.Worker interface.
func (op *caasOperator) Wait() error {
	return op.catacomb.Wait()
}
