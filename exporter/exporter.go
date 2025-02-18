package exporter

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/richartkeil/nplan/core"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Export(path string, scan *core.Scan) {
	cols := 5
	width := 210
	height := 120
	padding := 30

	cells := make([]MxCell, 0)
	cells = append(cells, MxCell{
		Id: "0",
	})
	cells = append(cells, MxCell{
		Id:     "1",
		Parent: "0",
	})
	for i, host := range scan.Hosts {
		cells = append(cells, MxCell{
			Id:     uuid.NewString(),
			Value:  getHostValue(host),
			Parent: "1",
			Style:  "rounded=1;whiteSpace=wrap;html=1;arcSize=2",
			Vertex: "1",
			MxGeometry: &MxGeometry{
				X:      fmt.Sprint((i % cols) * (width + padding)),
				Y:      fmt.Sprint((i / cols) * (height + padding)),
				Width:  fmt.Sprint(width),
				Height: fmt.Sprint(height),
				As:     "geometry",
			},
		})
	}

	mxFile := MxFile{
		Diagram: &Diagram{
			Id:   uuid.NewString(),
			Name: "Network Plan",
			MxGraphModel: &MxGraphModel{
				Root: &Root{
					MxCells: cells,
				},
				Dx:       "3000",
				Dy:       "2000",
				Grid:     "1",
				GridSize: "10",
				Guides:   "1",
				Tooltips: "1",
				Connect:  "1",
				Arrows:   "1",
			},
		},
	}

	output, err := xml.MarshalIndent(mxFile, "", "  ")
	check(err)

	os.WriteFile(path, output, 0644)
}

func getHostValue(host core.Host) string {
	serviceColor := "#bbb"

	value := ""
	if host.Hostname != "" {
		value += fmt.Sprintf("<strong>%v</strong><br>", host.Hostname)
	}
	if host.IPv4 != "" {
		value += fmt.Sprintf("%v<br>", host.IPv4)
	}
	if host.IPv6 != "" {
		value += fmt.Sprintf("%v<br>", host.IPv6)
	}

	value += "<br>"

	for _, port := range host.Ports {
		value += fmt.Sprintf(
			":%v - %v<br><span style=\"color: %v\">(%v)</span><br>",
			port.Number,
			port.ServiceName,
			serviceColor,
			port.ServiceVersion,
		)
	}

	return value
}
