## DMI设计方案

### KubeEdge跟EdgeX对比

|              | KubeEdge                                 | EdgeX                                                        |
| ------------ | ---------------------------------------- | ------------------------------------------------------------ |
| 历史数据存储 | 无                                       | 支持边缘本地数据库对接                                       |
| 安全         | 无                                       | Security Services支持存放加密数据  支持API网关               |
| 数据规则引擎 | 无                                       | 支持，默认kuiper                                             |
| 云边协同     | 支持                                     | 需要依赖第三方云边平台如openyurt                             |
| 南向设备协议 | BlueTooth、ModBus、OPCUA、Onvif          | 十多种                                                       |
| 微服务化     |                                          | 支持，每个微服务组件可以替换                                 |
| 北向数据导出 | mqtt topic /data/update                  | Export Services  mqtt、HTTPS、ZMQ等                          |
| 设备Command  | 无                                       | Command Service支持具体命令定义和执行                        |
| 部署模式     | edgecore部署在边缘侧                     | 部署方式灵活，支持分布式，核心服务可以部署在云端、网关或边缘侧都可以 |
| 部署依赖     | 依赖K8s和CRI，kubeedge本身组件较为轻量化 | 不依赖底层硬件和操作系统，但EdgeX平台自身较重，需要资源较多  |
| 告警         |                                          | Supporting Services支持                                      |
| 扩展服务层   |                                          | Supporting Services支持                                      |
| 调度         |                                          | Supporting Services支持，EdgeX内部时钟，可以定时操作         |
| 日志         |                                          | Supporting Services支持                                      |
| SDK          | 未支持                                   | 提供GO和C的SDK                                               |
| API          | 未给出统一标准                           | 已提供                                                       |



### KubeEdge设备管理中存在的问题

1. 未适配EdgeX等第三方设备管理平台，用户迁移成本高
2. 设备通讯协议众多，缺乏标准化，难以维护，也难以为后续第三方厂商提供设备接入认证方法
3. 用户想接入新协议的设备，开发新的mapper难度大
4. 设备安全性不完善（待补充）
5. kubeedge内部南北向Device Resource 定义不统一
6. 功能不完善
   1. 缺乏对设备的监控指标：包括state，有接口但未用起来
   2. 缺乏对设备升级能力的检测：检测设备是否支持OTA升级能力
   3. 缺乏对设备管理平台的管理能力：检测设备管理平台状态、版本信息等
   4. 缺乏对device config的管理能力：目前device的config是以配置文件的形式保存在mapper本地的
   5. 缺乏对边缘设备产生的流式数据、历史数据的管理能力：目前是直接将边缘数据通过/data/update主题发布到mqtt broker，由用户应用自己获取使用。（工业相机一帧一帧的图片，适合mapper统一管理）
   6. 缺乏对边缘设备自动接入的检测能力：目前边缘设备需要先提供接入的配置文件，并在kubeedge中创建device模型，才能进行devicetwin建模（后续再对一下）

* //DMI对第三方设备平台的对接，整合进来，覆盖主流应用场景
* //KubeEdge本身的能力可以后续再进行优化

### UserStory

* 作为一名KubeEdge社区开发者，在开发设备管理相关特性时，
* 作为一名开发者或设备提供商，需要能够基于DMI提供的SDK，快速开发出某个协议的设备对应的mapper，并将该设备接入到KubeEdge中
* 作为一名设备使用者，在当前设备已经使用EdgeX平台管理起来的情况下，能够不做大量二次开发工作就能够将EdgeX对接到KubeEdge上
* 作为一名云端设备管理者，需要能够在云端获取指定设备当前的运行状态state，包括连接状态、CPU内存状态等
* 作为一名云端设备管理者，需要能够在云端更新指定设备运行状态，包括新建设备并连接、断开连接并删除设备等
* 作为一名云端设备管理者，需要能够获取所有接入设备的统计信息，如蓝牙设备mapper已接入设备数量
* 作为一名云端设备管理者，需要能够获取具体设备管理平台信息，如某个指定平台已接入设备数量、当前健康状态、组件当前版本
* 作为一名云端设备管理者，需要能够在云端对指定设备状态status进行查询和修改，如更新灯的开关状态
* 作为一名云端设备管理者，需要能够在云端获取设备对OTA升级能力的支持情况，并调用接口让设备升级
* 作为一名本地设备管理者，需要能够设备管理平台自动发现可连接设备，并通过shim将设备信息同步到云上
* //作为一名本地设备管理者，需要能够让mapper在未连接edgecore的情况下本地离线运行，连接本地设备并将数据data发布到mqtt broker或REST Server
* 作为一名设备数据使用者，需要能够以mqtt或REST的方式，获取设备实时推送的数据，并能够通过接口获取设备历史数据
* 原始用户场景待补充，比如通过摄像头进行，比如停车场，原子能力组合使用，路灯管理。社区例会收集一下原始场景。

