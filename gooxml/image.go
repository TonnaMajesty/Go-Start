package main

import (
	"io/ioutil"
	"log"

	"github.com/carmel/gooxml/common"
	"github.com/carmel/gooxml/document"
	"github.com/carmel/gooxml/measurement"
)

var lorem = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin lobortis, lectus dictum feugiat tempus, sem neque finibus enim, sed eleifend sem nunc ac diam. Vestibulum tempus sagittis elementum`

func main() {
	doc := document.New()

	imgBytes, err := ioutil.ReadFile("/Users/tonnamajesty/Downloads/1425547475187859812-20230620080828.jpeg")
	img1, err := common.ImageFromBytes(imgBytes)

	//img1, err := common.ImageFromFile("/Users/tonnamajesty/Downloads/1425547475187859812-20230620080828.jpeg")
	//if err != nil {
	//	log.Fatalf("unable to create image: %s", err)
	//}

	//img2data, err := ioutil.ReadFile("/Users/tonnamajesty/Downloads/charge_test.jpeg")
	//if err != nil {
	//	log.Fatalf("unable to read file: %s", err)
	//}
	//img2, err := common.ImageFromBytes(img2data)
	//if err != nil {
	//	log.Fatalf("unable to create image: %s", err)
	//}

	img1ref, err := doc.AddImage(img1)
	if err != nil {
		log.Fatalf("unable to add image to document: %s", err)
	}
	//img2ref, err := doc.AddImage(img2)
	//if err != nil {
	//	log.Fatalf("unable to add image to document: %s", err)
	//}

	para := doc.AddParagraph()
	//anchored, err := para.AddRun().AddDrawingAnchored(img1ref)
	//if err != nil {
	//	log.Fatalf("unable to add anchored image: %s", err)
	//}
	//anchored.SetName("Gopher")
	//anchored.SetSize(2*measurement.Inch, 2*measurement.Inch)
	//anchored.SetOrigin(wml.WdST_RelFromHPage, wml.WdST_RelFromVTopMargin)
	//anchored.SetHAlignment(wml.WdST_AlignHCenter)
	//anchored.SetYOffset(3 * measurement.Inch)
	//anchored.SetTextWrapSquare(wml.WdST_WrapTextBothSides)

	run := para.AddRun()
	//for i := 0; i < 16; i++ {
	//	run.AddText(lorem)
	//
	//	// drop an inline image in
	//	if i == 13 {
	inl, err := run.AddDrawingInline(img1ref)
	if err != nil {
		log.Fatalf("unable to add inline image: %s", err)
	}
	inl.SetSize(measurement.Distance(512)*measurement.Point, measurement.Distance(290)*measurement.Point)
	//}
	//	if i == 15 {
	//		inl, err := run.AddDrawingInline(img2ref)
	//		if err != nil {
	//			log.Fatalf("unable to add inline image: %s", err)
	//		}
	//		inl.SetSize(1*measurement.Inch, 1*measurement.Inch)
	//	}
	//}
	doc.SaveToFile("/Users/tonnamajesty/Downloads/image.docx")
}
