package main

import (
	"bytes"
	"fmt"
	"mime/multipart"

	"github.com/go-courier/httptransport/transformers"
)

type a struct {
	b []string
}

var points = []string{
	"JSLAY011122170003", "JSLAY011122170003", "JSLAY011122170003", "JSLAY011122170003", "JSLAY011122170003", "JSLAY011122170003", "JSLAY011122170003", "JSLAY011122170003", "JSLAY011122170003", "JSLAY011123730007", "JSLAY011123730007", "JSLAY011123730007", "JSLAY011123730007", "JSLAY011123730007", "JSLAY011123730007", "JSLAY011123730007", "JSLAY011123730007", "JSLAY011123730007", "JSLAY011124150008", "JSLAY011124150008", "JSLAY011124150008", "JSLAY011124150008", "JSLAY011124150008", "JSLAY011124150008", "JSLAY011124150008", "JSLAY011124150008", "JSLAY011124150008", "JSLAY011125300011", "JSLAY011125300011", "JSLAY011125300011", "JSLAY011125300011", "JSLAY011125300011", "JSLAY011125300011", "JSLAY011125300011", "JSLAY011125300011", "JSLAY011125300011", "JSLAY011127790017", "JSLAY011127790017", "JSLAY011127790017", "JSLAY011127790017", "JSLAY011127790017", "JSLAY011127790017", "JSLAY011127790017", "JSLAY011127790017", "JSLAY011127790017", "JSLAY011128520019", "JSLAY011128520019", "JSLAY011128520019", "JSLAY011128520019", "JSLAY011128520019", "JSLAY011128520019", "JSLAY011128520019", "JSLAY011128520019", "JSLAY011128520019", "JSLAY011129770022", "JSLAY011129770022", "JSLAY011129770022", "JSLAY011129770022", "JSLAY011129770022", "JSLAY011129770022", "JSLAY011129770022", "JSLAY011129770022", "JSLAY011129770022", "JSLAY011130170023", "JSLAY011130170023", "JSLAY011130170023", "JSLAY011130170023", "JSLAY011130170023", "JSLAY011130170023", "JSLAY011130170023", "JSLAY011130170023", "JSLAY011130170023", "JSLAY011130900025", "JSLAY011130900025", "JSLAY011130900025", "JSLAY011130900025", "JSLAY011130900025", "JSLAY011130900025", "JSLAY011130900025", "JSLAY011130900025", "JSLAY011130900025", "JSLAY011132560029", "JSLAY011132560029", "JSLAY011133710032", "JSLAY011133710032", "JSLAY011133710032", "JSLAY011133710032", "JSLAY011133710032", "JSLAY011133710032", "JSLAY011133710032", "JSLAY011133710032", "JSLAY011133710032", "JSLAY011134540034", "JSLAY011134540034", "JSLAY011134540034", "JSLAY011134540034", "JSLAY011134540034", "JSLAY011134540034", "JSLAY011134540034", "JSLAY011134540034", "JSLAY011134540034", "JSLAY011135380036", "JSLAY011135380036", "JSLAY011135380036", "JSLAY011135380036", "JSLAY011135380036", "JSLAY011135380036", "JSLAY011135380036", "JSLAY011135380036", "JSLAY011135380036", "JSLAY011136940040", "JSLAY011136940040", "JSLAY011136940040", "JSLAY011136940040", "JSLAY011136940040", "JSLAY011136940040", "JSLAY011136940040", "JSLAY011136940040", "JSLAY011136940040", "JSLAY011137360041", "JSLAY011137360041", "JSLAY011137360041", "JSLAY011137360041", "JSLAY011137360041", "JSLAY011137360041", "JSLAY011137360041", "JSLAY011137360041", "JSLAY011137360041", "JSLAY011140980050", "JSLAY011142130053", "JSLAY011142130053", "JSLAY011142130053", "JSLAY011142130053", "JSLAY011142130053", "JSLAY011142130053", "JSLAY011142130053", "JSLAY011142130053", "JSLAY011142130053", "JSLAY011142540054", "JSLAY011142540054", "JSLAY011142540054", "JSLAY011142540054", "JSLAY011142540054", "JSLAY011142540054", "JSLAY011142540054", "JSLAY011142540054", "JSLAY011142540054", "JSLAY011144110058", "JSLAY011144110058", "JSLAY011144110058", "JSLAY011144110058", "JSLAY011144110058", "JSLAY011144110058", "JSLAY011144110058", "JSLAY011144110058", "JSLAY011144110058", "JSLAY011145360061", "JSLAY011145360061", "JSLAY011145360061", "JSLAY011145360061", "JSLAY011145360061", "JSLAY011145360061", "JSLAY011145360061", "JSLAY011145360061", "JSLAY011145360061", "JSLAY011146920065", "JSLAY011146920065", "JSLAY011146920065", "JSLAY011146920065", "JSLAY011146920065", "JSLAY011146920065", "JSLAY011146920065", "JSLAY011146920065", "JSLAY011146920065", "JSLAY011147340066", "JSLAY011147340066", "JSLAY011147340066", "JSLAY011147340066", "JSLAY011147340066", "JSLAY011147340066", "JSLAY011147340066", "JSLAY011147340066", "JSLAY011147340066", "JSLAY011148900070", "JSLAY011148900070", "JSLAY011148900070", "JSLAY011148900070", "JSLAY011148900070", "JSLAY011148900070", "JSLAY011148900070", "JSLAY011148900070", "JSLAY011148900070", "JSLAY011150950075", "JSLAY011150950075", "JSLAY011150950075", "JSLAY011150950075", "JSLAY011150950075", "JSLAY011150950075", "JSLAY011150950075", "JSLAY011150950075", "JSLAY011150950075", "JSLAY011156570089", "JSLAY011156570089", "JSLAY011156570089", "JSLAY011156570089", "JSLAY011156570089", "JSLAY011156570089", "JSLAY011156570089", "JSLAY011156570089", "JSLAY011156570089", "JSLAY011159700097", "JSLAY011159700097", "JSLAY011159700097", "JSLAY011159700097", "JSLAY011159700097", "JSLAY011159700097", "JSLAY011159700097", "JSLAY011159700097", "JSLAY011159700097", "JSLAY011161760102", "JSLAY011161760102", "JSLAY011161760102", "JSLAY011161760102", "JSLAY011161760102", "JSLAY011161760102", "JSLAY011161760102", "JSLAY011161760102", "JSLAY011161760102", "JSLAY011162590104", "JSLAY011162590104", "JSLAY011162590104", "JSLAY011162590104", "JSLAY011162590104", "JSLAY011162590104", "JSLAY011162590104", "JSLAY011162590104", "JSLAY011162590104", "JSLAY011162910105", "JSLAY011163330106", "JSLAY011163330106", "JSLAY011163330106", "JSLAY011163330106", "JSLAY011163330106", "JSLAY011163330106", "JSLAY011163330106", "JSLAY011163330106", "JSLAY011163330106", "JSLAY011164570109", "JSLAY011164570109", "JSLAY011164570109", "JSLAY011164570109", "JSLAY011164570109", "JSLAY011164570109", "JSLAY011164570109", "JSLAY011164570109", "JSLAY011164570109", "JSLAY011166550114", "JSLAY011166550114", "JSLAY011166550114", "JSLAY011166550114", "JSLAY011166550114", "JSLAY011166550114", "JSLAY011166550114", "JSLAY011166550114", "JSLAY011166550114", "JSLAY011167700117", "JSLAY011167700117", "JSLAY011167700117", "JSLAY011167700117", "JSLAY011167700117", "JSLAY011167700117", "JSLAY011167700117", "JSLAY011167700117", "JSLAY011167700117", "JSLAY011168950120", "JSLAY011168950120", "JSLAY011168950120", "JSLAY011168950120", "JSLAY011168950120", "JSLAY011168950120", "JSLAY011168950120", "JSLAY011168950120", "JSLAY011168950120", "JSLAY011169780122", "JSLAY011169780122", "JSLAY011169780122", "JSLAY011169780122", "JSLAY011169780122", "JSLAY011169780122", "JSLAY011169780122", "JSLAY011169780122", "JSLAY011169780122", "JSLAY011171330126", "JSLAY011171330126", "JSLAY011171330126", "JSLAY011171330126", "JSLAY011171330126", "JSLAY011171330126", "JSLAY011171330126", "JSLAY011171330126", "JSLAY011171330126", "JSLAY011171740127", "JSLAY011171740127", "JSLAY011171740127", "JSLAY011171740127", "JSLAY011171740127", "JSLAY011171740127", "JSLAY011171740127", "JSLAY011171740127", "JSLAY011171740127", "JSLAY011177370141", "JSLAY011177370141", "JSLAY011177370141", "JSLAY011177370141", "JSLAY011177370141", "JSLAY011177370141", "JSLAY011177370141", "JSLAY011177370141", "JSLAY011177370141", "JSLAY011178100143", "JSLAY011178100143", "JSLAY011178100143", "JSLAY011178100143", "JSLAY011178100143", "JSLAY011178100143", "JSLAY011178100143", "JSLAY011178100143", "JSLAY011178100143", "JSLJS0427V5911594", "JSLJS0427V5911594", "JSLSJ091464530045", "JSLSJ091464530045", "JSLSJ091464530045", "JSLSJ091466100049", "JSLSJ091466100049", "JSLZT022488160003", "JSLZT022488160003", "JSLZT022488160003", "JSLZT022488160003", "JSLZT022488160003", "JSLZT022488160003", "JSLZT022488160003", "JSLZT022488160003", "JSLZT022488160003", "JSLZT022488990005", "JSLZT022488990005", "JSLZT022488990005", "JSLZT022488990005", "JSLZT022488990005", "JSLZT022488990005", "JSLZT022488990005", "JSLZT022488990005", "JSLZT022488990005", "JSLAY011149320071", "JSLAY011149320071", "JSLAY011149320071", "JSLAY011149320071", "JSLAY011149320071", "JSLAY011149320071", "JSLAY011149320071", "JSLAY011149320071", "JSLAY011149320071", "JSLAY011150120073", "JSLAY011150120073", "JSLAY011150120073", "JSLAY011150120073", "JSLAY011150120073", "JSLAY011150120073", "JSLAY011150120073", "JSLAY011150120073", "JSLAY011150120073", "JSLAY011160100098", "JSLAY011160100098", "JSLJS0427V5161537", "JSLJS0427V5161537", "JSLJS0427V5161537", "JSLJS0427V5161537", "JSLJS0427V5161537", "JSLJS0427V5161537", "JSLJS0427V5161537", "JSLJS0427V5161537", "JSLJS0427V5161537", "JSLJS0427V5311600", "JSLJS0427V5311600", "JSLJS0427V5311600", "JSLJS0427V5311600", "JSLJS0427V5311600", "JSLJS0427V5311600", "JSLJS0427V5311600", "JSLJS0427V5311600", "JSLJS0427V5311600", "JSLJS0427V5331595", "JSLJS0427V5331595", "JSLJS0427V5331595", "JSLJS0427V5331595", "JSLJS0427V5331595", "JSLJS0427V5331595", "JSLJS0427V5331595", "JSLJS0427V5331595", "JSLJS0427V5331595", "JSLJS0427V5761531", "JSLJS0427V5761531", "JSLJS0427V5761531", "JSLJS0427V5761531", "JSLJS0427V5761531", "JSLJS0427V5761531", "JSLJS0427V5761531", "JSLJS0427V5761531", "JSLJS0427V5761531", "JSLSJ091405700148", "JSLSJ091405700148", "JSLSJ091417760178", "JSLSJ091417760178", "JSLSJ091417760178", "JSLSJ091426910201", "JSLSJ091426910201", "JSLSJ091426910201", "JSLSJ091430530210", "JSLSJ091430530210", "JSLSJ091442180239", "JSLSJ091442180239", "JSLSJ091457770278", "JSLSJ091457770278", "JSLSJ091457770278", "JSJINSUG201842150", "JSJINSUG201842260", "JSLAY011121340001", "JSLAY011121340001", "JSLAY011126540014", "JSLAY011126540014", "JSLAY011128940020", "JSLAY011128940020", "JSLAY011131320026", "JSLAY011131320026", "JSLAY011131320026", "JSLAY011131320026", "JSLAY011131320026", "JSLAY011131320026", "JSLAY011131320026", "JSLAY011131320026", "JSLAY011131320026", "JSLAY011134960035", "JSLAY011134960035", "JSLAY011134960035", "JSLAY011134960035", "JSLAY011134960035", "JSLAY011134960035", "JSLAY011134960035", "JSLAY011134960035", "JSLAY011134960035", "JSLAY011145770062", "JSLAY011145770062", "JSLAY011145770062", "JSLAY011145770062", "JSLAY011145770062", "JSLAY011145770062", "JSLAY011145770062", "JSLAY011145770062", "JSLAY011145770062", "JSLAY011152510079", "JSLAY011152510079", "JSLAY011152510079", "JSLAY011152510079", "JSLAY011152510079", "JSLAY011152510079", "JSLAY011152510079", "JSLAY011152510079", "JSLAY011152510079", "JSLAY011152930080", "JSLAY011152930080", "JSLAY011152930080", "JSLAY011152930080", "JSLAY011152930080", "JSLAY011152930080", "JSLAY011152930080", "JSLAY011152930080", "JSLAY011152930080", "JSLAY011157310091", "JSLAY011157310091", "JSLAY011158140093", "JSLAY011158140093", "JSLAY011166140113", "JSLAY011166140113", "JSLAY011176950140", "JSLAY011176950140", "JSLAY011178930145", "JSLAY011178930145", "JSLJS0427V5121542", "JSLJS0427V5121542", "JSLJS0427V5121542", "JSLJS0427V5121542", "JSLJS0427V5121542", "JSLJS0427V5121542", "JSLJS0427V5121542", "JSLJS0427V5121542", "JSLJS0427V5121542", "JSLSJ091404140144", "JSLSJ091404140144", "JSLSJ091404140144", "JSLSJ091404970146", "JSLSJ091404970146", "JSLSJ091404970146", "JSLSJ091428990206", "JSLSJ091428990206", "JSLSJ091428990206", "JSLSJ091469330057", "JSLSJ091469330057", "JSLSJ091482930091", "JSLSJ091482930091", "JSLSJ091489390107", "JSLSJ091489390107", "JSLSJ091489390107", "JSLSJ091497390127", "JSLSJ091497390127", "JSLSJ091497390127", "JSLSW031967720001", "JSLSW031968140002", "JSLSW3322V5C16006", "JSLSW3322V5C16006", "JSLSW3322V5C18001", "JSLSW3322V5C18001", "JSLSW3322V5C33004", "JSLSW3322V5C33004", "JSLSW3322V5C59002", "JSLSW3322V5C59002", "JSLSW3322V5C74005", "JSLSW3322V5C74005", "JSLSW3322V5C91003", "JSLSW3322V5C91003", "JSLZT022489310006", "JSLZT022489310006", "JSLZT022489310006", "JSLZT022489310006", "JSLZT022489310006", "JSLZT022489310006", "JSLZT022489310006", "JSLZT022489310006", "JSLZT022489310006", "JSLAY031514940067", "JSLAY031514940067", "JSLAY031514940067", "JSLJS0427V5571578", "JSLJS0427V5571578", "JSLSW3407V5C34003", "JSLSW3407V5C34003", "JSLSW3407V5C50001", "JSLSW3407V5C50001", "JSLSW3407V5C92002", "JSLSW3407V5C92002", "JSLYY3214V5C35054", "JSLYY3214V5C35054", "JSLYY3214V5C37094", "JSLYY3214V5C37094", "BRSIPC4GJS1780175", "BRSIPC4GJS1780192", "BRSIPC4GJS1780193", "BRSIPC4GJS1780196", "JDNY2012173418008", "JDNY2012173418008", "JDNY2012173491010", "JDNY2012173491010", "JDNY2012173699015", "JDNY2012173699015", "JDNY2012173814018", "JDNY2012173814018", "JSLAY031590500006", "JSLAY031590920007", "JSLAY031591750009", "JSLJS0427V5012637", "JSLJS0427V5012637", "JSLJS0427V5012657", "JSLJS0427V5012657", "JSLJS0427V5013777", "JSLJS0427V5013777", "JSLJS0427V5014672", "JSLJS0427V5014672", "JSLJS0427V5015752", "JSLJS0427V5015752", "JSLJS0427V5016627", "JSLJS0427V5016627", "JSLJS0427V5018707", "JSLJS0427V5018707", "JSLJS0427V5019742", "JSLJS0427V5019742", "JSLJS0427V5019827", "JSLJS0427V5019827", "JSLJS0427V5019912", "JSLJS0427V5019912", "JSLJS0427V5030755", "JSLJS0427V5030755", "JSLJS0427V5030840", "JSLJS0427V5030840", "JSLJS0427V5031735", "JSLJS0427V5031735", "JSLJS0427V5033605", "JSLJS0427V5033605", "JSLJS0427V5033710", "JSLJS0427V5033710", "JSLJS0427V5033730", "JSLJS0427V5033730", "JSLJS0427V5034745", "JSLJS0427V5034745", "JSLJS0427V5035620", "JSLJS0427V5035620", "JSLJS0427V5035685", "JSLJS0427V5035685", "JSLJS0427V5036910", "JSLJS0427V5036910", "JSLJS0427V5037640", "JSLJS0427V5037640", "JSLJS0427V5037660", "JSLJS0427V5037660", "JSLJS0427V5039655", "JSLJS0427V5039655", "JSLJS0427V5039675", "JSLJS0427V5039675", "JSLJS0427V5051723", "JSLJS0427V5051723", "JSLJS0427V5052803", "JSLJS0427V5052803", "JSLJS0427V5054818", "JSLJS0427V5054818", "JSLJS0427V5054883", "JSLJS0427V5054883", "JSLJS0427V5055613", "JSLJS0427V5055613", "JSLJS0427V5056773", "JSLJS0427V5056773", "JSLJS0427V5057628", "JSLJS0427V5057628", "JSLJS0427V5070616", "JSLJS0427V5070616", "JSLJS0427V5071861", "JSLJS0427V5071861", "JSLJS0427V5072631", "JSLJS0427V5072631", "JSLJS0427V5073856", "JSLJS0427V5073856", "JSLJS0427V5075851", "JSLJS0427V5075851", "JSLJS0427V5075871", "JSLJS0427V5075871", "JSLJS0427V5076621", "JSLJS0427V5076621", "JSLJS0427V5076621", "JSLJS0427V5076621", "JSLJS0427V5076621", "JSLJS0427V5076621", "JSLJS0427V5076621", "JSLJS0427V5076621", "JSLJS0427V5076621", "JSLJS0427V5078681", "JSLJS0427V5078681", "JSLJS0427V5092849", "JSLJS0427V5092849", "JSLJS0427V5093664", "JSLJS0427V5093664", "JSLJS0427V5094909", "JSLJS0427V5094909", "JSLJS0427V5097654", "JSLJS0427V5097654", "JSLJS0427V5098754", "JSLJS0427V5098754", "JSLJS0427V5098754", "JSLJS0427V5098774", "JSLJS0427V5098774", "JSLJS0427V5098839", "JSLJS0427V5098839", "JSLJS0427V5181552", "JSLJS0427V5181552", "JSLJS0427V5311000", "JSLJS0427V5311000", "JSLJS0427V5701566", "JSLJS0427V5701566", "JSLJS0427V5951589", "JSLJS0427V5951589", "JSLSJ1126V6S11019", "JSLSJ1126V6S11019", "JSLSJ1126V6S13099", "JSLSJ1126V6S13099", "JSLSJ1126V6S15034", "JSLSJ1126V6S15034", "JSLSJ1126V6S15054", "JSLSJ1126V6S15054", "JSLSJ1126V6S17029", "JSLSJ1126V6S17029", "JSLSJ1126V6S17069", "JSLSJ1126V6S17069", "JSLSJ1126V6S30037", "JSLSJ1126V6S30037", "JSLSJ1126V6S30057", "JSLSJ1126V6S30057", "JSLSJ1126V6S32052", "JSLSJ1126V6S32052", "JSLSJ1126V6S36022", "JSLSJ1126V6S36022", "JSLSJ1126V6S36042", "JSLSJ1126V6S36042", "JSLSJ1126V6S38062", "JSLSJ1126V6S38062", "JSLSJ1126V6S38082", "JSLSJ1126V6S38082", "JSLSJ1126V6S38102", "JSLSJ1126V6S38102", "JSLSJ1126V6S50065", "JSLSJ1126V6S50065", "JSLSJ1126V6S52085", "JSLSJ1126V6S52085", "JSLSJ1126V6S56055", "JSLSJ1126V6S56055", "JSLSJ1126V6S58030", "JSLSJ1126V6S58030", "JSLSJ1126V6S71038", "JSLSJ1126V6S71038", "JSLSJ1126V6S71098", "JSLSJ1126V6S71098", "JSLSJ1126V6S73033", "JSLSJ1126V6S73033", "JSLSJ1126V6S73053", "JSLSJ1126V6S73053", "JSLSJ1126V6S79063", "JSLSJ1126V6S79063", "JSLSJ1126V6S79103", "JSLSJ1126V6S79103", "JSLSJ1126V6S90031", "JSLSJ1126V6S90031", "JSLSJ1126V6S90051", "JSLSJ1126V6S90051", "JSLSJ1126V6S94041", "JSLSJ1126V6S94041", "JSLSJ1126V6S94086", "JSLSJ1126V6S94086", "JSLSJ1126V6S96101", "JSLSJ1126V6S96101", "JSLSJ1126V6S98036", "JSLSJ1126V6S98036", "JSLYY3214V5C12111", "JSLYY3214V5C12111", "JSLYY3214V5C14126", "JSLYY3214V5C14126", "JSLYY3214V5C14171", "JSLYY3214V5C14171", "JSLYY3214V5C31169", "JSLYY3214V5C31169", "JSLYY3214V5C51157", "JSLYY3214V5C51157", "JSLYY3214V5C53132", "JSLYY3214V5C53132", "JSLYY3214V5C53152", "JSLYY3214V5C53152", "JSLYY3214V5C93118", "JSLYY3214V5C93118", "JSLYY3214V5C97173", "JSLYY3214V5C97173", "V020020HY00008230", "V020020HY00008230", "V020020HY00008231", "V020020HY00008231", "V020020HY00008232", "V020020HY00008232", "V020020HY00008233", "V020020HY00008233", "ZYPHOTO0000002119", "ZYPHOTO0000002154", "ZYPHOTO0000002182", "ZYPHOTO0000002184", "ZYPHOTO0000002186", "ZYPHOTO0000002189", "ZYPHOTO0000002191", "ZYPHOTO0000002193", "ZYPHOTO0000002194", "ZYPHOTO0000002196", "ZYPHOTO0000002198", "ZYPHOTO0000002199", "ZYPHOTO0000002204", "ZYPHOTO0000002210", "ZYPHOTO0000002211", "ZYPHOTO0000002212", "ZYPHOTO0000002213", "ZYPHOTO0000002214", "ZYPHOTO0000002215", "ZYPHOTO0000002216", "ZYPHOTO0000002217", "ZYPHOTO0000002218", "ZYPHOTO0000002220", "ZYPHOTO0000002222", "ZYPHOTO0000002223", "ZYPHOTO0000002225", "ZYPHOTO0000002227", "ZYPHOTO0000002231", "ZYPHOTO0000002232", "ZYPHOTO0000002234", "ZYPHOTO0000002235", "ZYPHOTO0000002237", "ZYPHOTO0000002242", "ZYPHOTO0000002245", "ZYPHOTO0000002250", "ZYPHOTO0000002252", "ZYPHOTO0000002261", "ZYPHOTO0000002264", "ZYPHOTO0000002265", "ZYPHOTO0000002266", "ZYPHOTO0000002268", "ZYPHOTO0000002271", "ZYPHOTO0000002272", "ZYPHOTO0000002277", "ZYPHOTO0000002281", "ZYPHOTO0000002282", "ZYPHOTO0000002286", "ZYPHOTO0000002287", "ZYPHOTO0000002288", "ZYPHOTO0000002290", "ZYPHOTO0000002291", "ZYPHOTO0000002313", "ZYPHOTO0000002315", "ZYPHOTO0000002321", "ZYPHOTO0000002322", "ZYPHOTO0000002325", "ZYPHOTO0000002326", "ZYPHOTO0000002327", "ZYPHOTO0000002328", "ZYPHOTO0000002330", "ZYPHOTO0000002331", "ZYPHOTO0000002334", "ZYPHOTO0000002335", "ZYPHOTO0000002336", "ZYPHOTO0000002338", "ZYPHOTO0000002339", "ZYPHOTO0000002340", "ZYPHOTO0000002345", "ZYPHOTO0000002348", "ZYPHOTO0000002351", "ZYPHOTO0000002352", "ZYPHOTO0000002353", "ZYPHOTO0000002354", "ZYPHOTO0000002355", "ZYPHOTO0000002371", "ZYPHOTO0000002372", "ZYPHOTO0000002374", "ZYPHOTO0000002375", "ZYPHOTO0000002376", "ZYPHOTO0000002378", "ZYPHOTO0000002396", "ZYPHOTO0000003040", "ZYPHOTO0000003134", "ZYPHOTO0000003136", "ZYPHOTO0000003151", "ZYPHOTO0000003156", "ZYPHOTO0000003161", "ZYPHOTO0000003162", "ZYPHOTO0000003165", "ZYPHOTO0000003167", "ZYPHOTO0000003174", "ZYPHOTO0000003178", "ZYPHOTO0000003180", "ZYPHOTO0000003183", "ZYPHOTO0000003187", "ZYPHOTO0000003191", "ZYPHOTO0000003192", "ZYPHOTO0000003195", "ZYPHOTO0000003196", "ZYPHOTO0000003197",
}

