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

	. "github.com/snapcore/secboot/efi/preinstall"
	internal_efi "github.com/snapcore/secboot/internal/efi"
	"github.com/snapcore/secboot/internal/efitest"
	"github.com/snapcore/secboot/internal/testutil"
	. "gopkg.in/check.v1"
)

type tpmARM64Suite struct{}

var _ = Suite(&tpmARM64Suite{})

func makeArm64TPMDevice(driver string) internal_efi.SysfsDevice {
	parent := efitest.NewMockSysfsDevice(
		"/sys/devices/platform/firmware-tpm",
		map[string]string{"DRIVER": driver},
		"platform",
		nil,
		nil,
	)
	return efitest.NewMockSysfsDevice(
		"/sys/devices/platform/firmware-tpm/tpm/tpm0",
		map[string]string{"DEVNAME": "tpm0"},
		"tpm",
		nil,
		parent,
	)
}

func (s *tpmARM64Suite) TestIsTPMDiscreteOPTEE(c *C) {
	for _, driver := range []string{"optee-ftpm", "ftpm-tee"} {
		env := efitest.NewMockHostEnvironmentWithOpts(efitest.WithSysfsDevices(makeArm64TPMDevice(driver)))

		discrete, err := IsTPMDiscrete(env)
		c.Check(err, IsNil, Commentf("driver %q", driver))
		c.Check(discrete, testutil.IsFalse, Commentf("driver %q", driver))
	}
}

func (s *tpmARM64Suite) TestIsTPMDiscreteDGXSpark(c *C) {
	env := efitest.NewMockHostEnvironmentWithOpts(
		efitest.WithARM64Environment("NVIDIA", "DGX Spark"),
		efitest.WithSysfsDevices(makeArm64TPMDevice("tpm_crb")),
	)

	discrete, err := IsTPMDiscrete(env)
	c.Check(err, IsNil)
	c.Check(discrete, testutil.IsTrue)
}

func (s *tpmARM64Suite) TestIsTPMDiscreteUnsupportedSystemVendor(c *C) {
	env := efitest.NewMockHostEnvironmentWithOpts(
		efitest.WithARM64Environment("ACME", "Unknown"),
		efitest.WithSysfsDevices(makeArm64TPMDevice("tpm_crb")),
	)

	_, err := IsTPMDiscrete(env)
	c.Check(err, ErrorMatches, `unsupported platform: unsupported system vendor: ACME`)
	var upe *UnsupportedPlatformError
	c.Check(errors.As(err, &upe), testutil.IsTrue)
}
