package preinstall_test

import "strings"

// Mocked values for the MEI HFSTS registers
// (usually exposed as: /sys/devices/pci0000:00/0000:00:16.0/mei/mei0/fw_status)
// Format: HFSTS1..6 separated by NL (\n) characters
var (
	fwStatusFPFNotLockedCSME11          = toFwStatus("94000245 09F10506 00000020 00004000 00041F03 87E003CB")
	fwStatusNoManufLockCSME11           = toFwStatus("94000245 09F10506 00000020 00004000 00041F03 C7C003CB")
	fwStatusUnsupportedNoFVMECSME11     = toFwStatus("94000245 09F10506 00000020 00004000 00041F03 C7E00002")
	fwStatusInvalidProfileCSME11        = toFwStatus("94000245 09F10506 00000020 00004000 00041F03 C7E0024A")
	fwStatusFVECSME11                   = toFwStatus("94000245 09F10506 00000020 00004000 00041F03 C7E002CB")
	fwStatusUnsupportedVMProfileCSME11  = toFwStatus("94000245 09F10506 00000020 00004000 00041F03 C7E0030A")
	fwStatusBase                        = toFwStatus("94000245 09F10506 00000020 00004000 00041F03 C7E003CB")
	fwStatusBootGuardDisable            = toFwStatus("94000245 09F10506 00000020 00004000 00041F03 D7E003CB")
	fwStatusManufacturingMode           = toFwStatus("94000255 09F10506 00000020 00004000 00041F03 C7E003CB")
	fwStatusOperationModeDebug          = toFwStatus("94020245 09F10506 00000020 00004000 00041F03 C7E003CB")
	fwStatusOperationModeOverrideJumper = toFwStatus("94040245 09F10506 00000020 00004000 00041F03 C7E003CB")
	fwStatusNoFVMEProfileCSME18         = toFwStatus("A4000245 09110500 00000020 00000000 02E21F03 40200000")
	fwStatusUnsupportedVMProfileCSME18  = toFwStatus("A4000245 09110500 00000020 00000000 02EE1F03 40200000")
	fwStatusFVECSME18                   = toFwStatus("A4000245 09110500 00000020 00000000 02F21F03 40200000")
	fwStatusACMNotActive                = toFwStatus("A4000245 09110500 00000020 00000000 02F61E02 40200000")
	fwStatusACMNotDone                  = toFwStatus("A4000245 09110500 00000020 00000000 02F61E03 40200000")
	fwStatusInvalidProfileCSME18        = toFwStatus("A4000245 09110500 00000020 00000000 02F61F01 40200000")
	fwStatusFPFNotLockedCSME18          = toFwStatus("A4000245 09110500 00000020 00000000 02F61F03 00200000")
	fwStatusNoManufLockCSME18           = toFwStatus("A4000245 09110500 00000020 00000000 02F61F03 40000000")
	fwStatusFVMECSME18                  = toFwStatus("A4000245 09110500 00000020 00000000 02F61F03 40200000")
	fwStatusNoSPIProtectionCSME18       = toFwStatus("A4000255 09110500 00000020 00000000 02F61F03 40200000")

	// Malformed registers
	fwStatusRegular        = toFwStatus("94000245 09F10506 00000020 00004000 00041F03 C7E0034B")
	fwStatusTooMany        = toFwStatus("94000245 09F10506 00000020 00004000 00041F03 C7E003CB 00000000")
	fwStatusInvalidLineLen = toFwStatus("94000245 09F10506 00000020 00004000 00041F03 7E003CB")
	fwStatusInvalidLine    = toFwStatus("94000245 09F10506 00000020 00004000 00041F03 G7E003CB")
	fwStatusNotEnough      = toFwStatus("94000245 09F10506 00000020 00004000 00041F03")
)

// toFwStatus converts to the format used by the Kernel to expose fw_status
func toFwStatus(str string) []byte {
	bytes := []byte(strings.Replace(str, " ", "\n", -1) + "\n")
	return bytes
}