var points2 = []string{
	"BRSIPC4GJS1780175", "BRSIPC4GJS1780192", "BRSIPC4GJS1780193", "BRSIPC4GJS1780196", "JDNY2012173418008", "JDNY2012173491010", "JDNY2012173699015", "JDNY2012173814018", "JSJINSUG201842150", "JSJINSUG201842260", "JSLAY01112217000", "JSLAY01112373000", "JSLAY01112415000", "JSLAY01112530001", "JSLAY01112977002", "JSLAY01113090002", "JSLAY011121340001", "JSLAY011122170003", "JSLAY011123730007", "JSLAY011124150008", "JSLAY011125300011", "JSLAY011126540014", "JSLAY011127790017", "JSLAY011128520019", "JSLAY011128940020", "JSLAY011129770022", "JSLAY011130170023", "JSLAY011130900025", "JSLAY011131320026", "JSLAY011132560029", "JSLAY011133710032", "JSLAY011134540034", "JSLAY011134960035", "JSLAY011135380036", "JSLAY011136940040", "JSLAY011137360041", "JSLAY011140980050", "JSLAY011142130053", "JSLAY011142540054", "JSLAY011144110058", "JSLAY011145360061", "JSLAY011145770062", "JSLAY011146920065", "JSLAY011147340066", "JSLAY011148900070", "JSLAY011149320071", "JSLAY011150120073", "JSLAY011150950075", "JSLAY011152510079", "JSLAY011152930080", "JSLAY011156570089", "JSLAY011157310091", "JSLAY011158140093", "JSLAY011159700097", "JSLAY011160100098", "JSLAY011161760102", "JSLAY011162590104", "JSLAY011162910105", "JSLAY011163330106", "JSLAY011164570109", "JSLAY011166140113", "JSLAY011166550114", "JSLAY011167700117", "JSLAY011168950120", "JSLAY011169780122", "JSLAY011171330126", "JSLAY011171740127", "JSLAY011176950140", "JSLAY011177370141", "JSLAY011178100143", "JSLAY011178930145", "JSLAY031514940067", "JSLAY031590500006", "JSLAY031590920007", "JSLAY031591750009", "JSLJS0427V5012637", "JSLJS0427V5012657", "JSLJS0427V5013777", "JSLJS0427V5013777", "JSLJS0427V5014672", "JSLJS0427V5015752", "JSLJS0427V5015752", "JSLJS0427V5016627", "JSLJS0427V5018707", "JSLJS0427V5019742", "JSLJS0427V5019742", "JSLJS0427V5019827", "JSLJS0427V5019827", "JSLJS0427V5019912", "JSLJS0427V5030755", "JSLJS0427V5030755", "JSLJS0427V5030840", "JSLJS0427V5031735", "JSLJS0427V5031735", "JSLJS0427V5033605", "JSLJS0427V5033710", "JSLJS0427V5033730", "JSLJS0427V5033730", "JSLJS0427V5034745", "JSLJS0427V5034745", "JSLJS0427V5035620", "JSLJS0427V5035685", "JSLJS0427V5036910", "JSLJS0427V5036910", "JSLJS0427V5037640", "JSLJS0427V5037660", "JSLJS0427V5039655", "JSLJS0427V5039675", "JSLJS0427V5051723", "JSLJS0427V5052803", "JSLJS0427V5054818", "JSLJS0427V5054883", "JSLJS0427V5054883", "JSLJS0427V5055613", "JSLJS0427V5056773", "JSLJS0427V5057628", "JSLJS0427V5070616", "JSLJS0427V5071861", "JSLJS0427V5072631", "JSLJS0427V5073856", "JSLJS0427V5073856", "JSLJS0427V5075851", "JSLJS0427V5075871", "JSLJS0427V5076621", "JSLJS0427V5078681", "JSLJS0427V5092849", "JSLJS0427V5093664", "JSLJS0427V5094909", "JSLJS0427V5097654", "JSLJS0427V5098754", "JSLJS0427V5098774", "JSLJS0427V5098839", "JSLJS0427V5121542", "JSLJS0427V5161537", "JSLJS0427V5181552", "JSLJS0427V5311000", "JSLJS0427V5311600", "JSLJS0427V5331595", "JSLJS0427V5571578", "JSLJS0427V5701566", "JSLJS0427V5761531", "JSLJS0427V5911594", "JSLJS0427V5951589", "JSLSJ1126V6S11019", "JSLSJ1126V6S13099", "JSLSJ1126V6S15034", "JSLSJ1126V6S15054", "JSLSJ1126V6S17029", "JSLSJ1126V6S17069", "JSLSJ1126V6S30037", "JSLSJ1126V6S30057", "JSLSJ1126V6S32052", "JSLSJ1126V6S36022", "JSLSJ1126V6S36042", "JSLSJ1126V6S38062", "JSLSJ1126V6S38082", "JSLSJ1126V6S38102", "JSLSJ1126V6S50065", "JSLSJ1126V6S52085", "JSLSJ1126V6S56055", "JSLSJ1126V6S58030", "JSLSJ1126V6S71038", "JSLSJ1126V6S71098", "JSLSJ1126V6S73033", "JSLSJ1126V6S73053", "JSLSJ1126V6S79063", "JSLSJ1126V6S79103", "JSLSJ1126V6S90031", "JSLSJ1126V6S90051", "JSLSJ1126V6S94041", "JSLSJ1126V6S94086", "JSLSJ1126V6S96101", "JSLSJ1126V6S98036", "JSLSJ091404140144", "JSLSJ091404970146", "JSLSJ091405700148", "JSLSJ091417760178", "JSLSJ091426910201", "JSLSJ091428990206", "JSLSJ091428990206", "JSLSJ091430530210", "JSLSJ091430530210", "JSLSJ091442180239", "JSLSJ091464530045", "JSLSJ091466100049", "JSLSJ091469330057", "JSLSJ091482930091", "JSLSJ091489390107", "JSLSJ091497390127", "JSLSW3322V5C16006", "JSLSW3322V5C18001", "JSLSW3322V5C59002", "JSLSW3322V5C74005", "JSLSW3322V5C91003", "JSLSW3407V5C34003", "JSLSW3407V5C50001", "JSLSW031967720001", "JSLSW031968140002", "JSLYY3214V5C12111", "JSLYY3214V5C14126", "JSLYY3214V5C14171", "JSLYY3214V5C31169", "JSLYY3214V5C35054", "JSLYY3214V5C37094", "JSLYY3214V5C51157", "JSLYY3214V5C53132", "JSLYY3214V5C53152", "JSLYY3214V5C93118", "JSLYY3214V5C97173", "JSLZT022488160003", "JSLZT022488990005", "JSLZT022489310006", "V020020HY00008230", "V020020HY00008231", "V020020HY00008232", "V020020HY00008233", "ZYPHOTO0000002119", "ZYPHOTO0000002154", "ZYPHOTO0000002182", "ZYPHOTO0000002184", "ZYPHOTO0000002186", "ZYPHOTO0000002189", "ZYPHOTO0000002191", "ZYPHOTO0000002193", "ZYPHOTO0000002194", "ZYPHOTO0000002196", "ZYPHOTO0000002198", "ZYPHOTO0000002199", "ZYPHOTO0000002204", "ZYPHOTO0000002210", "ZYPHOTO0000002211", "ZYPHOTO0000002212", "ZYPHOTO0000002213", "ZYPHOTO0000002214", "ZYPHOTO0000002215", "ZYPHOTO0000002216", "ZYPHOTO0000002217", "ZYPHOTO0000002218", "ZYPHOTO0000002220", "ZYPHOTO0000002222", "ZYPHOTO0000002223", "ZYPHOTO0000002225", "ZYPHOTO0000002227", "ZYPHOTO0000002231", "ZYPHOTO0000002232", "ZYPHOTO0000002234", "ZYPHOTO0000002235", "ZYPHOTO0000002237", "ZYPHOTO0000002242", "ZYPHOTO0000002245", "ZYPHOTO0000002250", "ZYPHOTO0000002252", "ZYPHOTO0000002261", "ZYPHOTO0000002264", "ZYPHOTO0000002265", "ZYPHOTO0000002266", "ZYPHOTO0000002268", "ZYPHOTO0000002271", "ZYPHOTO0000002272", "ZYPHOTO0000002277", "ZYPHOTO0000002281", "ZYPHOTO0000002282", "ZYPHOTO0000002286", "ZYPHOTO0000002287", "ZYPHOTO0000002288", "ZYPHOTO0000002290", "ZYPHOTO0000002291", "ZYPHOTO0000002313", "ZYPHOTO0000002315", "ZYPHOTO0000002321", "ZYPHOTO0000002322", "ZYPHOTO0000002325", "ZYPHOTO0000002326", "ZYPHOTO0000002327", "ZYPHOTO0000002328", "ZYPHOTO0000002330", "ZYPHOTO0000002331", "ZYPHOTO0000002334", "ZYPHOTO0000002335", "ZYPHOTO0000002336", "ZYPHOTO0000002338", "ZYPHOTO0000002339", "ZYPHOTO0000002340", "ZYPHOTO0000002345", "ZYPHOTO0000002348", "ZYPHOTO0000002351", "ZYPHOTO0000002352", "ZYPHOTO0000002353", "ZYPHOTO0000002354", "ZYPHOTO0000002355", "ZYPHOTO0000002371", "ZYPHOTO0000002372", "ZYPHOTO0000002374", "ZYPHOTO0000002375", "ZYPHOTO0000002376", "ZYPHOTO0000002378", "ZYPHOTO0000002396", "ZYPHOTO0000003134", "ZYPHOTO0000003136", "ZYPHOTO0000003151", "ZYPHOTO0000003156", "ZYPHOTO0000003161", "ZYPHOTO0000003165", "ZYPHOTO0000003167", "ZYPHOTO0000003174", "ZYPHOTO0000003178", "ZYPHOTO0000003180", "ZYPHOTO0000003183", "ZYPHOTO0000003187", "ZYPHOTO0000003188", "ZYPHOTO0000003189", "ZYPHOTO0000003191", "ZYPHOTO0000003192", "ZYPHOTO0000003195", "ZYPHOTO0000003196", "ZYPHOTO0000003197",
}