### DMI方案

#### 架构

整体KubeEdge边缘侧使用DeviceManager来管理设备，DeviceTwin下沉为独立模块，用于兼容老的设备管理接口及我们自己的mapper接入，并新增一个DMI shim用于对接EdgeX等第三方平台。DMI作为一套边缘框架，指导Device Manager、DeviceTwin、DMI shim与mapper的开发，保持接口统一和数据格式一致。

在这个架构中，可以把DeviceManager看做kubelet，DMI shim就是CRI的shim，EdgeX平台或DeviceTwin作为containerd，每个Device作为一个pod或container，跟k8s中的Device CR对应。

CSI

![image-20220301153831857](C:\Users\z00525294\AppData\Roaming\Typora\typora-user-images\image-20220301153831857.png)

#### EdgeX 对接

EdgeX的定位是边缘侧设备管理平台，自身已经包括边缘设备的增删改查、边缘设备协议转换接入、数据筛选、数据存储、数据传输、边缘设备COMMAND接口等能力。EdgeX自身功能比较强大，目前也已经有不少用户在使用，但EdgeX自身没有做云平台，可以通过接口对接其他云平台实现云边协同。

新增一个独立外部模块DMI shim，用于对接EdgeX，实际相当于一个接口转换和数据格式清洗的模块。state和status类型可以采用定时向下拉取的方式获取，跟mapper的方式保持一致。缺点是state和status数据存储在三个db中，云上k8s etcd、边缘侧的sqlite以及EdgeX的本地数据库，可能存在数据同步问题。（注：设备场景对于status的实时性要求和快速响应要求不清楚该如何取舍，即是否每次查询都走一次完整链路，要求准确但不要求速度，还是有一定响应速度要求，只在边缘数据库侧拉取结果即可？）从DeviceManager的角度来看，是应该同步到云上etcd的，这样用户只需要看到k8s里面的数据即可，更符合数字孪生的定义。



EdgeX与mapper模式的区别在于

* EdgeX可以使用REST方式，一来一回双向通信，而且EdgeX本地提供数据库，可以不用中转
* mapper使用mqtt方式，只支持单向通信，双向的话需要做很多额外工作量，而且本地不提供数据库，只能通过DeviceTwin转存在sqlite中，或实时获取状态数据





在对接EdgeX的接口中，应该新增一个字段，选择是否mapper直连设备还是通过第三方设备管理平台转接。如果是标注为通过EdgeX进行管理，则用户新建的Device CR应该包含EdgeX Device相关格式。

对接EdgeX需要将KubeEdge的Device相关CR跟EdgeX的Device相关CR适配起来。

![image-20220218114012561](C:\Users\z00525294\AppData\Roaming\Typora\typora-user-images\image-20220218114012561.png)

#### config管理

目前mapper连接device的接入信息是以config文件的形式写在mapper所在节点上，mapper初始化的时候会加载该数据到内存。这种方式下，新增device的操作会比较麻烦。这里希望把config也变成一种k8s的资源类型，当然这里也只能是一种泛化地定义方式，由云上下发到边缘节点，然后边缘侧由deviceTwin保存到边缘数据库中做持久化，调用新建设备的接口，向mapper下发配置文件，mapper每次初始化时，也会调用特定接口拉取这些config到内存（有点难）。但也要兼容mapper本地配置，mapper启动后先从DeviceTwin拉取该mapperID对应的device config列表，如果拉取失败，则读取本地config配置。

#### 提供DMI SDK（南向）

#### 

### DMI南向接口设计

设备相关数据类型分为state、status和data三种，state代表device的连接状态、运行状态，status和data均为数据，status为可读写类型数据，如灯颜色，data为只读数据，如温度湿度，故可规定如下几类接口：

* 配置相关

* state的读写

* status的读写

* data的读取（可以跟上面合并）

  ![image-20220228162635019](C:\Users\z00525294\AppData\Roaming\Typora\typora-user-images\image-20220228162635019.png)

