//go:build amd64

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

	"github.com/canonical/cpuid"
	"github.com/canonical/go-tpm2"
	tpm2_testutil "github.com/canonical/go-tpm2/testutil"
	. "github.com/snapcore/secboot/efi/preinstall"
	internal_efi "github.com/snapcore/secboot/internal/efi"
	"github.com/snapcore/secboot/internal/efitest"
	"github.com/snapcore/secboot/internal/testutil"
	"github.com/snapcore/secboot/internal/tpm2_device"
	. "gopkg.in/check.v1"
)

func (s *runChecksContextSuite) TestRunBadHostSecurityErrorMissingIntelMEModule(c *C) {
	// Test case where host security checks fail because the intel ME kernel module is missing.
	devices := []internal_efi.SysfsDevice{
		efitest.NewMockSysfsDevice("/sys/devices/virtual/iommu/dmar0", nil, "iommu", nil, nil),
		efitest.NewMockSysfsDevice("/sys/devices/virtual/iommu/dmar1", nil, "iommu", nil, nil),
		efitest.NewMockSysfsDevice("/sys/devices/pci0000:00:16:0", map[string]string{"PCI_CLASS": "78000", "PCI_ID": "8086:7E70"}, "pci", nil, nil),
	}

	errs := s.testRun(c, &testRunChecksContextRunParams{
		env: efitest.NewMockHostEnvironmentWithOpts(
			efitest.WithVirtMode(internal_efi.VirtModeNone, internal_efi.DetectVirtModeAll),
			efitest.WithTPMDevice(newTpmDevice(tpm2_testutil.NewTransportBackedDevice(s.Transport, false, 1), nil, tpm2_device.ErrNoPPI)),
			efitest.WithLog(efitest.NewLog(c, &efitest.LogOptions{Algorithms: []tpm2.HashAlgorithmId{tpm2.HashAlgorithmSHA256}})),
			efitest.WithAMD64Environment("GenuineIntel", 0x6, []uint64{cpuid.SDBG, cpuid.SMX}, 4, map[uint32]uint64{0x13a: (3 << 1), 0xc80: 0x40000000}),
			efitest.WithSysfsDevices(devices...),
			efitest.WithMockVars(efitest.MockVars{}.SetSecureBoot(false)),
		),
		tpmPropertyModifiers: map[tpm2.Property]uint32{
			tpm2.PropertyNVCountersMax:     0,
			tpm2.PropertyPSFamilyIndicator: 1,
			tpm2.PropertyManufacturer:      uint32(tpm2.TPMManufacturerINTC),
		},
		enabledBanks: []tpm2.HashAlgorithmId{tpm2.HashAlgorithmSHA256},
		actions:      []actionAndArgs{{action: ActionNone}},
	})
	c.Assert(errs, HasLen, 1)
	c.Check(errs[0], ErrorMatches, `error with system security: encountered an error when checking Intel BootGuard configuration: the kernel module "mei_me" must be loaded`)
	c.Check(errs[0], DeepEquals, NewWithKindAndActionsError(ErrorKindInternal, nil, nil, errs[0].Unwrap()))
	c.Check(errors.Is(errs[0], MissingKernelModuleError("mei_me")), testutil.IsTrue)
}

func (s *runChecksContextSuite) TestRunBadHostSecurityErrorMissingMSR(c *C) {
	// Test case where host security checks fail because the MSR kernel module is missing.
	meiAttrs := map[string][]byte{
		"fw_ver": []byte(`0:16.1.27.2176
0:16.1.27.2176
0:16.0.15.1624
`),
		"fw_status": fwStatusBase,
	}
	devices := []internal_efi.SysfsDevice{
		efitest.NewMockSysfsDevice("/sys/devices/virtual/iommu/dmar0", nil, "iommu", nil, nil),
		efitest.NewMockSysfsDevice("/sys/devices/virtual/iommu/dmar1", nil, "iommu", nil, nil),
		efitest.NewMockSysfsDevice("/sys/devices/pci0000:00/0000:00:16.0/mei/mei0", map[string]string{"DEVNAME": "mei0"}, "mei", meiAttrs, efitest.NewMockSysfsDevice(
			"/sys/devices/pci0000:00:16:0", map[string]string{"DRIVER": "mei_me"}, "pci", nil, nil,
		)),
	}

	errs := s.testRun(c, &testRunChecksContextRunParams{
		env: efitest.NewMockHostEnvironmentWithOpts(
			efitest.WithVirtMode(internal_efi.VirtModeNone, internal_efi.DetectVirtModeAll),
			efitest.WithTPMDevice(newTpmDevice(tpm2_testutil.NewTransportBackedDevice(s.Transport, false, 1), nil, tpm2_device.ErrNoPPI)),
			efitest.WithLog(efitest.NewLog(c, &efitest.LogOptions{Algorithms: []tpm2.HashAlgorithmId{tpm2.HashAlgorithmSHA256}})),
			efitest.WithAMD64Environment("GenuineIntel", 0x6, []uint64{cpuid.SDBG, cpuid.SMX}, 0, nil),
			efitest.WithSysfsDevices(devices...),
			efitest.WithMockVars(efitest.MockVars{}.SetSecureBoot(false)),
		),
		tpmPropertyModifiers: map[tpm2.Property]uint32{
			tpm2.PropertyNVCountersMax:     0,
			tpm2.PropertyPSFamilyIndicator: 1,
			tpm2.PropertyManufacturer:      uint32(tpm2.TPMManufacturerINTC),
		},
		enabledBanks: []tpm2.HashAlgorithmId{tpm2.HashAlgorithmSHA256},
		actions:      []actionAndArgs{{action: ActionNone}},
	})
	c.Assert(errs, HasLen, 1)
	c.Check(errs[0], ErrorMatches, `error with system security: encountered an error when checking Intel CPU debugging configuration: the kernel module "msr" must be loaded`)
	c.Check(errs[0], DeepEquals, NewWithKindAndActionsError(ErrorKindInternal, nil, nil, errs[0].Unwrap()))
	c.Check(errors.Is(errs[0], MissingKernelModuleError("msr")), testutil.IsTrue)
}
