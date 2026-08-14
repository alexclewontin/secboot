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

	"github.com/pilebones/go-udev/netlink"
	internal_efi "github.com/snapcore/secboot/internal/efi"
)

// isTPMDiscrete determines whether the default TPM is discrete. OP-TEE firmware
// TPMs are identified by their backing kernel driver. Other implementations use
// platform-specific knowledge.
func isTPMDiscrete(env internal_efi.HostEnvironment) (bool, error) {
	isOpteefTPM, err := isTPMFirmwareOptee(env)
	if err != nil {
		return false, err
	}
	if isOpteefTPM {
		return false, nil
	}

	arm64Env, err := env.ARM64()
	if err != nil {
		return false, err
	}

	systemVendor, err := arm64Env.SystemVendor()
	if err != nil {
		return false, &UnsupportedPlatformError{fmt.Errorf("cannot determine system vendor: %w", err)}
	}

	switch systemVendor {
	case "NVIDIA":
		return isTPMDiscreteNvidia(arm64Env)
	default:
		return false, &UnsupportedPlatformError{fmt.Errorf("unsupported system vendor: %s", systemVendor)}
	}
}

// isTPMFirmwareOptee determines whether the default TPM is an OP-TEE firmware TPM,
// as identified by its backing kernel driver
func isTPMFirmwareOptee(env internal_efi.HostEnvironment) (bool, error) {
	devices, err := env.EnumerateDevices(&netlink.RuleDefinition{
		Env: map[string]string{
			"SUBSYSTEM": "tpm",
			"DEVNAME":   "tpm0",
		},
	})
	if err != nil {
		return false, fmt.Errorf("cannot enumerate TPM devices: %w", err)
	}
	if len(devices) != 1 {
		return false, fmt.Errorf("internal error: expected one tpm0 device, found %d", len(devices))
	}

	parent, err := devices[0].Parent()
	if err != nil {
		return false, fmt.Errorf("cannot obtain parent of tpm0 device: %w", err)
	}
	if parent == nil {
		return false, fmt.Errorf("internal error: tpm0 device has no parent")
	}

	switch parent.Properties()["DRIVER"] {
	case "optee-ftpm", "ftpm-tee":
		return true, nil
	default:
		return false, nil
	}
}

// isTPMDiscreteNvidia determines whether the default TPM is discrete on NVIDIA systems
func isTPMDiscreteNvidia(env internal_efi.HostEnvironmentARM64) (bool, error) {
	productName, err := env.ProductName()
	if err != nil {
		return false, &UnsupportedPlatformError{fmt.Errorf("cannot determine product name: %w", err)}
	}

	switch productName {
	// We just happen to know that the NVIDIA DGX Spark has a dTPM
	case nvidiaDGXSparkProductName:
		return true, nil
	default:
		return false, &UnsupportedPlatformError{fmt.Errorf("unsupported NVIDIA product: %s", productName)}
	}
}
