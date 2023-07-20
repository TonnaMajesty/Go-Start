package main

import (
	"fmt"
	"log"

	"github.com/carmel/gooxml/document"
)

var lorem2 = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin lobortis, lectus dictum feugiat tempus, sem neque finibus enim, sed eleifend sem nunc ac diam. Vestibulum tempus sagittis elementum`

func main() {
	// When Word saves a document, it removes all unused styles.  This means to
	// copy the styles from an existing document, you must first create a
	// document that contains text in each style of interest.  As an example,
	// see the template.docx in this directory.  It contains a paragraph set in
	// each style that Word supports by default.
	doc, err := document.OpenTemplate("/Users/tonnamajesty/Downloads/变电站巡检报告模板.docx")
	if err != nil {
		log.Fatalf("error opening Windows Word 2016 document: %s", err)
	}

	// We can now print out all styles in the document, verifying that they
	// exist.
	for _, s := range doc.Styles.Styles() {
		fmt.Println("style ", s.Name(), "id ", s.StyleID(), "type ", s.Type())
	}

	nds := doc.Numbering.Definitions()
	fmt.Println(len(nds))
}
