//go:build arm64

// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2026 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package preinstall_test

import (
	"errors"

	"github.com/canonical/go-tpm2"
	. "github.com/snapcore/secboot/efi/preinstall"
	internal_efi "github.com/snapcore/secboot/internal/efi"
	"github.com/snapcore/secboot/internal/efitest"
	"github.com/snapcore/secboot/internal/testutil"
	. "gopkg.in/check.v1"
)

type hostSecurityARM64Suite struct{}

var _ = Suite(&hostSecurityARM64Suite{})

type hostSecurityARM64ErrorEnv struct {
	systemVendor string
	product      string
}

func (e *hostSecurityARM64ErrorEnv) SystemVendor() (string, error) {
	return e.systemVendor, nil
}

func (e *hostSecurityARM64ErrorEnv) ProductName() (string, error) {
	return e.product, nil
}

func makeArm64IOMMUDevices() []internal_efi.SysfsDevice {
	return []internal_efi.SysfsDevice{
		efitest.NewMockSysfsDevice("/sys/devices/platform/soc@0/8000000.iommu", nil, "iommu", nil, nil),
	}
}

func makeArm64PCRResults(c *C) *PCRBankResults {
	return NewPCRBankResults(tpm2.HashAlgorithmSHA256, 0, [8]PcrResults{
		MakePCRResults(
			false,
			make(tpm2.Digest, 32),
			testutil.DecodeHexString(c, "a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5"),
			testutil.DecodeHexString(c, "a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5a5"),
			nil,
		),
	})
}

func (s *hostSecurityARM64Suite) TestCheckHostSecurityErrNotARM64(c *C) {
	env := efitest.NewMockHostEnvironmentWithOpts()

	_, err := CheckHostSecurity(env, nil)
	c.Check(err, ErrorMatches, `unsupported platform: cannot obtain ARM64 environment: not a ARM64 host`)
	var upe *UnsupportedPlatformError
	c.Check(errors.As(err, &upe), testutil.IsTrue)
}

func (s *hostSecurityARM64Suite) TestCheckHostSecurityErrUnknownVendor(c *C) {
	env := efitest.NewMockHostEnvironmentWithOpts(efitest.WithARM64Environment("ACME", exampleARM64ProductName))

	_, err := CheckHostSecurity(env, nil)
	c.Check(err, ErrorMatches, `unsupported platform: unsupported system vendor: ACME`)
	var upe *UnsupportedPlatformError
	c.Check(errors.As(err, &upe), testutil.IsTrue)
}

func (s *hostSecurityARM64Suite) TestCheckHostSecurityUEFIDebuggerFinding(c *C) {
	env := efitest.NewMockHostEnvironmentWithOpts(
		efitest.WithARM64Environment(exampleARM64SystemVendor, exampleARM64ProductName),
		efitest.WithSysfsDevices(makeArm64IOMMUDevices()...),
	)
	log := efitest.NewLog(c, &efitest.LogOptions{FirmwareDebugger: true})

	integrity, err := CheckHostSecurity(env, log)
	c.Check(integrity, Equals, PlatformFirmwareIntegrityMeasured)
	c.Check(err, ErrorMatches, `the platform firmware contains a debugging endpoint enabled`)
	var tmpl CompoundError
	c.Assert(err, Implements, &tmpl)
	c.Check(err.(CompoundError).Unwrap(), DeepEquals, []error{ErrUEFIDebuggingEnabled})
}

func (s *hostSecurityARM64Suite) TestCheckHostSecurityInsufficientDMAProtection(c *C) {
	env := efitest.NewMockHostEnvironmentWithOpts(
		efitest.WithARM64Environment(exampleARM64SystemVendor, exampleARM64ProductName),
		efitest.WithSysfsDevices(makeArm64IOMMUDevices()...),
	)
	log := efitest.NewLog(c, &efitest.LogOptions{DMAProtection: efitest.DMAProtectionDisabled})

	integrity, err := CheckHostSecurity(env, log)
	c.Check(integrity, Equals, PlatformFirmwareIntegrityMeasured)
	c.Check(err, ErrorMatches, `the platform firmware indicates that DMA protections are insufficient`)
	var tmpl CompoundError
	c.Assert(err, Implements, &tmpl)
	c.Check(err.(CompoundError).Unwrap(), DeepEquals, []error{ErrInsufficientDMAProtection})
}

