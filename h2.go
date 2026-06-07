package illutls

import (
	http2 "github.com/bogdanfinn/fhttp/http2"
)

// ApplyH2Settings configures an fhttp http2.Transport to match the
// browser profile's HTTP/2 SETTINGS frame parameters.
func ApplyH2Settings(t *http2.Transport, s H2Settings, windowUpdate uint32, priority H2Priority) {
	// For bogdanfinn/fhttp, we must use Settings map and SettingsOrder
	// to explicitly control the SETTINGS frame.
	if t.Settings == nil {
		t.Settings = make(map[http2.SettingID]uint32)
	}
	if s.HeaderTableSize > 0 {
		t.HeaderTableSize = s.HeaderTableSize
		t.Settings[http2.SettingHeaderTableSize] = s.HeaderTableSize
		t.SettingsOrder = append(t.SettingsOrder, http2.SettingHeaderTableSize)
	}
	// SettingEnablePush (2) is sent as 0
	t.Settings[http2.SettingEnablePush] = s.EnablePush
	t.SettingsOrder = append(t.SettingsOrder, http2.SettingEnablePush)
	if s.MaxConcurrentStreams > 0 {
		t.Settings[http2.SettingMaxConcurrentStreams] = s.MaxConcurrentStreams
		t.SettingsOrder = append(t.SettingsOrder, http2.SettingMaxConcurrentStreams)
	}
	if s.InitialWindowSize > 0 {
		t.InitialWindowSize = s.InitialWindowSize
		t.Settings[http2.SettingInitialWindowSize] = s.InitialWindowSize
		t.SettingsOrder = append(t.SettingsOrder, http2.SettingInitialWindowSize)
	}
	if s.MaxHeaderListSize > 0 {
		t.Settings[http2.SettingMaxHeaderListSize] = s.MaxHeaderListSize
		t.SettingsOrder = append(t.SettingsOrder, http2.SettingMaxHeaderListSize)
	}
	if s.NoRFC7540Priorities > 0 {
		t.Settings[http2.SettingID(9)] = s.NoRFC7540Priorities
		t.SettingsOrder = append(t.SettingsOrder, http2.SettingID(9))
	}
	if windowUpdate > 0 {
		t.ConnectionFlow = windowUpdate
	}
	if priority.Weight > 0 {
		t.HeaderPriority = &http2.PriorityParam{
			Weight:    priority.Weight,
			StreamDep: priority.DependsOn,
			Exclusive: priority.Exclusive,
		}
	}
}
