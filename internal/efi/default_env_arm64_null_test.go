//go:build !arm64

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
	. "github.com/snapcore/secboot/internal/efi"
	. "gopkg.in/check.v1"
)

type defaultEnvARM64Suite struct{}

var _ = Suite(&defaultEnvARM64Suite{})

func (s *defaultEnvARM64Suite) TestNotARM64Host(c *C) {
	_, err := DefaultEnv.ARM64()
	c.Check(err, Equals, ErrNotARM64Host)
}
