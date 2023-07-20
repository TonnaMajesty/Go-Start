package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/carmel/gooxml/color"
	"github.com/carmel/gooxml/common"
	"github.com/carmel/gooxml/document"
	"github.com/carmel/gooxml/measurement"
	"github.com/carmel/gooxml/schema/soo/ofc/sharedTypes"
	"github.com/carmel/gooxml/schema/soo/wml"
)

func main() {
	doc, _ := document.OpenTemplate("/Users/tonnamajesty/www/go/industai/srv-transmission-line/internal/resource/inspection_task_report.docx")
	//bodyx := doc.X()
	//fmt.Println(bodyx)
	//bodyx.Body.SectPr = wml.NewCT_SectPr()
	//bodyx.Body.SectPr.PgSz = wml.NewCT_PageSz()
	//bodyx.Body.SectPr.PgSz.OrientAttr = wml.ST_PageOrientationLandscape
	//var (
	//	height uint64 = 11906
	//	width  uint64 = 16838
	//)
	//bodyx.Body.SectPr.PgSz.HAttr = &sharedTypes.ST_TwipsMeasure{
	//	ST_UnsignedDecimalNumber: &height,
	//}
	//bodyx.Body.SectPr.PgSz.WAttr = &sharedTypes.ST_TwipsMeasure{
	//	ST_UnsignedDecimalNumber: &width,
	//}

	// 添加标题
	titleP := doc.AddParagraph()
	titleP.Properties().SetAlignment(wml.ST_JcCenter)
	//titleP.Properties().Spacing().SetBefore(3.5 * measurement.Centimeter)
	//titleP.Properties().Spacing().SetAfter(3.5 * measurement.Centimeter)
	//titleP.Properties().Spacing().SetLineSpacing(260*measurement.Twips, wml.ST_LineSpacingRuleAuto)

	titleR := titleP.AddRun()
	titleR.Properties().SetSize(16) // 初号
	titleR.Properties().SetFontFamily("微软雅黑")
	titleR.AddText("输电线路智能巡检报告-异常事件")

	analyzeTimeP := doc.AddParagraph()
	analyzeTimeP.Properties().SetAlignment(wml.ST_JcEnd)

	timeR := analyzeTimeP.AddRun()
	timeR.Properties().SetSize(8)
	timeR.Properties().SetFontFamily("微软雅黑")
	timeR.AddText(fmt.Sprintf("分析开始时间：%s", "2020-01-01 00:00:00"))

	// summary
	summaryTable := doc.AddTable()
	summaryTable.Properties().SetWidthPercent(100)

	summaryTableBorders := summaryTable.Properties().Borders()
	summaryTableBorders.SetBottom(wml.ST_BorderSingle, color.Black, 1*measurement.Point)
	summaryTableBorders.SetTop(wml.ST_BorderSingle, color.Black, 1*measurement.Point)
	summaryTableBorders.SetLeft(wml.ST_BorderSingle, color.White, 1*measurement.Point)
	summaryTableBorders.SetRight(wml.ST_BorderSingle, color.White, 1*measurement.Point)

	summaryHeaderRow := summaryTable.AddRow()
	summaryHeaderCell := summaryHeaderRow.AddCell()

	summaryHeaderCellP := summaryHeaderCell.AddParagraph()
	//summaryHeaderCellP.Properties().Spacing().SetLineSpacing(26*measurement.Twips, wml.ST_LineSpacingRuleAuto)
	//summaryHeaderCellP.Properties().Spacing().SetAfter(50 * measurement.Twips)
	summaryHeaderCellR := summaryHeaderCellP.AddRun()
	summaryHeaderCellR.Properties().SetFontFamily("微软雅黑")
	summaryHeaderCellR.Properties().SetSize(11)

	//summaryHeaderCellR.Properties().SetCharacterSpacing(26 * measurement.Twips)
	//summaryHeaderCellR.Properties().SetBold(true)
	summaryHeaderCellR.AddText("输电线路名称：110kV洞太一线")
	summaryHeaderCellR.AddTab()
	summaryHeaderCellR.AddTab()
	summaryHeaderCellR.AddTab()
	summaryHeaderCellR.AddText("杆塔数量：100基")
	summaryHeaderCellR.AddBreak()
	summaryHeaderCellR.AddText("巡视点位数：2000张")
	summaryHeaderCellR.AddBreak()

	doc.AddParagraph().AddRun().AddBreak()

	detailTitleP := doc.AddParagraph()
	detailTitleP.Properties().SetAlignment(wml.ST_JcStart)
	detailTitleR := detailTitleP.AddRun()
	detailTitleR.Properties().SetBold(true)
	detailTitleR.AddText("异常事件详情:")

	// 添加一个表格
	table := doc.AddTable()
	table.Properties().SetWidthPercent(100)

	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Black, 1*measurement.Point)

	// 添加第一行
	headerRow := table.AddRow()
	cell1 := headerRow.AddCell()
	cell1.AddParagraph().AddRun().AddText("Row 1, Column 1")

	// 添加第二行
	imageRow := table.AddRow()
	cell2 := imageRow.AddCell()

	imgP := cell2.AddParagraph()
	imgFile, err := os.Open("/Users/tonnamajesty/Downloads/110kV洞太二线019_2_1088_1920.jpg")
	if err != nil {
		fmt.Println(err)
	}
	imgByte, err := ioutil.ReadAll(imgFile)
	if err != nil {
		fmt.Println(err)
	}
	img, err := common.ImageFromBytes(imgByte)
	if err != nil {
		fmt.Println(err)
	}
	imgRef, err := doc.AddImage(img)
	if err != nil {
		fmt.Println(err)
	}
	inlineDrawing, err := imgP.AddRun().AddDrawingInline(imgRef)
	if err != nil {
		fmt.Println(err)
	}
	inlineDrawing.SetSize(measurement.Distance(512)*measurement.Point, measurement.Distance(288)*measurement.Point)

	err = doc.SaveToFile("/Users/tonnamajesty/Downloads/inspection-task-report.docx")
	if err != nil {
		fmt.Println(err)
	}

}

