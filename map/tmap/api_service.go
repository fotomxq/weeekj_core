package MapTMap

import (
	"encoding/xml"
	"fmt"
	CoreSQL2 "gitee.com/weeekj/weeekj_core/v5/core/sql2"
)

// GetGeocoder 逆地理编码查询
// url: http://api.tianditu.gov.cn/geocoder?postStr={'lon':116.37304,'lat':39.92594,'ver':1}&type=geocode&tk=您的密钥
// 参考文档：http://lbs.tianditu.gov.cn/server/geocoding.html
func GetGeocoder(lon float64, lat float64) (data []byte, err error) {
	url := fmt.Sprint("http://api.tianditu.gov.cn/geocoder?postStr={'lon':", lon, ",'lat':", lat, ",'ver':1}&type=geocode")
	data, err = postURL(url)
	return
}

// GetGeocoderInterface 地理编码接口
// url: http://api.tianditu.gov.cn/geocoder?ds={"keyWord":"延庆区北京市延庆区延庆镇莲花池村前街50夕阳红养老院"}&tk=您的密钥
// 参考文档：http://lbs.tianditu.gov.cn/server/geocodinginterface.html
func GetGeocoderInterface(keyWord string) (data []byte, err error) {
	url := fmt.Sprint("http://api.tianditu.gov.cn/geocoder?ds={\"keyWord\":\"", keyWord, "\"}")
	data, err = postURL(url)
	return
}

// GetSearchV2Normal 地名搜索V2.0普通搜索示例
// url: http://api.tianditu.gov.cn/v2/search?postStr={"keyWord":"北京大学","level":12,"mapBound":"116.02524,39.83833,116.65592,39.99185","queryType":1,"start":0,"count":10}&type=query&tk=您的密钥
// 参考文档：http://lbs.tianditu.gov.cn/server/search2.html
func GetSearchV2Normal(keyWord string, level int, mapBound string, queryType int, start int, count int) (data []byte, err error) {
	url := fmt.Sprint("http://api.tianditu.gov.cn/v2/search?postStr={\"keyWord\":\"", keyWord, "\",\"level\":", level, ",\"mapBound\":\"", mapBound, "\",\"queryType\":", queryType, ",\"start\":", start, ",\"count\":", count, "}&type=query")
	data, err = postURL(url)
	return
}

// GetSearchV2View 视野内搜索示例
// url: http://api.tianditu.gov.cn/v2/search?postStr={"keyWord":"医院","level":12,"mapBound":"116.02524,39.83833,116.65592,39.99185","start":0,"count":10}&type=query&tk=您的密钥
// 参考文档：http://lbs.tianditu.gov.cn/server/search2.html
func GetSearchV2View(keyWord string, level int, mapBound string, start int, count int) (data []byte, err error) {
	url := fmt.Sprint("http://api.tianditu.gov.cn/v2/search?postStr={\"keyWord\":\"", keyWord, "\",\"level\":", level, ",\"mapBound\":\"", mapBound, "\",\"queryType\":2,\"start\":", start, ",\"count\":", count, "}&type=query")
	data, err = postURL(url)
	return
}

// GetSearchV2Around 周边搜索示例
// url: http://api.tianditu.gov.cn/v2/search?postStr={"keyWord":"公园","level":12,"queryRadius":5000,"pointLonlat":"116.48016,39.93136","queryType":3,"start":0,"count":10}&type=query&tk=您的密钥
// 参考文档：http://lbs.tianditu.gov.cn/server/search2.html
func GetSearchV2Around(keyWord string, level int, queryRadius int, pointLonlat string, start int, count int) (data []byte, err error) {
	url := fmt.Sprint("http://api.tianditu.gov.cn/v2/search?postStr={\"keyWord\":\"", keyWord, "\",\"level\":", level, ",\"queryRadius\":", queryRadius, ",\"pointLonlat\":\"", pointLonlat, "\",\"queryType\":3,\"start\":", start, ",\"count\":", count, "}&type=query")
	data, err = postURL(url)
	return
}

// GetSearchV2Polygon 多边形搜索
// url: http://api.tianditu.gov.cn/v2/search?postStr={"keyWord":"学校","polygon":"118.93232636500011,27.423305726000024,118.93146426300007,27.30976105800005,118.80356153600007,27.311829507000027,118.80469010700006,27.311829508000073,118.8046900920001,27.32381604300008,118.77984777400002,27.32381601800006,118.77984779100007,27.312213007000025,118.76792266100006,27.31240586100006,118.76680145600005,27.429347074000077,118.93232636500011,27.423305726000024","queryType":10,"start":0,"count":10}&type=query&tk=您的密钥
// 参考文档：http://lbs.tianditu.gov.cn/server/search2.html
func GetSearchV2Polygon(keyWord string, polygon string, start int, count int) (data []byte, err error) {
	url := fmt.Sprint("http://api.tianditu.gov.cn/v2/search?postStr={\"keyWord\":\"", keyWord, "\",\"polygon\":\"", polygon, "\",\"queryType\":10,\"start\":", start, ",\"count\":", count, "}&type=query")
	data, err = postURL(url)
	return
}

// GetSearchV2Specify 行政区划区域搜索
// url: http://api.tianditu.gov.cn/v2/search?postStr={"keyWord":"商厦","queryType":12,"start":0,"count":10,"specify":"156110108"}&type=query&tk=您的密钥
// 参考文档：http://lbs.tianditu.gov.cn/server/search2.html
func GetSearchV2Specify(keyWord string, start int, count int, specify string) (data []byte, err error) {
	url := fmt.Sprint("http://api.tianditu.gov.cn/v2/search?postStr={\"keyWord\":\"", keyWord, "\",\"queryType\":12,\"start\":", start, ",\"count\":", count, ",\"specify\":\"", specify, "\"}&type=query")
	data, err = postURL(url)
	return
}

