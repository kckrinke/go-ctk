package gtkdoc2ctk

import (
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/kckrinke/go-cdk"
	"github.com/kckrinke/go-cdk/utils"
)

func ProcessObjectHierarchy(src *GtkSource, s *goquery.Selection) {
	pre := s.Find("pre.screen")
	text, err := pre.Html()
	if err != nil {
		cdk.Error(err)
		return
	}
	text = rxStripLineArt.ReplaceAllString(text, "")
	text = rxStripTags.ReplaceAllString(text, "")
	lines := strings.Split(text, "\n")
	last := ""
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			if line == "GObject" {
				line = "Object"
			} else if line == "GInterface" {
				line = "CInterface"
			} else if line == "GInitiallyUnowned" {
				src.Hierarchy = []string{}
				continue
			} else if strings.HasPrefix(line, "Gdk") {
				line = line[3:]
			} else if strings.HasPrefix(line, "Gtk") {
				line = line[3:]
			}
			src.Hierarchy = append(src.Hierarchy, line)
			if src.Parent == "" && line == src.Name {
				src.Parent = last
			}
			last = line
			if !utils.StringSliceHasValue(src.Hierarchy, line) {
				src.Hierarchy = append(src.Hierarchy, line)
			}
		}
	}
}