func main1() {
	//doc := document.New()
	doc, _ := document.OpenTemplate("/Users/tonnamajesty/Downloads/变电站巡检报告模板.docx")

	addHead(doc, "巡视点位设置原则", "2")
	addHead(doc, "巡视点位设置原则", "3")

	content1_1_slice := []string{
		`能从整体各方位观察主变箱体、基础、散热器、油枕及其各部连接处有无渗漏油痕迹，外观有无锈蚀、破损、变形、异物等其余异常情况，确保主变四面均有覆盖，且视角良好。能观测到除各部件间遮挡处、下方视角盲区处等少数区域外的大部分区域，其余区域列为全面巡视重点检查部位。`,
		`能从细节观测主变套管油位、本体油位、有载油位指示，能通过图像，判断油位分界面、或油位表计指针位置。`,
		`对主变本体、有载瓦斯继电器，及其附属集气盒，能从细节观测其窥视孔，排查出内部是否集气。`,
	}

	for _, c := range content1_1_slice {
		content1_1_p := doc.AddParagraph()
		content1_1_p.Properties().SetAlignment(wml.ST_JcLeft)
		content1_1_p.SetNumberingDefinitionByID(3)
		content1_1 := content1_1_p.AddRun()
		content1_1.AddText(c)
	}

	// 下一页
	doc.AddParagraph().AddRun().AddPageBreak()

	addHead(doc, "灵敏度配置原则", "2")

	{
		table_lmd := doc.AddTable()

		table_lmd.Properties().SetWidthPercent(100)
		// set borers
		borders := table_lmd.Properties().Borders()
		borders.SetAll(wml.ST_BorderSingle, color.Black, 1*measurement.Point)

		// 表头
		row := table_lmd.AddRow()

		lmd_bt_slice := []string{
			"区域", "灵敏度范围", "适用区域",
		}

		for _, s := range lmd_bt_slice {
			cell := row.AddCell()
			p := cell.AddParagraph()
			p.Properties().SetAlignment(wml.ST_JcCenter)
			r := p.AddRun()
			r.Properties().SetBold(true)
			r.Properties().SetSize(11)
			r.AddText(s)
		}

		lmd_slice := []string{
			`高灵敏度（红）`, `0.7<灵敏度≤1.0`, `重点设备区域：需细节、重点观测的设备部件；如主变套管、仪表表计、开关指示等。`,
			`较高灵敏（橙）`, `0.5<灵敏度≤0.7`, `各类常规监测设施（呼吸器、瓦斯继电器、绝缘子等）和封闭式设备外观表面（GIS外观、变压器外观等）；若受到环境光照影响，易误报，可酌情调整至中灵敏度区域。`,
			`中灵敏度（黄）`, `0.3<灵敏度≤0.5`, `易误报封闭式设备外观表面区域（GIS外观、变压器外观等），易受光照影响的液晶屏区域。各类指示灯、把手等区域。`,
		}

		for i, s := range lmd_slice {
			if i%3 == 0 {
				row = table_lmd.AddRow()
			}
			cell := row.AddCell()
			p := cell.AddParagraph()
			r := p.AddRun()
			r.Properties().SetSize(11)
			r.AddText(s)
		}
	}

	doc.AddParagraph().AddRun().AddPageBreak()

	addHead(doc, "调优情况", "2")
	addHead(doc, "站点信息", "3")

	addUl(doc, []string{"调优站点：内江-110kv双苏变电站", "调优时间：2022/12/24", "监控点数：20", "监控点分布图："})

	img1, _ := common.ImageFromFile("/Users/tonnamajesty/Downloads/station_map.png")
	img1ref, _ := doc.AddImage(img1)

	img1_p := doc.AddParagraph()
	img1_p.Properties().SetAlignment(wml.ST_JcCenter)
	img1_r := img1_p.AddRun()
	img1_r.Properties().SetVerticalAlignment(sharedTypes.ST_VerticalAlignRunSuperscript)
	img1_inl, _ := img1_r.AddDrawingInline(img1ref)
	img1_inl.SetSize(12.5*measurement.Centimeter, 15.42*measurement.Centimeter)

	doc.AddParagraph().AddRun().AddPageBreak()

	addHead(doc, "巡检点位列表", "2")

	table_xjd := doc.AddTable()

	table_xjd.Properties().SetWidthPercent(100)
	// set borers
	borders := table_xjd.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Black, 1*measurement.Point)

	// 表头
	row_xjd_bt := table_xjd.AddRow()

	xjd_bt_slice := []string{
		"序号", "预置点位名称", "摄像头机位",
	}

	for _, s := range xjd_bt_slice {
		cell := row_xjd_bt.AddCell()
		cell.Properties().SetShading(wml.ST_ShdSolid, color.LightGray, color.Auto)
		p := cell.AddParagraph()
		p.Properties().SetAlignment(wml.ST_JcCenter)
		r := p.AddRun()
		r.Properties().SetBold(true)
		r.Properties().SetSize(11)
		r.AddText(s)
	}

	row_xjd_jg := table_xjd.AddRow()
	cell_xjd := row_xjd_jg.AddCell()
	cell_xjd.Properties().SetColumnSpan(3)
	cell_xjd.Properties().SetShading(wml.ST_ShdSolid, color.LightGray, color.Auto)
	cell_xjd.Properties().SetWidthPercent(100)
	cell_xjd_p := cell_xjd.AddParagraph()
	cell_xjd_p.Properties().SetAlignment(wml.ST_JcCenter)
	cell_xjd_r := cell_xjd_p.AddRun()
	cell_xjd_r.AddText("1号主变本体间隔")

	xjd_slice := []struct {
		Name   string
		Camera string
	}{
		{"1号主变本体全景1", "110kV双苏站110kV开关场东5号云台"},
		{"1号主变本体全景2", "110kV双苏站1号主变前侧6号云台"},
		{"1号主变本体全景3", "110kV双苏站主控楼顶西北12号球机"},
		{"1号主变本体全景顶部", "110kV双苏站主控楼顶东北13号球机"},
	}

	for i, s := range xjd_slice {
		row := table_xjd.AddRow()
		cell := row.AddCell()
		p := cell.AddParagraph()
		r := p.AddRun()
		r.Properties().SetSize(11)
		r.AddText(strconv.Itoa(i + 1))

		cell = row.AddCell()
		p = cell.AddParagraph()
		r = p.AddRun()
		r.Properties().SetSize(11)
		r.AddText(s.Name)

		cell = row.AddCell()
		p = cell.AddParagraph()
		r = p.AddRun()
		r.Properties().SetSize(11)
		r.AddText(s.Camera)

	}

	doc.AddParagraph().AddRun().AddPageBreak()

	addHead(doc, "巡检点位明细", "2")
	addHead(doc, "号主变本体间隔", "3")
	addHead(doc, "号主变本体全景1", "4")

	table_dw := doc.AddTable()
	table_dw.Properties().SetWidthPercent(100)
	// set borers
	table_dw_borders := table_dw.Properties().Borders()
	table_dw_borders.SetAll(wml.ST_BorderSingle, color.Black, 1*measurement.Point)

	row_dw_bt := table_dw.AddRow()
	dw_bt_slice := []string{"点位截图", "摄像头参数"}

	for _, s := range dw_bt_slice {
		cell := row_dw_bt.AddCell()
		cell.Properties().SetShading(wml.ST_ShdSolid, color.LightGray, color.Auto)
		p := cell.AddParagraph()
		p.Properties().SetAlignment(wml.ST_JcLeft)
		r := p.AddRun()
		r.Properties().SetBold(true)
		r.Properties().SetSize(11)
		r.AddText(s)
	}

	content_r := table_dw.AddRow()
	img_cell := content_r.AddCell()

	img_dw, _ := common.ImageFromFile("/Users/tonnamajesty/Downloads/dianwei.png")
	img_dw_ref, _ := doc.AddImage(img_dw)

	img_dw_inl, _ := img_cell.AddParagraph().AddRun().AddDrawingInline(img_dw_ref)
	img_dw_inl.SetSize(12*measurement.Centimeter, 9*measurement.Centimeter)

	cs_cell := content_r.AddCell()

	cs_cell_slice := map[string]string{
		"所用摄像头": "110kV双苏站110kV开关场东5号云台",
		"相机参数":  "",
		"图像调节":  "亮度:100    对比度:50   饱和度:50    锐度:50",
		"曝光模式":  "自动",
		"最大光圈":  "50",
		"最小光圈":  "23",
		"最大快门":  "1\\25",
		"聚焦模式":  "自动",
	}
	for k, v := range cs_cell_slice {
		cs_cell_k_p := cs_cell.AddParagraph()
		cs_cell_k_r := cs_cell_k_p.AddRun()
		cs_cell_k_r.Properties().SetBold(true)
		cs_cell_k_r.AddText(k + ":")

		cs_cell_v_p := cs_cell.AddParagraph()
		cs_cell_v_r := cs_cell_v_p.AddRun()
		cs_cell_v_r.AddText(v)
	}

	doc.SaveToFile("simple.docx")

}

func addHead(doc *document.Document, content string, style string) {
	head_p := doc.AddParagraph()
	head_p.SetStyle(style)
	head := head_p.AddRun()
	head.AddText(content)
	head.Properties().SetBold(true)
}

func addUl(doc *document.Document, content []string) {
	for _, c := range content {
		p := doc.AddParagraph()
		p.SetNumberingDefinitionByID(3)
		r := p.AddRun()
		r.AddText(c)
	}
}
