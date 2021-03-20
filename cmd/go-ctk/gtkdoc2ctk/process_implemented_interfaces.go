package gtkdoc2ctk

import (
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/kckrinke/go-cdk/utils"
)

func ProcessImplementedInterfaces(src *GtkSource, s *goquery.Selection) {
	links := s.Find("p > a.link")
	links.Each(func(i int, selection *goquery.Selection) {
		value := selection.Text()
		value = strings.TrimSpace(value)
		if strings.HasPrefix(value, "Gdk") || strings.HasPrefix(value, "Gtk") {
			value = value[3:]
		}
		if !utils.StringSliceHasValue(src.Implements, value) {
			src.Implements = append(src.Implements, value)
		}
	})
}
