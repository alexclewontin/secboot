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
	. "gopkg.in/check.v1"
)

type runChecksHostCapabilities uint32

const (
	runChecksHostCapabilityValid runChecksHostCapabilities = 1 << iota

	// Platform/runtime topology.
	runChecksHostCapabilityVirtualMachine
	runChecksHostCapabilityNotVirtualMachine
	runChecksHostCapabilityNoKernelIOMMU
	runChecksHostCapabilityFirmwareTPM
	runChecksHostCapabilityDiscreteTPM

	// Hardware security properties.
	runChecksHostCapabilityInsufficientHWRootOfTrust

	// Discrete-TPM locality properties.
	runChecksHostCapabilityStartupLocality0AccessibleFromOS
	runChecksHostCapabilityStartupLocality3InaccessibleFromOS
	runChecksHostCapabilityStartupLocality3AccessibleFromOS
	runChecksHostCapabilityStartupLocality4InaccessibleFromOS
	runChecksHostCapabilityStartupLocality4AccessibleFromOS
)

// Platform identity capabilities do not imply security properties. In
// particular, HasIntelBootGuard does not imply that any TPM startup locality is
// inaccessible from the OS; fixtures must declare locality properties
// separately.

// runChecksHostFixture abstracts host-specific details from scenario logic in
// the shared RunChecks and RunChecksContext tests.
type runChecksHostFixture struct {
	name                    string
	capabilities            runChecksHostCapabilities
	environment             efitest.MockHostEnvironmentOption
	virtualizationMode      string
	virtualizationDetection internal_efi.DetectVirtMode
	devices                 []internal_efi.SysfsDevice

	// additionalExpectedFlags allows a host fixture to inject fixture-specific
	// CheckResultFlags to the expected value. For example, platforms that validate
	// platformFirmwareIntegrity via measured boot instead of reading a fused key will
	// report RequireLockToPlatformFirmware, which needs to be OR'd against the
	// test-specific expected flags.
	additionalExpectedFlags CheckResultFlags
}

func runChecksHostFixturesFor(c *C, required runChecksHostCapabilities) []runChecksHostFixture {
	var matches []runChecksHostFixture
	for _, fixture := range runChecksPlatformHostFixtures() {
		if fixture.capabilities&required == required {
			matches = append(matches, fixture)
		}
	}
	if len(matches) == 0 {
		c.Skip("no platform fixture provides the required host capabilities")
	}
	return matches
}

func (f *runChecksHostFixture) newEnvironment(options ...efitest.MockHostEnvironmentOption) internal_efi.HostEnvironment {
	base := []efitest.MockHostEnvironmentOption{
		f.environment,
		efitest.WithVirtMode(f.virtualizationMode, f.virtualizationDetection),
	}
	devices := append([]internal_efi.SysfsDevice(nil), f.devices...)
	base = append(base, efitest.WithSysfsDevices(devices...))
	return efitest.NewMockHostEnvironmentWithOpts(append(base, options...)...)
}
