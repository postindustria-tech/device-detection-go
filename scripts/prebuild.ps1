/*
This powershell script builds the device-detection-cxx binaries using CMake as
described in device-detection-cxx README.md.

The resulting binaries can be kept at the same location and be linked to the
go package using Cgo LDFLAGS.

This powershell script is used on both Windows and Linux.
*/