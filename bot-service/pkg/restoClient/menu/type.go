package menu

import evotorrestogo "github.com/softc24/evotor-resto-go"

type MenuResponse evotorrestogo.MenuItem
type Menu []evotorrestogo.MenuItem

func (m *MenuResponse) CanAddToOrder() bool {
	if m.Group {
		return true
	}

	if !m.Group && m.Type == "NORMAL" {
		return true
	}
	//
	//if m.Group && (m.IsUnavailable == "" || m.IsUnavailable == "0") {
	//	return true
	//}
	//
	//if !m.Group && (m.IsUnavailable == "" || m.IsUnavailable == "0") && m.Type == "NORMAL" {
	//	return true
	//}

	return false
}

func (m *MenuResponse) ImageNotEmpty() bool {

	if m.ImageURL != "" {
		return true
	}
	return false
}
