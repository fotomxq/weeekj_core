package ToolsInstall

func InstallIOT() (err error) {
	/** 等待完成该模块改造后重新放开
	//上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	//写入action.json配置数据
	type dataInstallActionData struct {
		//标识码
		Mark string `json:"mark"`
		//名称
		Name string `json:"name"`
		//描述
		Des string `json:"des"`
		//默认过期时间
		ExpireTime string `json:"expireTime"`
		//过期时间
		//配置信息
		ConfigData []IOTCore.FieldsActionConfigType `json:"configData"`
		//连接方式
		// mqtt-client MQTT对设备推送
		// mqtt-group MQTT对设备组推送
		// none 设备驱动自主发现、业务代码其他渠道推送
		ConnectType string `json:"connectType"`
	}
	type dataInstallAction struct {
		Data []dataInstallActionData
	}
	var data dataInstallAction
	err = loadConfigFile(fmt.Sprint(".", CoreFile.Sep, "device", CoreFile.Sep, "action.json"), &data)
	if err != nil {
		return
	}
	for _, v := range data.Data {
		_, err = IOTCore.GetAction(ctx, &IOTCore.ArgsGetAction{
			Mark: v.Mark,
		})
		if err == nil {
			err = IOTCore.UpdateAction(ctx, &IOTCore.ArgsActionUpdate{
				Mark:        v.Mark,
				Name:        v.Name,
				Des:         v.Des,
				ExpireTime:  v.ExpireTime,
				ConfigData:  v.ConfigData,
				ConnectType: v.ConnectType,
			})
			if err != nil {
				return
			}
			continue
		}
		_, err = IOTCore.CreateAction(ctx, &IOTCore.ArgsActionCreate{
			Mark:        v.Mark,
			Name:        v.Name,
			Des:         v.Des,
			ExpireTime:  v.ExpireTime,
			ConfigData:  v.ConfigData,
			ConnectType: v.ConnectType,
		})
		if err != nil {
			return
		}
		CoreLog.Info("create new action, ", v.Mark, " ", v.Name)
	}
	//写入group数据
	type dataInstallGroupData struct {
		//标识码
		Mark string `json:"mark"`
		//名称
		Name string `json:"name"`
		//支持动作的标识码
		Action []string `json:"action"`
		//默认设备掉线过期时间
		ExpireTime string `json:"expireTime"`
		//使用类型
		// public 公共设备 / private 私有设备
		UseType string `json:"useType"`
		//扩展信息
		Infos []IOTCore.FieldsGroupInfoType `json:"infos"`
	}
	type dataInstallGroup struct {
		Data []dataInstallGroupData
	}
	var dataGroup dataInstallGroup
	err = loadConfigFile(fmt.Sprint(".", CoreFile.Sep, "device", CoreFile.Sep, "group.json"), &dataGroup)
	if err != nil {
		return
	}
	for _, v := range dataGroup.Data {
		var vGroupData IOTCore.FieldsGroupType
		vGroupData, err = IOTCore.GetGroup(ctx, &IOTCore.ArgsGetGroup{
			Mark: v.Mark,
		})
		if err == nil {
			err = IOTCore.UpdateGroup(ctx, &IOTCore.ArgsGroupUpdate{
				ID:         vGroupData.ID.Hex(),
				Mark:       v.Mark,
				Name:       v.Name,
				Action:     v.Action,
				ExpireTime: v.ExpireTime,
				UseType:    v.UseType,
				Infos:      v.Infos,
			})
			if err != nil {
				return
			}
			continue
		}
		_, err = IOTCore.CreateGroup(ctx, &IOTCore.ArgsGroupCreate{
			Mark:       v.Mark,
			Name:       v.Name,
			Action:     v.Action,
			ExpireTime: v.ExpireTime,
			UseType:    v.UseType,
			Infos:      v.Infos,
		})
		if err != nil {
			return
		}
		CoreLog.Info("create new group, ", v.Mark, " ", v.Name)
	}
	 */
	//反馈成功
	return
}
