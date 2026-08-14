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
	"github.com/canonical/cpuid"
	internal_efi "github.com/snapcore/secboot/internal/efi"
	"github.com/snapcore/secboot/internal/efitest"
)

func runChecksPlatformHostFixtures() []runChecksHostFixture {
	intelDevices := func(status []byte, withIOMMU bool) []internal_efi.SysfsDevice {
		attrs := map[string][]byte{
			"fw_ver": []byte(`0:16.1.27.2176
0:16.1.27.2176
0:16.0.15.1624
`),
			"fw_status": status,
		}
		var devices []internal_efi.SysfsDevice
		if withIOMMU {
			devices = append(devices,
				efitest.NewMockSysfsDevice("/sys/devices/virtual/iommu/dmar0", nil, "iommu", nil, nil),
				efitest.NewMockSysfsDevice("/sys/devices/virtual/iommu/dmar1", nil, "iommu", nil, nil),
			)
		}
		return append(devices, efitest.NewMockSysfsDevice("/sys/devices/pci0000:00/0000:00:16.0/mei/mei0", map[string]string{"DEVNAME": "mei0"}, "mei", attrs, efitest.NewMockSysfsDevice(
			"/sys/devices/pci0000:00:16:0", map[string]string{"DRIVER": "mei_me"}, "pci", nil, nil,
		)))
	}

	return []runChecksHostFixture{
		{
			name: "intel-ptt",
			capabilities: runChecksHostCapabilityValid |
				runChecksHostCapabilityNotVirtualMachine |
				runChecksHostCapabilityFirmwareTPM,
			environment:             efitest.WithAMD64Environment("GenuineIntel", 0x6, []uint64{cpuid.SDBG, cpuid.SMX}, 4, map[uint32]uint64{0xc80: 0x40000000, 0x13a: (3 << 1)}),
			virtualizationMode:      internal_efi.VirtModeNone,
			virtualizationDetection: internal_efi.DetectVirtModeAll,
			devices:                 intelDevices(fwStatusBase, true),
		},
		{
			name: "intel-dtpm-smx",
			capabilities: runChecksHostCapabilityValid |
				runChecksHostCapabilityNotVirtualMachine |
				runChecksHostCapabilityDiscreteTPM |
				runChecksHostCapabilityStartupLocality0AccessibleFromOS |
				runChecksHostCapabilityStartupLocality3InaccessibleFromOS |
				runChecksHostCapabilityStartupLocality4InaccessibleFromOS,
			environment:             efitest.WithAMD64Environment("GenuineIntel", 0x6, []uint64{cpuid.SDBG, cpuid.SMX}, 4, map[uint32]uint64{0xc80: 0x40000000, 0x13a: (2 << 1)}),
			virtualizationMode:      internal_efi.VirtModeNone,
			virtualizationDetection: internal_efi.DetectVirtModeAll,
			devices:                 intelDevices(fwStatusBase, true),
		},
		{
			name: "intel-dtpm-no-smx",
			capabilities: runChecksHostCapabilityValid |
				runChecksHostCapabilityNotVirtualMachine |
				runChecksHostCapabilityDiscreteTPM |
				runChecksHostCapabilityStartupLocality0AccessibleFromOS |
				runChecksHostCapabilityStartupLocality3AccessibleFromOS |
				runChecksHostCapabilityStartupLocality4AccessibleFromOS,
			environment:             efitest.WithAMD64Environment("GenuineIntel", 0x6, []uint64{cpuid.SDBG}, 4, map[uint32]uint64{0xc80: 0x40000000, 0x13a: (2 << 1)}),
			virtualizationMode:      internal_efi.VirtModeNone,
			virtualizationDetection: internal_efi.DetectVirtModeAll,
			devices:                 intelDevices(fwStatusBase, true),
		},
		{
			name:                    "intel-ptt-no-kernel-iommu",
			capabilities:            runChecksHostCapabilityNotVirtualMachine | runChecksHostCapabilityNoKernelIOMMU | runChecksHostCapabilityFirmwareTPM,
			environment:             efitest.WithAMD64Environment("GenuineIntel", 0x6, []uint64{cpuid.SDBG, cpuid.SMX}, 4, map[uint32]uint64{0xc80: 0x40000000, 0x13a: (3 << 1)}),
			virtualizationMode:      internal_efi.VirtModeNone,
			virtualizationDetection: internal_efi.DetectVirtModeAll,
			devices:                 intelDevices(fwStatusBase, false),
		},
		{
			name:                    "intel-mfg-mode",
			capabilities:            runChecksHostCapabilityNotVirtualMachine | runChecksHostCapabilityInsufficientHWRootOfTrust,
			environment:             efitest.WithAMD64Environment("GenuineIntel", 0x6, []uint64{cpuid.SDBG, cpuid.SMX}, 4, map[uint32]uint64{0xc80: 0x40000000, 0x13a: (3 << 1)}),
			virtualizationMode:      internal_efi.VirtModeNone,
			virtualizationDetection: internal_efi.DetectVirtModeAll,
			devices:                 intelDevices(fwStatusManufacturingMode, true),
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
