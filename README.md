# frpCracker
一款golang编写的，批量检测frp server未授权访问、弱token的工具

![image-20230501143507003](/Users/coco/Desktop/Code/frp/frpCracker/assets/image-20230501143507003.png)



# 使用说明

将要检测的目标放入target.txt，如127.0.0.1:7000

检测的弱口令放入token.txt，一行一个

```shell
./frpCracker -l target.txt -tl token.txt
```

若不指定 -tl参数，则为检测未授权访问

```shell
./frpCracker -l target.txt
```



完整参数

```
Usage of ./frpCracker:
  -l string
    	输入文件，支持IP:Port格式，一行一个
  -o string
    	输出文件 (default "result.txt")
  -t int
    	线程数量 (default 20)
  -tl string
    	Token的输入文件，一行一个
  -v string
    	指定爆破时客户端的版本 (default "0.48.0")
```



# 免责声明

本工具仅面向**合法授权**的企业安全建设行为，如您需要测试本工具的可用性，请自行搭建靶机环境。

在使用本工具进行检测时，您应确保该行为符合当地的法律法规，并且已经取得了足够的授权。**请勿对非授权目标进行扫描。**

如您在使用本工具的过程中存在任何非法行为，您需自行承担相应后果，我们将不承担任何法律及连带责任。

在安装并使用本工具前，请您**务必审慎阅读、充分理解各条款内容**，限制、免责条款或者其他涉及您重大权益的条款可能会以加粗、加下划线等形式提示您重点注意。 除非您已充分阅读、完全理解并接受本协议所有条款，否则，请您不要安装并使用本工具。您的使用行为或者您以其他任何明示或者默示方式表示接受本协议的，即视为您已阅读并同意本协议的约束。
