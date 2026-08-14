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

package efi_test

import (
	"errors"
	"fmt"

	. "github.com/snapcore/secboot/internal/efi"
	"github.com/snapcore/secboot/internal/testutil"
	. "gopkg.in/check.v1"
)

type defaultEnvARM64Suite struct{}

var _ = Suite(&defaultEnvARM64Suite{})

func (s *defaultEnvARM64Suite) TestSystemVendor(c *C) {
	restore := MockOsReadFile(func(path string) ([]byte, error) {
		switch path {
		case "/sys/class/dmi/id/sys_vendor":
			return []byte("\tExample Vendor\n"), nil
		case "/sys/class/dmi/id/product_name":
			return []byte("Example Product\n"), nil
		default:
			return nil, fmt.Errorf("unexpected path: %s", path)
		}
	})
	defer restore()

	arm64, err := DefaultEnv.ARM64()
	c.Assert(err, IsNil)

	vendor, err := arm64.SystemVendor()
	c.Check(err, IsNil)
	c.Check(vendor, Equals, "Example Vendor")
}

func (s *defaultEnvARM64Suite) TestSystemVendorReadError(c *C) {
	restore := MockOsReadFile(func(path string) ([]byte, error) {
		if path == "/sys/class/dmi/id/sys_vendor" {
			return nil, errors.New("some error")
		}
		return []byte("Example Product\n"), nil
	})
	defer restore()

	arm64, err := DefaultEnv.ARM64()
	c.Assert(err, IsNil)

	_, err = arm64.SystemVendor()
	c.Check(err, ErrorMatches, "cannot read /sys/class/dmi/id/sys_vendor: some error")
	c.Check(errors.Is(err, errors.New("some error")), testutil.IsFalse)
}

func (s *defaultEnvARM64Suite) TestProductName(c *C) {
	restore := MockOsReadFile(func(path string) ([]byte, error) {
		switch path {
		case "/sys/class/dmi/id/sys_vendor":
			return []byte("Example Vendor\n"), nil
		case "/sys/class/dmi/id/product_name":
			return []byte("\tExample Product\n"), nil
		default:
			return nil, fmt.Errorf("unexpected path: %s", path)
		}
	})
	defer restore()

	arm64, err := DefaultEnv.ARM64()
	c.Assert(err, IsNil)

	productName, err := arm64.ProductName()
	c.Check(err, IsNil)
	c.Check(productName, Equals, "Example Product")
}

func (s *defaultEnvARM64Suite) TestProductNameReadError(c *C) {
	restore := MockOsReadFile(func(path string) ([]byte, error) {
		if path == "/sys/class/dmi/id/product_name" {
			return nil, errors.New("some error")
		}
		return []byte("Example Vendor\n"), nil
	})
	defer restore()

	arm64, err := DefaultEnv.ARM64()
	c.Assert(err, IsNil)

	_, err = arm64.ProductName()
	c.Check(err, ErrorMatches, "cannot read /sys/class/dmi/id/product_name: some error")
	c.Check(errors.Is(err, errors.New("some error")), testutil.IsFalse)
}
