# IMChat

>
>即时通讯项目，支持好友、单聊、群聊
- api文档：https://www.apifox.cn/apidoc/shared-b851ba04-a8e5-481c-acda-858e9070396c

# 项目构成

```
api:用于定义接口函数 
cache:redis操作 
conf:初始化配置以及存放配置信息 
middlware：应用中间件 
model:数据库模型与初始化 
pkg--e:错误处理类 
   --util:工具类 
router：路由组 
serializer:序列化函数包 
service：服务模块 
```

# 演示流程

- 注册登录
- 单聊
- 群聊（即时）
- 添加好友
- 查看历史记录

# 待完善功能
- 群聊功能的完善，包括审核进群的申请信息，群成员管理（这部分实现了也就可以发送离线消息了）。
- 未读消息的批量处理
