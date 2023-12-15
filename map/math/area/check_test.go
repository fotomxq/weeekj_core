package MapMathArea

import (
	"encoding/json"
	"testing"

	MapMathArgs "github.com/fotomxq/weeekj_core/v5/map/math/args"
)

func TestCheckXYInArea(t *testing.T) {
	b := CheckXYInArea(&ArgsCheckXYInArea{
		Point: MapMathArgs.ParamsPoint{
			PointType: "WGS-84",
			Longitude: 5.3,
			Latitude:  5.2,
		}, Area: ParamsArea{
			ID:        1,
			PointType: "GCJ-02",
			Points: []ParamsAreaPoint{
				//该范围大致为矩形结构，四个点组成
				{
					Longitude: 1.1,
					Latitude:  1.2,
				},
				{
					Longitude: 7.6,
					Latitude:  1.5,
				},
				{
					Longitude: 8.6,
					Latitude:  9.6,
				},
				{
					Longitude: 1.3,
					Latitude:  9.5,
				},
			},
		},
	})
	if !b {
		t.Error(b)
	}
	//超级复杂的结构体，求是否在范围内
	pointsByte := []byte("{\"points\":[{\"longitude\":112.615482198,\"latitude\":37.858581395908445},{\"longitude\":112.61642574884002,\"latitude\":37.86290457688844},{\"longitude\":112.61775583114479,\"latitude\":37.86682780510099},{\"longitude\":112.61790588816632,\"latitude\":37.868653828076475},{\"longitude\":112.61815257805404,\"latitude\":37.872480256901746},{\"longitude\":112.61855757327987,\"latitude\":37.87435148028163},{\"longitude\":112.61844758437485,\"latitude\":37.87612103180482},{\"longitude\":112.61838051081406,\"latitude\":37.87802603663723},{\"longitude\":112.61754096105699,\"latitude\":37.88010035724216},{\"longitude\":112.61329116851095,\"latitude\":37.88464567833344},{\"longitude\":112.61103752620522,\"latitude\":37.88718919184703},{\"longitude\":112.60981385216121,\"latitude\":37.89007130211149},{\"longitude\":112.60946935594086,\"latitude\":37.89244853315518},{\"longitude\":112.61021920099859,\"latitude\":37.89582473343393},{\"longitude\":112.61055120818321,\"latitude\":37.897504309266616},{\"longitude\":112.610733011663,\"latitude\":37.89909918624958},{\"longitude\":112.61052857704465,\"latitude\":37.90132050068013},{\"longitude\":112.60959458157424,\"latitude\":37.9030676759442},{\"longitude\":112.60481371164326,\"latitude\":37.90842845472225},{\"longitude\":112.60260776236657,\"latitude\":37.91460141773439},{\"longitude\":112.60126011997465,\"latitude\":37.91908113056412},{\"longitude\":112.59853918656711,\"latitude\":37.92410221611611},{\"longitude\":112.59723445951943,\"latitude\":37.93059543698284},{\"longitude\":112.59558640971784,\"latitude\":37.936478840506076},{\"longitude\":112.59325171440844,\"latitude\":37.94215870697221},{\"longitude\":112.5908741037548,\"latitude\":37.947838134477294},{\"longitude\":112.5887378919125,\"latitude\":37.95062804647866},{\"longitude\":112.58724541023378,\"latitude\":37.9535532093001},{\"longitude\":112.58601042062048,\"latitude\":37.95586917249172},{\"longitude\":112.585462076515,\"latitude\":37.95832041064279},{\"longitude\":112.58572612741966,\"latitude\":37.96010504097604},{\"longitude\":112.58616183970128,\"latitude\":37.961889627943265},{\"longitude\":112.58754824839536,\"latitude\":37.965594006351985},{\"longitude\":112.58867716502402,\"latitude\":37.969467357216274},{\"longitude\":112.58860445201401,\"latitude\":37.97361114358078},{\"longitude\":112.5894889942557,\"latitude\":37.98135679500386},{\"longitude\":112.59057499554012,\"latitude\":37.98536461229896},{\"longitude\":112.59243347302083,\"latitude\":37.989236919954564},{\"longitude\":112.5950215113536,\"latitude\":37.99395454161357},{\"longitude\":112.59821036450569,\"latitude\":37.998468949403815},{\"longitude\":112.60161379437898,\"latitude\":38.002543465106285},{\"longitude\":112.60404507016767,\"latitude\":38.00485116268588},{\"longitude\":112.60538945409473,\"latitude\":38.006072614175764},{\"longitude\":112.60580415399284,\"latitude\":38.006700239623456},{\"longitude\":112.60578970044855,\"latitude\":38.00759837418592},{\"longitude\":112.60389308530841,\"latitude\":38.00926321169799},{\"longitude\":112.6033331008349,\"latitude\":38.01010208005929},{\"longitude\":112.6031164391153,\"latitude\":38.011143814781896},{\"longitude\":112.60294269273993,\"latitude\":38.01308155074273},{\"longitude\":112.60242562361066,\"latitude\":38.015661638352},{\"longitude\":112.6016060620733,\"latitude\":38.017187074976285},{\"longitude\":112.60070066984747,\"latitude\":38.018712479851196},{\"longitude\":112.598203239888,\"latitude\":38.02101942085869},{\"longitude\":112.59472382806246,\"latitude\":38.02334481045479},{\"longitude\":112.5907723474503,\"latitude\":38.02523064808315},{\"longitude\":112.5868852398545,\"latitude\":38.02659245401012},{\"longitude\":112.58282647088174,\"latitude\":38.02751477015391},{\"longitude\":112.57932560138408,\"latitude\":38.027938454621456},{\"longitude\":112.57599639326338,\"latitude\":38.02815930833961},{\"longitude\":112.56895173892383,\"latitude\":38.028127748417084},{\"longitude\":112.5653653321043,\"latitude\":38.02825507743564},{\"longitude\":112.56177892528478,\"latitude\":38.02848382014379},{\"longitude\":112.55899718116973,\"latitude\":38.02869565988172},{\"longitude\":112.55617252171044,\"latitude\":38.02873847672171},{\"longitude\":112.55397193685178,\"latitude\":38.0284026758213},{\"longitude\":112.55161398176108,\"latitude\":38.02894833499687},{\"longitude\":112.54934185735885,\"latitude\":38.02905453487327},{\"longitude\":112.54642839163546,\"latitude\":38.029740192166685},{\"longitude\":112.54269953437154,\"latitude\":38.02980047160069},{\"longitude\":112.53897067710761,\"latitude\":38.029691730899046},{\"longitude\":112.53314911007885,\"latitude\":38.02923761953354},{\"longitude\":112.52741337373857,\"latitude\":38.02844545982023},{\"longitude\":112.52154889136557,\"latitude\":38.02761948653662},{\"longitude\":112.515941901058,\"latitude\":38.02632022690286},{\"longitude\":112.51123613297943,\"latitude\":38.02464062336292},{\"longitude\":112.50670202627782,\"latitude\":38.022690523526094},{\"longitude\":112.50298331111674,\"latitude\":38.02026705565201},{\"longitude\":112.49947917267684,\"latitude\":38.01726874348897},{\"longitude\":112.49514891386036,\"latitude\":38.01091447310927},{\"longitude\":112.49197736933831,\"latitude\":38.0038495197988},{\"longitude\":112.49054180078213,\"latitude\":37.99921770739519},{\"longitude\":112.48936372429137,\"latitude\":37.99438268066173},{\"longitude\":112.48769421681766,\"latitude\":37.98971757087527},{\"longitude\":112.48458704531197,\"latitude\":37.983724535103235},{\"longitude\":112.4813940431178,\"latitude\":37.97769718179354},{\"longitude\":112.47777188748125,\"latitude\":37.9721767938992},{\"longitude\":112.47389223977927,\"latitude\":37.96692665532358},{\"longitude\":112.47022350171585,\"latitude\":37.96352135444181},{\"longitude\":112.46818554673348,\"latitude\":37.95970986087338},{\"longitude\":112.46792831660252,\"latitude\":37.95828621844894},{\"longitude\":112.46818607060243,\"latitude\":37.95679487312628},{\"longitude\":112.46880860502836,\"latitude\":37.95506662866384},{\"longitude\":112.46988175056879,\"latitude\":37.95347370044938},{\"longitude\":112.47305800991137,\"latitude\":37.94989011150992},{\"longitude\":112.47453885122206,\"latitude\":37.948115172321266},{\"longitude\":112.47589094650004,\"latitude\":37.946238662734864},{\"longitude\":112.47755443995823,\"latitude\":37.942519344038786},{\"longitude\":112.4778068295401,\"latitude\":37.94035500757586},{\"longitude\":112.47784464240078,\"latitude\":37.93819060737211},{\"longitude\":112.47660219289367,\"latitude\":37.932501112589144},{\"longitude\":112.47602388348434,\"latitude\":37.92938539823615},{\"longitude\":112.47604638889436,\"latitude\":37.92579562616597},{\"longitude\":112.47588755182926,\"latitude\":37.922458026542756},{\"longitude\":112.47512789994482,\"latitude\":37.91878172449898},{\"longitude\":112.47300778135661,\"latitude\":37.91247817687452},{\"longitude\":112.47274375140671,\"latitude\":37.908586643944325},{\"longitude\":112.4724797214568,\"latitude\":37.904017668964045},{\"longitude\":112.47137884229426,\"latitude\":37.899956369569885},{\"longitude\":112.46950548693542,\"latitude\":37.89630123338189},{\"longitude\":112.46872601184998,\"latitude\":37.894652949794676},{\"longitude\":112.46828985951845,\"latitude\":37.89286916104128},{\"longitude\":112.4681755722687,\"latitude\":37.891237734531536},{\"longitude\":112.4684188360535,\"latitude\":37.89055748048133},{\"longitude\":112.46943457603459,\"latitude\":37.88994495705907},{\"longitude\":112.4705004607514,\"latitude\":37.889126321859806},{\"longitude\":112.47154838724066,\"latitude\":37.88888634568157},{\"longitude\":112.47216716028754,\"latitude\":37.88830767755853},{\"longitude\":112.47043654873971,\"latitude\":37.88775419198699},{\"longitude\":112.46926475867633,\"latitude\":37.88581117133405},{\"longitude\":112.46986766897146,\"latitude\":37.88324233889381},{\"longitude\":112.47128597080712,\"latitude\":37.88019919744106},{\"longitude\":112.47246823824946,\"latitude\":37.87754548393391},{\"longitude\":112.47304969087247,\"latitude\":37.874756173235895},{\"longitude\":112.47303032867615,\"latitude\":37.87181431156946},{\"longitude\":112.472324320972,\"latitude\":37.86880457613251},{\"longitude\":112.4713179058582,\"latitude\":37.866133512444115},{\"longitude\":112.47035440608863,\"latitude\":37.86322518721438},{\"longitude\":112.47021090790633,\"latitude\":37.862029823065946},{\"longitude\":112.47022834226493,\"latitude\":37.86080902815422},{\"longitude\":112.47019883796577,\"latitude\":37.85860455801826},{\"longitude\":112.47022566005592,\"latitude\":37.85451732560189},{\"longitude\":112.47214612171052,\"latitude\":37.84830768707056},{\"longitude\":112.47472104236488,\"latitude\":37.84522385231535},{\"longitude\":112.47815426990394,\"latitude\":37.84261435294711},{\"longitude\":112.48308953449134,\"latitude\":37.83912357955987},{\"longitude\":112.4847203175724,\"latitude\":37.83717477498551},{\"longitude\":112.48557862445716,\"latitude\":37.83522591892508},{\"longitude\":112.48610366687183,\"latitude\":37.832513057351505},{\"longitude\":112.48555582568054,\"latitude\":37.82980009602177},{\"longitude\":112.48248603746299,\"latitude\":37.82484845272313},{\"longitude\":112.47995269104842,\"latitude\":37.81998122871294},{\"longitude\":112.47870747551326,\"latitude\":37.81754749634335},{\"longitude\":112.47812778308992,\"latitude\":37.816245844329636},{\"longitude\":112.4779128710926,\"latitude\":37.814876363507864},{\"longitude\":112.47779317751531,\"latitude\":37.813494143440586},{\"longitude\":112.47788806065921,\"latitude\":37.81211189749517},{\"longitude\":112.47827094599609,\"latitude\":37.809482949777674},{\"longitude\":112.47814957603816,\"latitude\":37.8070573479433},{\"longitude\":112.47768119528894,\"latitude\":37.80557325992082},{\"longitude\":112.47723427221183,\"latitude\":37.80407218812318},{\"longitude\":112.47677930250768,\"latitude\":37.802908051700236},{\"longitude\":112.47658182486896,\"latitude\":37.801608260908445},{\"longitude\":112.47662038162355,\"latitude\":37.80024062794617},{\"longitude\":112.47687351509933,\"latitude\":37.79887296966116},{\"longitude\":112.47771237596874,\"latitude\":37.79639191329197},{\"longitude\":112.47806743338708,\"latitude\":37.79504961778886},{\"longitude\":112.47820791408424,\"latitude\":37.79377512319197},{\"longitude\":112.47816030487422,\"latitude\":37.7916075973264},{\"longitude\":112.47725438877944,\"latitude\":37.789440007873644},{\"longitude\":112.47568328484897,\"latitude\":37.787323228242165},{\"longitude\":112.47466169849042,\"latitude\":37.786383521868885},{\"longitude\":112.47359719678764,\"latitude\":37.785579469244574},{\"longitude\":112.47069571718578,\"latitude\":37.783360826865334},{\"longitude\":112.46818047568206,\"latitude\":37.780972525541},{\"longitude\":112.46731982186441,\"latitude\":37.77963418981305},{\"longitude\":112.46654499873523,\"latitude\":37.77795663171762},{\"longitude\":112.46628515973691,\"latitude\":37.776516479799376},{\"longitude\":112.46628281280402,\"latitude\":37.77487277226742},{\"longitude\":112.4666532929242,\"latitude\":37.773127261938946},{\"longitude\":112.46753875717525,\"latitude\":37.77144955618634},{\"longitude\":112.46913802430038,\"latitude\":37.76809403047927},{\"longitude\":112.47043688401584,\"latitude\":37.765213314361695},{\"longitude\":112.47147825166587,\"latitude\":37.76219677767211},{\"longitude\":112.47226212725047,\"latitude\":37.758823867792714},{\"longitude\":112.47271843805913,\"latitude\":37.757205214032304},{\"longitude\":112.47338932558898,\"latitude\":37.7555865248553},{\"longitude\":112.47440353587274,\"latitude\":37.75366242157076},{\"longitude\":112.47606147632007,\"latitude\":37.7521454501451},{\"longitude\":112.4783846046031,\"latitude\":37.750501200907614},{\"longitude\":112.48109397098426,\"latitude\":37.74923017982551},{\"longitude\":112.4904843847454,\"latitude\":37.747801795948746},{\"longitude\":112.50021812126045,\"latitude\":37.7467805959197},{\"longitude\":112.50956561967735,\"latitude\":37.74486010785324},{\"longitude\":112.51453977629546,\"latitude\":37.74342474643445},{\"longitude\":112.51968559429054,\"latitude\":37.742328720284505},{\"longitude\":112.5241447667778,\"latitude\":37.742348339666385},{\"longitude\":112.52868976995353,\"latitude\":37.74318242371454},{\"longitude\":112.53297728106384,\"latitude\":37.74510242508033},{\"longitude\":112.53709313079719,\"latitude\":37.7473617166671},{\"longitude\":112.54534628793601,\"latitude\":37.75140504545784},{\"longitude\":112.54891496703033,\"latitude\":37.75274800101568},{\"longitude\":112.5515395085514,\"latitude\":37.755448153319286},{\"longitude\":112.55276930138473,\"latitude\":37.75584948307612},{\"longitude\":112.55571570798759,\"latitude\":37.75577579124952},{\"longitude\":112.56385487124328,\"latitude\":37.758314670943776},{\"longitude\":112.5694694052637,\"latitude\":37.759109080112985},{\"longitude\":112.57521268531684,\"latitude\":37.75946241318194},{\"longitude\":112.58042287632827,\"latitude\":37.75973649033415},{\"longitude\":112.5856330673397,\"latitude\":37.76007842268836},{\"longitude\":112.58954708084468,\"latitude\":37.76060563156978},{\"longitude\":112.59337526366119,\"latitude\":37.76147211201093},{\"longitude\":112.59733219251041,\"latitude\":37.762542144360765},{\"longitude\":112.60094579860572,\"latitude\":37.76391749885141},{\"longitude\":112.60736700698737,\"latitude\":37.765997577846164},{\"longitude\":112.61002909943466,\"latitude\":37.7666320857112},{\"longitude\":112.61273410722617,\"latitude\":37.76713088882409},{\"longitude\":112.61669070079927,\"latitude\":37.76754818990631},{\"longitude\":112.62077604040508,\"latitude\":37.767677130031245},{\"longitude\":112.62511887207631,\"latitude\":37.767602522360406},{\"longitude\":112.62941878840331,\"latitude\":37.76766361319409},{\"longitude\":112.63399765446786,\"latitude\":37.767504193840935},{\"longitude\":112.63874818190936,\"latitude\":37.768158963904},{\"longitude\":112.64349870935087,\"latitude\":37.7698314437524},{\"longitude\":112.64600271910433,\"latitude\":37.77065070787182},{\"longitude\":112.64712597794835,\"latitude\":37.771238432354785},{\"longitude\":112.64769133731727,\"latitude\":37.77206361071429},{\"longitude\":112.6474858130515,\"latitude\":37.77366471337249},{\"longitude\":112.64702279672031,\"latitude\":37.77513009672995},{\"longitude\":112.64669757887725,\"latitude\":37.778264295092725},{\"longitude\":112.64759209558372,\"latitude\":37.78432879165564},{\"longitude\":112.64681626662616,\"latitude\":37.78740324780157},{\"longitude\":112.64501046940688,\"latitude\":37.78986711890549},{\"longitude\":112.64071222946052,\"latitude\":37.79445549223585},{\"longitude\":112.63271656438712,\"latitude\":37.80259717953443},{\"longitude\":112.628461239785,\"latitude\":37.80653205968931},{\"longitude\":112.62360510036353,\"latitude\":37.8101615848414},{\"longitude\":112.62016516730193,\"latitude\":37.81352394150184},{\"longitude\":112.6172402183712,\"latitude\":37.817699789056356},{\"longitude\":112.61547398373489,\"latitude\":37.82299408299194},{\"longitude\":112.61490937873725,\"latitude\":37.82842358418254},{\"longitude\":112.61476722165946,\"latitude\":37.84087428237429},{\"longitude\":112.61496838733558,\"latitude\":37.85332287914327}]}")
	//pointsByte := []byte("{\"points\":[{\"longitude\":112.50597180426121,\"latitude\":37.91169840227521},{\"longitude\":112.48468579351902,\"latitude\":37.76854143068962},{\"longitude\":112.65222729742527,\"latitude\":37.770712548328554},{\"longitude\":112.64192761480808,\"latitude\":37.91224013519099}]}")
	//pointsByte := []byte("{\"points\":[{\"longitude\":112.50874480037078,\"latitude\":37.85782490806801},{\"longitude\":112.50874480037078,\"latitude\":37.85777096797502},{\"longitude\":112.50908710037083,\"latitude\":37.85777096797502},{\"longitude\":112.50908710037083,\"latitude\":37.85782490806801}]}")
	type pointsType struct {
		Points []ParamsAreaPoint `json:"points"`
	}
	var points1 pointsType
	err := json.Unmarshal(pointsByte, &points1)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	b = CheckXYInArea(&ArgsCheckXYInArea{
		Point: MapMathArgs.ParamsPoint{
			PointType: "GCJ-02",
			Longitude: 112.53961743414402,
			Latitude:  37.86997298978658,
		}, Area: ParamsArea{
			ID:        1,
			PointType: "GCJ-02",
			Points:    points1.Points,
		},
	})
	if !b {
		t.Error(b)
	}
	b = CheckXYInArea(&ArgsCheckXYInArea{
		Point: MapMathArgs.ParamsPoint{
			PointType: "GCJ-02",
			Longitude: 112.51562,
			Latitude:  37.85929,
		}, Area: ParamsArea{
			ID:        1,
			PointType: "GCJ-02",
			Points:    points1.Points,
		},
	})
	if !b {
		t.Error(b)
	}
}
