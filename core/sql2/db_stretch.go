package CoreSQL2

/**
数据表自动扩容伸缩设计
1. 自动根据ID、mark、以及其他业务关键信息，将表分拆处理
2. 具备一定规则性的处理机制，外部模块不需要考虑具体分表的设计，只需要关注业务本身和提供业务关键信息后获取数据即可
3. 搭配client模块设置完成该设计，即启动自动伸缩后，将自动根据设置进行拆分的内部处理，不过也对外暴露获取表名方法，方便特殊的行为调用处理
TODO: 根据整体设计思路，完成数据表的扩容模块设计
*/

//根据关键信息获取当前模块的表