func (s *hostSecurityARM64Suite) TestCheckHostSecurityNoIOMMU(c *C) {
	env := efitest.NewMockHostEnvironmentWithOpts(
		efitest.WithARM64Environment(exampleARM64SystemVendor, exampleARM64ProductName),
		efitest.WithSysfsDevices(),
	)
	log := efitest.NewLog(c, &efitest.LogOptions{})

	integrity, err := CheckHostSecurity(env, log)
	c.Check(integrity, Equals, PlatformFirmwareIntegrityMeasured)
	c.Check(err, ErrorMatches, `no kernel IOMMU support was detected`)
	var tmpl CompoundError
	c.Assert(err, Implements, &tmpl)
	c.Check(err.(CompoundError).Unwrap(), DeepEquals, []error{ErrNoKernelIOMMU})
}

func (s *hostSecurityARM64Suite) TestCheckHostSecurityMultipleRecoverableErrors(c *C) {
	env := efitest.NewMockHostEnvironmentWithOpts(
		efitest.WithARM64Environment(exampleARM64SystemVendor, exampleARM64ProductName),
		efitest.WithSysfsDevices(),
	)
	log := efitest.NewLog(c, &efitest.LogOptions{FirmwareDebugger: true})

	integrity, err := CheckHostSecurity(env, log)
	c.Check(integrity, Equals, PlatformFirmwareIntegrityMeasured)
	c.Check(err, ErrorMatches, `2 errors detected:
- the platform firmware contains a debugging endpoint enabled
- no kernel IOMMU support was detected
`)
	var tmpl CompoundError
	c.Assert(err, Implements, &tmpl)
	c.Check(err.(CompoundError).Unwrap(), DeepEquals, []error{ErrUEFIDebuggingEnabled, ErrNoKernelIOMMU})
}

func (s *hostSecurityARM64Suite) TestCheckDiscreteTPMPartialResetAttackMitigationStatusUnknownForUnsupportedVendor(c *C) {
	env := efitest.NewMockHostEnvironmentWithOpts(
		efitest.WithARM64Environment(exampleARM64SystemVendor, exampleARM64ProductName),
		efitest.WithSysfsDevices(makeArm64TPMDevice("tpm_crb")),
	)

	status, err := CheckDiscreteTPMPartialResetAttackMitigationStatus(env, makeArm64PCRResults(c))
	c.Check(status, Equals, DtpmPartialResetAttackMitigationUnknown)
	c.Check(err, ErrorMatches, `error with TPM2 device: unsupported platform: unsupported system vendor: `+exampleARM64SystemVendor)
	var tpmErr *TPM2DeviceError
	c.Check(errors.As(err, &tpmErr), testutil.IsTrue)
	c.Check(status, Equals, DtpmPartialResetAttackMitigationUnknown)
}

func (s *hostSecurityARM64Suite) TestCheckDiscreteTPMPartialResetAttackMitigationStatusNotRequiredForOPTEE(c *C) {
	env := efitest.NewMockHostEnvironmentWithOpts(
		efitest.WithARM64Environment(exampleARM64SystemVendor, exampleARM64ProductName),
		efitest.WithSysfsDevices(makeArm64TPMDevice("optee-ftpm")),
	)

	status, err := CheckDiscreteTPMPartialResetAttackMitigationStatus(env, makeArm64PCRResults(c))
	c.Check(err, IsNil)
	c.Check(status, Equals, DtpmPartialResetAttackMitigationNotRequired)
}
