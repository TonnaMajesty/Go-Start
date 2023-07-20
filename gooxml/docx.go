package main

import (
	"github.com/carmel/gooxml/document"
)

func AddHead(doc *document.Document, content string, style string) {
	headP := doc.AddParagraph()
	headP.SetStyle(style)
	head := headP.AddRun()
	head.AddText(content)
	head.Properties().SetBold(true)
}

func AddUl(doc *document.Document, content []string) {
	nd := doc.Numbering.Definitions()[0]
	for _, c := range content {
		p := doc.AddParagraph()

		p.SetNumberingDefinition(nd)
		r := p.AddRun()
		r.AddText(c)
	}
}
