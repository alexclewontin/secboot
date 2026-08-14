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

package preinstall

import (
	"errors"
	"fmt"

	"github.com/canonical/tcglog-parser"
	internal_efi "github.com/snapcore/secboot/internal/efi"
)

// checkHostSecurityARM64Platform selects the platform-specific firmware
// integrity check. Tests replace this to supply synthetic platforms.
var checkHostSecurityARM64Platform = func(env internal_efi.HostEnvironmentARM64, systemVendor string) (platformFirmwareIntegrityConfig, error) {
	switch systemVendor {
	case "NVIDIA":
		return checkHostSecurityNVIDIA(env)
	default:
		return platformFirmwareIntegrityNone, &UnsupportedPlatformError{fmt.Errorf("unsupported system vendor: %s", systemVendor)}
	}
}

// checkHostSecurity is the main entry point for verifying that the host security
// is sufficient. Errors that can't be resolved or which should prevent further checks from running
// are returned immediately and without any wrapping. Errors that can be resolved and which shouldn't
// prevent further checks from running are returned wrapped in [joinError].
func checkHostSecurity(env internal_efi.HostEnvironment, log *tcglog.Log) (platformFirmwareIntegrityConfig, error) {
	arm64Env, err := env.ARM64()
	if err != nil {
		return platformFirmwareIntegrityNone, &UnsupportedPlatformError{fmt.Errorf("cannot obtain ARM64 environment: %w", err)}
	}

	systemVendor, err := arm64Env.SystemVendor()
	if err != nil {
		return platformFirmwareIntegrityNone, &UnsupportedPlatformError{fmt.Errorf("cannot determine system vendor: %w", err)}
	}

	integrity, err := checkHostSecurityARM64Platform(arm64Env, systemVendor)
	if err != nil {
		return platformFirmwareIntegrityNone, err
	}

	return checkHostSecurityARM64Generic(env, log, integrity)
}

func checkHostSecurityARM64Generic(env internal_efi.HostEnvironment, log *tcglog.Log, integrity platformFirmwareIntegrityConfig) (platformFirmwareIntegrityConfig, error) {
	var errs []error

	if err := checkSecureBootPolicyPCRForDegradedFirmwareSettings(log); err != nil {
		var ce CompoundError
		if !errors.As(err, &ce) {
			return platformFirmwareIntegrityNone, fmt.Errorf("encountered an error whilst checking the TCG log for degraded firmware settings: %w", err)
		}
		errs = append(errs, ce.Unwrap()...)
	}

	if err := checkForKernelIOMMU(env); err != nil {
		switch {
		case errors.Is(err, ErrNoKernelIOMMU):
			errs = append(errs, err)
		default:
			return platformFirmwareIntegrityNone, fmt.Errorf("encountered an error whilst checking sysfs to determine that kernel IOMMU support is enabled: %w", err)
		}
	}

	if len(errs) > 0 {
		return integrity, joinErrors(errs...)
	}

	return integrity, nil
}

// checkDiscreteTPMPartialResetAttackMitigationStatus determines whether a partial mitigation
// against discrete TPM reset attacks should be enabled.
func checkDiscreteTPMPartialResetAttackMitigationStatus(env internal_efi.HostEnvironment, _ *pcrBankResults) (discreteTPMPartialResetAttackMitigationStatus, error) {
	discreteTPM, err := isTPMDiscrete(env)
	if err != nil {
		return dtpmPartialResetAttackMitigationUnknown, &TPM2DeviceError{err}
	}
	if !discreteTPM {
		return dtpmPartialResetAttackMitigationNotRequired, nil
	}

	// ARM64 has no generic mechanism to establish that the TPM startup locality
	// is protected by the hardware root of trust, so PCR0 binding cannot be relied
	// on to mitigate an independent reset of a discrete TPM.
	return dtpmPartialResetAttackMitigationUnavailable, nil
}
