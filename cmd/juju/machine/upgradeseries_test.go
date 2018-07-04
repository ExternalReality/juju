// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package machine_test

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/juju/cmd"
	"github.com/juju/cmd/cmdtesting"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/cmd/juju/machine"
	"github.com/juju/juju/cmd/juju/machine/mocks"
)

type UpgradeSeriesSuite struct {
	ctx *cmd.Context
}

var _ = gc.Suite(&UpgradeSeriesSuite{})

const machineArg = "1"
const seriesArg = "xenial"

func (s *UpgradeSeriesSuite) SetUpTest(c *gc.C) {}

func runUpgradeSeriesCommand(c *gc.C, args ...string) error {
	err := runUpgradeSeriesCommandWithConfirmation(c, "y", args...)
	return err
}

func runUpgradeSeriesCommandWithConfirmation(c *gc.C, confirmation string, args ...string) error {
	mockController := gomock.NewController(c)
	mockUpgradeSeriesAPI := mocks.NewMockUpgradeMachineSeriesAPI(mockController)
	// Stub for CLI arg testing
	mockUpgradeSeriesAPI.EXPECT().UpgradeSeriesPrepare(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	return runUpgradeSeriesCommandGivenMock(c, mockUpgradeSeriesAPI, confirmation, args...)
}

func runUpgradeSeriesCommandGivenMock(c *gc.C, mockUpgradeSeriesAPI *mocks.MockUpgradeMachineSeriesAPI, confirmation string, args ...string) error {
	var stdin, stdout, stderr bytes.Buffer
	ctx, err := cmd.DefaultContext()
	s.ctx = ctx
	c.Assert(err, jc.ErrorIsNil)
	s.ctx.Stderr = &stderr
	s.ctx.Stdout = &stdout
	s.ctx.Stdin = &stdin
	stdin.WriteString(confirmation)

	com := machine.NewUpgradeSeriesCommandForTest(mockUpgradeSeriesAPI)

	err = cmdtesting.InitCommand(com, args)
	if err != nil {
		return err
	}

	err = com.Run(ctx)
	if err != nil {
		return err
	}

	if stderr.String() != "" {
		return errors.New(stderr.String())
	}

	return nil
}

func (s *UpgradeSeriesSuite) TestPrepareCommand(c *gc.C) {
	mockController := gomock.NewController(c)
	defer mockController.Finish()
	mockUpgradeSeriesAPI := mocks.NewMockUpgradeMachineSeriesAPI(mockController)
	mockUpgradeSeriesAPI.EXPECT().UpgradeSeriesPrepare(machineArg, seriesArg, gomock.Eq(false))

	err := runUpgradeSeriesCommandGivenMock(c, mockUpgradeSeriesAPI, "y", machine.PrepareCommand, machineArg, seriesArg)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *UpgradeSeriesSuite) TestPrepareCommandShouldAcceptForceOption(c *gc.C) {
	mockController := gomock.NewController(c)
	defer mockController.Finish()
	mockUpgradeSeriesAPI := mocks.NewMockUpgradeMachineSeriesAPI(mockController)
	mockUpgradeSeriesAPI.EXPECT().UpgradeSeriesPrepare(machineArg, seriesArg, gomock.Eq(true))

	err := runUpgradeSeriesCommandGivenMock(c, mockUpgradeSeriesAPI, "y", machine.PrepareCommand, machineArg, seriesArg, "--force")
	c.Assert(err, jc.ErrorIsNil)
}

func (s *UpgradeSeriesSuite) TestPrepareCommandShouldAbortOnFailedConfirmation(c *gc.C) {
	err := runUpgradeSeriesCommandWithConfirmation(c, "n", machine.PrepareCommand, machineArg, seriesArg)
	c.Assert(err, gc.ErrorMatches, "upgrade series: aborted")
}

func (s *UpgradeSeriesSuite) TestUpgradeCommandShouldNotAcceptInvalidPrepCommands(c *gc.C) {
	invalidPrepCommand := "actuate"
	err := runUpgradeSeriesCommand(c, invalidPrepCommand, machineArg, seriesArg)
	c.Assert(err, gc.ErrorMatches, ".* \"actuate\" is an invalid upgrade-series command")
}

func (s *UpgradeSeriesSuite) TestUpgradeCommandShouldNotAcceptInvalidMachineArgs(c *gc.C) {
	invalidMachineArg := "machine5"
	err := runUpgradeSeriesCommand(c, machine.PrepareCommand, invalidMachineArg, seriesArg)
	c.Assert(err, gc.ErrorMatches, "\"machine5\" is an invalid machine name")
}

func (s *UpgradeSeriesSuite) TestPrepareCommandShouldOnlyAcceptSupportedSeries(c *gc.C) {
	BadSeries := "Combative Caribou"
	err := runUpgradeSeriesCommand(c, machine.PrepareCommand, machineArg, BadSeries)
	c.Assert(err, gc.ErrorMatches, ".* is an unsupported series")
}

func (s *UpgradeSeriesSuite) TestPrepareCommandShouldSupportSeriesRegardlessOfCase(c *gc.C) {
	capitalizedCaseXenial := "Xenial"
	err := runUpgradeSeriesCommand(c, machine.PrepareCommand, machineArg, capitalizedCaseXenial)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *UpgradeSeriesSuite) TestCompleteCommand(c *gc.C) {
	mockController := gomock.NewController(c)
	defer mockController.Finish()
	mockUpgradeSeriesAPI := mocks.NewMockUpgradeMachineSeriesAPI(mockController)
	mockUpgradeSeriesAPI.EXPECT().UpgradeSeriesComplete(machineArg)

	err := runUpgradeSeriesCommandGivenMock(c, mockUpgradeSeriesAPI, "y", machine.CompleteCommand, machineArg)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *UpgradeSeriesSuite) TestCompleteCommandDoesNotAcceptSeries(c *gc.C) {
	err := runUpgradeSeriesCommand(c, machine.CompleteCommand, machineArg, seriesArg)
	c.Assert(err, gc.ErrorMatches, "wrong number of arguments")
}

func (s *UpgradeSeriesSuite) TestPrepareCommandShouldPromptUserForConfirmation(c *gc.C) {
	err := s.runUpgradeSeriesCommandWithConfirmation(c, "y", machine.PrepareCommand, machineArg, seriesArg)
	c.Assert(err, jc.ErrorIsNil)
	confirmationMsg := fmt.Sprintf(machine.UpgradeSeriesConfirmationMsg, machineArg, seriesArg)
	c.Assert(s.ctx.Stdout.(*bytes.Buffer).String(), gc.Equals, confirmationMsg)
}

func (s *UpgradeSeriesSuite) TestPrepareCommandShouldAcceptAgreeAndNotPrompt(c *gc.C) {
	err := s.runUpgradeSeriesCommandWithConfirmation(c, "n", machine.PrepareCommand, machineArg, seriesArg, "--agree")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(s.ctx.Stdout.(*bytes.Buffer).String(), gc.Equals, ``)
}
