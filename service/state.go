// +build linux

package service

import (
	"errors"
	pb "github.com/mame82/P4wnP1_go/proto"
	"github.com/mame82/P4wnP1_go/service/datastore"
)

const (
	cSTORE_PATH = "/tmp/P4wnP1_store"
)

var ServiceState *GlobalServiceState

type GlobalServiceState struct {
	Store                 *datastore.Store
	EvMgr                 *EventManager
	UsbGM                 *UsbGadgetManager
	Led                   *LedState
	HidDevPath            map[string]string //stores device path for HID devices
	StoredNetworkSettings map[string]*pb.EthernetInterfaceSettings
	WifiSvc               *WiFiService
	BtSvc                 *BtService
}

func InitGlobalServiceState() (err error) {
	state := &GlobalServiceState{}
	ServiceState = state // store state in global variable

	state.Store,err = datastore.Open(cSTORE_PATH)
	if err != nil { return }
	state.StoredNetworkSettings = make(map[string]*pb.EthernetInterfaceSettings)

	/*
	state.StoredNetworkSettings[USB_ETHERNET_BRIDGE_NAME] = GetDefaultNetworkSettingsUSB()
	state.StoredNetworkSettings["wlan0"] = GetDefaultNetworkSettingsWiFi()
	*/
	//pre initialize Default settings for "wlan0" and USB_ETHERNET_BRIDGE_NAME ("usbeth")
	state.StoredNetworkSettings[USB_ETHERNET_BRIDGE_NAME] = &pb.EthernetInterfaceSettings{
		Name:       USB_ETHERNET_BRIDGE_NAME,
		Enabled:    false,
		Mode:       pb.EthernetInterfaceSettings_MANUAL,
		IpAddress4: "172.16.0.1",
		Netmask4:   "255.255.255.252",
	}
	state.StoredNetworkSettings["wlan0"] = &pb.EthernetInterfaceSettings{
		Name:       "wlan0",
		Enabled:    false,
		Mode:       pb.EthernetInterfaceSettings_MANUAL,
		IpAddress4: "172.24.0.1",
		Netmask4:   "255.255.255.0",
	}
	state.StoredNetworkSettings[BT_ETHERNET_BRIDGE_NAME] = &pb.EthernetInterfaceSettings{
		Name:       BT_ETHERNET_BRIDGE_NAME,
		Enabled:    true,
		Mode:       pb.EthernetInterfaceSettings_MANUAL,
		IpAddress4: "172.26.0.1",
		Netmask4:   "255.255.255.0",
	}

	state.WifiSvc = NewWifiService()

	state.HidDevPath = make(map[string]string) //should be initialized BEFORE UsbGadgetManager uses it
	state.EvMgr = NewEventManager(20)
	state.UsbGM, err = NewUSBGadgetManager()
	if err != nil {
		return
	}

	state.BtSvc = NewBtService()

	ledState, err := NewLed(false)
	if err != nil {
		return
	}
	state.Led = ledState

	return nil
}

func (state *GlobalServiceState) GetInterfaceSettingsByInterfaceName(ifname string) (*pb.EthernetInterfaceSettings, error) {
	for _, s := range state.StoredNetworkSettings {
		if s.Name == ifname {
			return s, nil
		}
	}
	return nil, errors.New("No settings for interface " + ifname + " found")
}

func (state *GlobalServiceState) StartService() {
	state.EvMgr.Start()
	// ToDo: Remove this, till the service is able to deploy startup settings
	state.BtSvc.StartNAP()
}

func (state *GlobalServiceState) StopService() {
	state.EvMgr.Stop()
}
