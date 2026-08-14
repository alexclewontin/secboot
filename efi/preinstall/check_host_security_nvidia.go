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
	"fmt"

	internal_efi "github.com/snapcore/secboot/internal/efi"
)

const nvidiaDGXSparkProductName = "DGX Spark"

func checkHostSecurityNVIDIA(env internal_efi.HostEnvironmentARM64) (platformFirmwareIntegrityConfig, error) {
	productName, err := env.ProductName()
	if err != nil {
		return platformFirmwareIntegrityNone, &UnsupportedPlatformError{fmt.Errorf("cannot determine platform product name: %w", err)}
	}

	switch productName {
	case nvidiaDGXSparkProductName:
		return checkHostSecurityNVIDIADGXSpark(env)
	default:
		return platformFirmwareIntegrityNone, &UnsupportedPlatformError{fmt.Errorf("unsupported NVIDIA product: %s", productName)}
	}
}

func checkHostSecurityNVIDIADGXSpark(env internal_efi.HostEnvironmentARM64) (platformFirmwareIntegrityConfig, error) {
	// TODO: Implement proper HW ROT fusing checks, once we have the documentation
	// from NVIDIA to do so. This will involve checking fuses and will return
	// platformFirmwareIntegrityVerified if set correctly.

	// TODO: Implement proper debug authentication checks, once we have the documentation
	// from NVIDIA to do so.
	return platformFirmwareIntegrityVerified, nil
}
