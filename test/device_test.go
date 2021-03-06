package test

import (
	"cryptsetup"
	"testing"
)

func Test_Device_Init_Works_If_Device_Is_Found(test *testing.T) {
	testWrapper := TestWrapper{test}

	device, err := cryptsetup.Init(DevicePath)
	testWrapper.AssertNoError(err)

	if device.Type() != "" {
		test.Error("Device should have no type.")
	}
}

func Test_Device_Init_Fails_If_Device_Is_Not_Found(test *testing.T) {
	testWrapper := TestWrapper{test}

	_, err := cryptsetup.Init("nonExistingDevicePath")
	testWrapper.AssertError(err)
	testWrapper.AssertErrorCodeEquals(err, -15)
}

func Test_Device_Free_Works(test *testing.T) {
	testWrapper := TestWrapper{test}

	device, err := cryptsetup.Init(DevicePath)
	testWrapper.AssertNoError(err)

	err = device.Format(
		cryptsetup.LUKS1{Hash: "sha256"},
		cryptsetup.GenericParams{Cipher: "aes", CipherMode: "xts-plain64", VolumeKeySize: 512 / 8},
	)
	testWrapper.AssertNoError(err)

	code := device.Dump()
	if code != 0 {
		test.Error("Dump() should have returned `0`.")
	}

	if device.Free() != true {
		test.Error("Free should have returned `true`.")
	}

	code = device.Dump()
	if code != -22 {
		test.Error("Dump() should have returned `-22`.")
	}
}

func Test_Device_Free_Doesnt_Fail_For_Empty_Device(test *testing.T) {
	device := &cryptsetup.Device{}

	if device.Free() != true {
		test.Error("Free should have returned `true`.")
	}

	if device.Free() != false {
		test.Error("Free should have returned `false`.")
	}
}

func Test_Device_Free_Doesnt_Fail_If_Called_Multiple_Times(test *testing.T) {
	testWrapper := TestWrapper{test}

	device, err := cryptsetup.Init(DevicePath)
	testWrapper.AssertNoError(err)

	err = device.Format(
		cryptsetup.LUKS1{Hash: "sha256"},
		cryptsetup.GenericParams{Cipher: "aes", CipherMode: "xts-plain64", VolumeKeySize: 512 / 8},
	)
	testWrapper.AssertNoError(err)

	if device.Free() != true {
		test.Error("Free should have returned `true`.")
	}

	if device.Free() != false {
		test.Error("Free should have returned `false`.")
	}
}

func Test_Device_Deactivate_Fails_If_Device_Is_Not_Active(test *testing.T) {
	testWrapper := TestWrapper{test}

	device, err := cryptsetup.Init(DevicePath)
	testWrapper.AssertNoError(err)

	err = device.Deactivate(DeviceName)
	testWrapper.AssertError(err)
	testWrapper.AssertErrorCodeEquals(err, -19)
}

func Test_Device_ActivateByPassphrase_Fails_If_Device_Has_No_Type(test *testing.T) {
	testWrapper := TestWrapper{test}

	device, err := cryptsetup.Init(DevicePath)
	testWrapper.AssertNoError(err)

	err = device.ActivateByPassphrase(DeviceName, 0, "testPassphrase", cryptsetup.CRYPT_ACTIVATE_READONLY)
	testWrapper.AssertError(err)
	testWrapper.AssertErrorCodeEquals(err, -22)
}

func Test_Device_ActivateByVolumeKey_Fails_If_Device_Has_No_Type(test *testing.T) {
	testWrapper := TestWrapper{test}

	genericParams := cryptsetup.GenericParams{
		Cipher:        "aes",
		CipherMode:    "xts-plain64",
		VolumeKey:     generateKey(32, test),
		VolumeKeySize: 32,
	}

	device, err := cryptsetup.Init(DevicePath)
	testWrapper.AssertNoError(err)

	err = device.ActivateByVolumeKey(DeviceName, genericParams.VolumeKey, genericParams.VolumeKeySize, cryptsetup.CRYPT_ACTIVATE_READONLY)
	testWrapper.AssertError(err)
	testWrapper.AssertErrorCodeEquals(err, -22)
}

func Test_Device_KeyslotAddByVolumeKey_Fails_If_Device_Has_No_Type(test *testing.T) {
	testWrapper := TestWrapper{test}

	device, err := cryptsetup.Init(DevicePath)
	testWrapper.AssertNoError(err)

	err = device.KeyslotAddByVolumeKey(0, "", "testPassphrase")
	testWrapper.AssertError(err)
	testWrapper.AssertErrorCodeEquals(err, -22)
}

func Test_Device_KeyslotAddByPassphrase_Fails_If_Device_Has_No_Type(test *testing.T) {
	testWrapper := TestWrapper{test}

	device, err := cryptsetup.Init(DevicePath)
	testWrapper.AssertNoError(err)

	err = device.KeyslotAddByPassphrase(0, "testPassphrase", "secondTestPassphrase")
	testWrapper.AssertError(err)
	testWrapper.AssertErrorCodeEquals(err, -22)
}

func Test_Device_KeyslotChangeByPassphrase_Fails_If_Device_Has_No_Type(test *testing.T) {
	testWrapper := TestWrapper{test}

	device, err := cryptsetup.Init(DevicePath)
	testWrapper.AssertNoError(err)

	err = device.KeyslotChangeByPassphrase(0, 0, "testPassphrase", "secondTestPassphrase")
	testWrapper.AssertError(err)
	testWrapper.AssertErrorCodeEquals(err, -22)
}