功能：

DeviceManager：

1. 创建device：
2. 删除device
3. 更新device的基本信息
4. 更新device的state信息
5. 更新device的status
6. 查询device的status
7. 查询device的status的具体字段
8. 查询device的基本信息
9. 查询device list
10. 查询device的state信息
11. 查询device的data导出信息

 

DeviceCommandManager

1. 在device上执行command（可选）
2. 获取框架上command信息（可选）

 

DeviceUpgradeManager

1. 获取device的升级能力信息

 

DeviceEventManager

1. 获取device event信息（可选）

 

DevicePlatformManager

1. 获取框架统计信息

2. 设备自动接入检测能力（可选）

3. 获取框架的版本信息

4. 获取框架的健康检查状态

5. 创建框架连接

   

仿照CRI的方式，给出DeviceManagerService的定义

````go
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
````

目前已定义DevicePlatformVersioner、DeviceManager、DeviceCommandManager、DeviceUpgradeManager、DeviceEventManager几种manager

```go
type DevicePlatformVersioner interface {
	// Version returns the device mapper name, device mapper version and device mapper API version
	Version(apiVersion string) (*dmiapi.VersionResponse, error)

	// list the API provided
	ListAPI()

	HealthCheck()

	RegisterPlatform()

	GetPlatform()
}
```



```go
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
	DeviceStatus(deviceID string) (*dmiapi.DeviceStatus, error)
	// UpdateDevice updates the status of the device.
	UpdateDeviceStatus(deviceID string, desiredDevice *dmiapi.Device) error

	GetDevice()

	GetDeviceDataInfo()
}
```



```go
type DeviceCommandManager interface {
	ListCommand()

	GetCommand()

	ExecCommand()
}
```



```go
type DeviceUpgradeManager interface {
	CheckUpgrade()

	UpgradeDevice()
}
```

```go
type DeviceEventManager interface {
	RegisterEvent()
	GetEvent()
	ListEvent()
}
```



### DMI北向CRD数据设计

KubeEdge不关心边缘Device是如何连接到mapper的，所以在Device本身的CRD设计中弱化了Address的内容，用户通过更新configmap，将Address直接下发给mapper，或者由Device管理员直接进行管理。

EdgeX还开放了对设备的具体操作接口，这一点KubeEdge目前还不支持

KubeEdge的设备接入模式跟EdgeX不太一样

* KubeEdge是先在mapper写入配置，接入设备，再从云端创建Device的数字孪生，并在DeviceTwin进行对齐
* EdgeX是直接通过云上调用Create接口，在EdgeX处新建一次设备
* 所以对接EdgeX需要create接口，对接mapper不需要
* 对接EdgeX，需要保留command接口



目前DeviceInstance对应EdgeX的Device，但是里面的具体定义不是很清晰，作为一个K8s资源类型，应该都是desired，所以Twins对应的应该是一次变更过程的中间变量，放在status里，感觉不够原生。

Property Visitors对应EdgeX的Device Service

DeviceModel对应EdgeX的Device Profile



DeviceModel：

* ID or Name
* Description
* Protocol
* Device Command
* Manufacturer
* Device Properties（定义property字段具体类型）status
* Device Data（定义data字段具体类型）
* Device State（定义state字段具体类型）



DeviceInstance：

* ID or Name
* Description
* Model（其实也可以放在label里面）
* PropertyVisitors（为了做到对property粒度的资源上报，mapper上报的topic也需要对应property级别）
* Data（这里只能获取data的metadata和接入方式，而不能获取data具体值）
* Properties（相当于spec）（使用twin主要是为了解决device设备状态更新不同步的问题）
* State（描述运行状态）
* Address

#### 遗留问题

* 目前只考虑了节点、DeviceManager与平台一对一的场景，后续要考虑每个节点上的DeviceManager与platform一对多甚至多对多的问题

* 考虑新接入框架的时候要如何扩展，不重启的情况下，或者把框架也作为一种k8s资源进行管理下发？需要有个配置参数来传入CRI远端框架的服务端点
* device的数据是从下往上推的，或者上面定时向下拉取，做到device信息实时更新
* 

###  其他待补充

#### 设备自动接入检查（待补充）

#### OTA升级相关接口（待补充）

#### 安全（待补充）

#### DMI框架自动生成（待补充）

#### data类型数据获取方式