// GetSearchV2DataTypes 数据分类搜索
// url: http://api.tianditu.gov.cn/v2/search?postStr={"queryType":13,"start":0,"count":5,"specify":"156110000","dataTypes":"法院,公园"}&type=query&tk=您的密钥
// 参考文档：http://lbs.tianditu.gov.cn/server/search2.html
func GetSearchV2DataTypes(start int, count int, specify string, dataTypes string) (data []byte, err error) {
	url := fmt.Sprint("http://api.tianditu.gov.cn/v2/search?postStr={\"queryType\":13,\"start\":", start, ",\"count\":", count, ",\"specify\":\"", specify, "\",\"dataTypes\":\"", dataTypes, "\"}&type=query")
	data, err = postURL(url)
	return
}

// GetSearchV2Statistics 统计搜索
// url: http://api.tianditu.gov.cn/v2/search?postStr={"keyWord":"学校","queryType":14,"specify":"156110108"}&type=query&tk=您的密钥
// 参考文档：http://lbs.tianditu.gov.cn/server/search2.html
func GetSearchV2Statistics(keyWord string, specify string) (data []byte, err error) {
	url := fmt.Sprint("http://api.tianditu.gov.cn/v2/search?postStr={\"keyWord\":\"", keyWord, "\",\"queryType\":14,\"specify\":\"", specify, "\"}&type=query")
	data, err = postURL(url)
	return
}

// DataGetDrive 驾车规划数据
type DataGetDrive struct {
	XMLName xml.Name `xml:"result"`
	Orig    string   `xml:"orig,attr" json:"orig"`
	Mid     string   `xml:"mid,attr" json:"mid"`
	Dest    string   `xml:"dest,attr" json:"dest"`
	Params  struct {
		Orig    string `xml:"orig" json:"orig"`
		Dest    string `xml:"dest" json:"dest"`
		Mid     string `xml:"mid" json:"mid"`
		Key     string `xml:"key" json:"key"`
		Width   string `xml:"width" json:"width"`
		Height  string `xml:"height" json:"height"`
		Style   string `xml:"style" json:"style"`
		Version string `xml:"version" json:"version"`
		Sort    string `xml:"sort" json:"sort"`
	} `xml:"parameters" json:"parameters"`
	Routes struct {
		Count string `xml:"count,attr" json:"count"`
		Time  string `xml:"time,attr" json:"time"`
		Items []struct {
			ID         string `xml:"id,attr" json:"id"`
			StrGuide   string `xml:"strguide" json:"strguide"`
			Signage    string `xml:"signage" json:"signage"`
			StreetName string `xml:"streetName" json:"streetName"`
			NextStreet string `xml:"nextStreetName" json:"nextStreetName"`
			TollStatus string `xml:"tollStatus" json:"tollStatus"`
			TurnLatLon string `xml:"turnlatlon" json:"turnLatLon"`
		} `xml:"item" json:"items"`
	} `xml:"routes" json:"routes"`
	Simple struct {
		Items []struct {
			ID             string  `xml:"id,attr" json:"id"`
			StrGuide       string  `xml:"strguide" json:"strguide"`
			StreetNames    string  `xml:"streetNames" json:"streetNames"`
			LastStreetName string  `xml:"lastStreetName" json:"lastStreetName"`
			LinkStreetName string  `xml:"linkStreetName" json:"linkStreetName"`
			Signage        string  `xml:"signage" json:"signage"`
			TollStatus     string  `xml:"tollStatus" json:"tollStatus"`
			TurnLatLon     string  `xml:"turnlatlon" json:"turnLatLon"`
			StreetLatLon   string  `xml:"streetLatLon" json:"streetLatLon"`
			StreetDistance float64 `xml:"streetDistance" json:"streetDistance"`
			SegmentNumber  string  `xml:"segmentNumber" json:"segmentNumber"`
		} `xml:"item" json:"items"`
	} `xml:"simple" json:"simple"`
	Distance    float64 `xml:"distance" json:"distance"`
	Duration    float64 `xml:"duration" json:"duration"`
	RouteLatLon string  `xml:"routelatlon" json:"routeLatLon"`
	MapInfo     struct {
		Center string `xml:"center" json:"center"`
		Scale  string `xml:"scale" json:"scale"`
	} `xml:"mapinfo" json:"mapInfo"`
}

// GetDrive 驾车规划
// url:http://api.tianditu.gov.cn/drive?postStr={"orig":"116.35506,39.92277","dest":"116.39751,39.90854","style":"0"}&type=search&tk=您的密钥
// 参考文档：http://lbs.tianditu.gov.cn/server/drive.html
// style: 默认0 （0：最快路线，1：最短路线，2：避开高速，3：步行）
func GetDrive(orig CoreSQL2.ArgsGPS, dest CoreSQL2.ArgsGPS, style string) (data DataGetDrive, err error) {
	url := fmt.Sprint("http://api.tianditu.gov.cn/drive?postStr={\"orig\":\"", fmt.Sprint(orig.Longitude, ",", orig.Latitude), "\",\"dest\":\"", fmt.Sprint(dest.Longitude, ",", dest.Latitude), "\",\"style\":\"", style, "\"}&type=search")
	var dataByte []byte
	dataByte, err = postURL(url)
	if err != nil {
		return
	}
	err = xml.Unmarshal(dataByte, &data)
	if err != nil {
		return
	}
	return
}
