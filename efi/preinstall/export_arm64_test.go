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

type PCRBankResults = pcrBankResults

var DtpmPartialResetAttackMitigationUnknown = dtpmPartialResetAttackMitigationUnknown

func RegisterARM64TestPlatform(systemVendor, productName string) {
	previous := checkHostSecurityARM64Platform
	checkHostSecurityARM64Platform = func(env internal_efi.HostEnvironmentARM64, vendor string) (platformFirmwareIntegrityConfig, error) {
		if vendor != systemVendor {
			return previous(env, vendor)
		}

		product, err := env.ProductName()
		if err != nil {
			return platformFirmwareIntegrityNone, &UnsupportedPlatformError{fmt.Errorf("cannot determine platform product name: %w", err)}
		}
		if product != productName {
			return previous(env, vendor)
		}

		return platformFirmwareIntegrityMeasured, nil
	}
}
