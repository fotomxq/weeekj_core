package BaseBPM

import (
	"time"
)

// FieldsSlot 插槽定义
type FieldsSlot struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//所属主题分类
	ThemeCategoryID int64 `db:"theme_category_id" json:"themeCategoryId" check:"id"`
	//所属主题
	// 插槽可用于的主题域
	ThemeID int64 `db:"theme_id" json:"themeId" check:"id"`
	//值类型
	// 插槽的值类型
	// input 输入框; text 文本域; radio 单选项; checkbox 多选项; select 下拉单选框;
	// date: 日期; time: 时间; datetime: 日期时间;
	// file: 文件ID; files 文件ID列; image: 图片; images 图片列; audio: 音频; video: 视频; videos 视频列; url: URL;
	// email: 邮箱; phone: 电话; id: ID; password: 密码;
	// code: 代码; html: HTML; markdown: Markdown; xml: XML; yaml: YAML;
	ValueType string `db:"value_type" json:"valueType" check:"des" min:"1" max:"3000"`
	//默认值
	// 插槽的默认值
	DefaultValue string `db:"default_value" json:"defaultValue" check:"des" min:"1" max:"3000"`
	//参数
	// 根据组件需求，自定义参数内容
	Params string `db:"params" json:"params"`
}
