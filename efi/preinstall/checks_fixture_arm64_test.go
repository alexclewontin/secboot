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
	. "github.com/snapcore/secboot/efi/preinstall"
	internal_efi "github.com/snapcore/secboot/internal/efi"
	"github.com/snapcore/secboot/internal/efitest"
)

const (
	exampleARM64SystemVendor = "Example Vendor"
	exampleARM64ProductName  = "Example OP-TEE Platform"
)

func init() {
	RegisterARM64TestPlatform(exampleARM64SystemVendor, exampleARM64ProductName)
}

func runChecksArm64TPMDevice(driver string) internal_efi.SysfsDevice {
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

func runChecksPlatformHostFixtures() []runChecksHostFixture {
	newDevices := func(tpmDriver string, withIOMMU bool) []internal_efi.SysfsDevice {
		devices := []internal_efi.SysfsDevice{runChecksArm64TPMDevice(tpmDriver)}
		if withIOMMU {
			devices = append([]internal_efi.SysfsDevice{efitest.NewMockSysfsDevice("/sys/devices/platform/smmu0", nil, "iommu", nil, nil)}, devices...)
		}
		return devices
	}

	return []runChecksHostFixture{
		{
			name: "example-arm64-optee-ftpm",
			capabilities: runChecksHostCapabilityValid |
				runChecksHostCapabilityNotVirtualMachine |
				runChecksHostCapabilityFirmwareTPM,
			environment:             efitest.WithARM64Environment(exampleARM64SystemVendor, exampleARM64ProductName),
			virtualizationMode:      internal_efi.VirtModeNone,
			virtualizationDetection: internal_efi.DetectVirtModeAll,
			devices:                 newDevices("optee-ftpm", true),
			additionalExpectedFlags: RequireLockToPlatformFirmware,
		},
		{
			name:                    "example-arm64-optee-ftpm-no-kernel-iommu",
			capabilities:            runChecksHostCapabilityNotVirtualMachine | runChecksHostCapabilityNoKernelIOMMU | runChecksHostCapabilityFirmwareTPM,
			environment:             efitest.WithARM64Environment(exampleARM64SystemVendor, exampleARM64ProductName),
			virtualizationMode:      internal_efi.VirtModeNone,
			virtualizationDetection: internal_efi.DetectVirtModeAll,
			devices:                 newDevices("optee-ftpm", false),
			additionalExpectedFlags: RequireLockToPlatformFirmware,
		},
		{
			name:                    "virtual-machine",
			capabilities:            runChecksHostCapabilityVirtualMachine,
			environment:             func(*efitest.MockHostEnvironment) {},
			virtualizationMode:      "qemu",
			virtualizationDetection: internal_efi.DetectVirtModeVM,
		},
	}
}
