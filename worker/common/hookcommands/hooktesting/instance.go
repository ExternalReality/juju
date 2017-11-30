// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package hooktesting

import (
	"github.com/juju/errors"

	"github.com/juju/juju/worker/common/hookcommands"
)

// Instance holds the values for the hook context.
type Instance struct {
	AvailabilityZone string
	RebootPriority   *hookcommands.RebootPriority
}

// ContextInstance is a test double for hooks.ContextInstance.
type ContextInstance struct {
	contextBase
	info *Instance
}

// AvailabilityZone implements hooks.ContextInstance.
func (c *ContextInstance) AvailabilityZone() (string, error) {
	c.stub.AddCall("AvailabilityZone")

	return c.info.AvailabilityZone, c.stub.NextErr()
}

// RequestReboot implements hooks.ContextInstance.
func (c *ContextInstance) RequestReboot(priority hookcommands.RebootPriority) error {
	c.stub.AddCall("RequestReboot", priority)
	if err := c.stub.NextErr(); err != nil {
		return errors.Trace(err)
	}

	c.info.RebootPriority = &priority
	return nil
}