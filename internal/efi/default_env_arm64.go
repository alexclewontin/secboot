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

package efi

import (
	"fmt"
	"strings"
)

var (
	dmiSystemVendorPath = "/sys/class/dmi/id/sys_vendor"
	dmiProductNamePath  = "/sys/class/dmi/id/product_name"
)

type defaultEnvARM64Impl struct{}

func readSMBIOSIdentity(path string) (string, error) {
	data, err := osReadFile(path)
	if err != nil {
		return "", fmt.Errorf("cannot read %s: %w", path, err)
	}

	value := strings.TrimSpace(string(data))
	if value == "" {
		return "", fmt.Errorf("%s is empty", path)
	}

	return value, nil
}

// SystemVendor implements [HostEnvironmentARM64.SystemVendor].
func (defaultEnvARM64Impl) SystemVendor() (string, error) {
	return readSMBIOSIdentity(dmiSystemVendorPath)
}

// ProductName implements [HostEnvironmentARM64.ProductName].
func (defaultEnvARM64Impl) ProductName() (string, error) {
	return readSMBIOSIdentity(dmiProductNamePath)
}

// ARM64 implements [HostEnvironment.ARM64].
func (defaultEnvImpl) ARM64() (HostEnvironmentARM64, error) {
	return defaultEnvARM64Impl{}, nil
}
