package dmi

import (
	dmiapi "github.com/kubeedge/dmi-api/pkg/apis/manager/v1"
)

// DevicePlatformVersioner contains methods for platform name, version and API version.
type DevicePlatformVersioner interface {
	// Version returns the device mapper name, device mapper version and device mapper API version
	Version(apiVersion string) (*dmiapi.VersionResponse, error)

	// list the API provided
	ListAPI()

	HealthCheck()

	RegisterPlatform()

	GetPlatform()
}

// DeviceCommandManager contains methods for retrieving the device
// statistics.
type DeviceCommandManager interface {
	ListCommand()

	GetCommand()

	ExecCommand()
}

// DeviceUpgradeManager contains methods for retrieving the device
// statistics.
type DeviceUpgradeManager interface {
	CheckUpgrade()

	UpgradeDevice()
}

// DeviceEventManager contains methods for retrieving the device
// statistics.
type DeviceEventManager interface {
	RegisterEvent()
	GetEvent()
	ListEvent()
}

// DeviceManager contains methods to manipulate devices managed by a
// device mapper. The methods are thread-safe.
// 这里暂时按照CRI的方式，DeviceConfig保存在内存中，mapper启动的时候，会去deviceTwin取一下。
type DeviceManager interface {
	// FetchDevice fetches config of devices from deviceTwin
	//FetchDevice()

	// CreateDevice creates a new device.
	CreateDevice(podSandboxID string, config *dmiapi.DeviceConfig) (string, error)
	// UpdateDeviceState update device state with a grace period (i.e., timeout).
	UpdateDeviceState(deviceID string, timeout int64) error

	GetDeviceState(deviceID string) (*dmiapi.DeviceState, error)

	// RemoveDevice removes the device from platform.
	RemoveDevice(deviceID string) error
	// ListDevices lists all devices by filters.
	ListDevices(filter *dmiapi.DeviceFilter) ([]*dmiapi.Device, error)
	// DeviceStatus returns the status of the device.
	GetDeviceStatus(deviceID string) (*dmiapi.DeviceStatus, error)
	// UpdateDevice updates the status of the device.
	UpdateDeviceStatus(deviceID string, desiredDevice *dmiapi.Device) error

	GetDevice()

	GetDeviceDataInfo()
}

// DeviceMapperService interface should be implemented by a device mapper.
// The methods should be thread-safe.
type DeviceManagerService interface {
	DevicePlatformVersioner
	DeviceManager
	DeviceCommandManager

	// fields to be extended
	DeviceUpgradeManager
	DeviceEventManager

	// UpdateDeviceMapperConfig updates device mapper configuration if specified
	UpdateDeviceMapperConfig(runtimeConfig *dmiapi.RuntimeConfig) error
	// Status returns the status of the device mapper.
	Status() (*dmiapi.RuntimeStatus, error)
}