type DelSampleRes struct {
	FailedImageIDs   []int `json:"failedImageIDs"`
	AffectedImageIDs []int `json:"affectedImageIDs"`
}

func (r *DelSampleRes) Format(s fmt.State, verb rune) {
	fmt.Fprintf(s, "%s", "hello")
}

type DetectDetails map[string][]DetectDetail

type DetectDetail struct {
	Color []uint8 `json:"color"` // b g r

}

func main() {

	// 假设编号为13，结果为1
	number := uint16(13)
	result := uint16(1)

	// 将编号左移2位（给结果预留2位），并加上结果
	finalValue := (number << 2) | result

	// 输出最终的uint16值
	fmt.Printf("Final uint16 value: %d\n", finalValue)

	// 提取编号和结果
	extractedNumber := (finalValue >> 2) & 0x3FFF
	extractedResult := finalValue & 0x3

	fmt.Printf("Extracted Number: %d\n", extractedNumber)
	fmt.Printf("Extracted Result: %d\n", extractedResult)

}

func BytesToFileHeader(fileData []byte) (*multipart.FileHeader, error) {
	reader := bytes.NewReader(fileData)
	// filename := uuid.New().String()
	filename := fmt.Sprintf("shm_%s_%d_%d_%d_%d", "shm-test", 30, 1000, 2, 1000)
	fileHeader, err := transformers.NewFileHeader(filename, filename, reader)
	if err != nil {
		return nil, err
	}

	return fileHeader, nil
}

type Entity struct {
	IsReport   bool       `json:"isreport"`
	Name       string     `json:"name"`
	Desc       string     `json:"desc"`
	ParentID   string     `json:"parent_id"`
	Points     [][]int    `json:"points"`
	Confidence float32    `json:"conf"`
	Alarm      bool       `json:"alarm"`
	TrackID    string     `json:"track_id"`
	Extral     string     `json:"extral"`
	Properties []Property `json:"property"`
}

type Property struct {
	IsReport            bool    `json:"isreport"`
	Name                string  `json:"name"`
	InnerName           string  `json:"inner_name"`
	Value               string  `json:"value"`
	Desc                string  `json:"desc"`
	Confidence          float32 `json:"conf"`
	Extral              string  `json:"extral"`
	VoteScore           float32 `json:"vote_score"`
	SingleAnalyzeResult float32 `json:"single_analyze_score"`
	Alarm               bool    `json:"alarm"`
}

type AnalysisResult struct {
	AnalysisDuration string         `json:"analysisDuration"`
	Duration         string         `json:"duration"`
	Alarm            bool           `json:"alarm"`
	EntityStruct     []Entity       `json:"entity_struct"`
	HistoryData      map[string]any `json:"history_data"`
}
