# AutoWeChat
基于GO开发的微信定时自动发送工具，通过爬虫拉取数据，定时向多人发送天气和每日一言消息

引用库：
>>>[itchat4go](https://github.com/newflydd/itchat4go)-微信个人号接口  
>>>[cron](https://github.com/robfig/cron)-定时任务  
>>>[simplejson](https://github.com/simplejson/simplejson)-json解析  
>>>[yaml](https://github.com/go-yaml/yaml)-yaml配置  

数据来源:  
天气来自[SOJSON](https://www.sojson.com/blog/305.html)  
[每日一言](https://api.ooopn.com/yan/api.php)

效果:
 ![image](https://github.com/hzy1754785321/AutoWeChat/blob/master/1.png)  
 一般情况下是显示地区，我是主动隐藏了地区  
 
 运行：
 安装好依赖库之后，自行修改配置文件，然后直接运行便可
