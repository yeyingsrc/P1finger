

<img src="./img/image-20240811182803001.png" alt="image-20240811182803001" style="zoom: 67%;" />

<h3 align="center">P1finger 一款红队行动下的重点资产指纹识别工具</h3>



# 功能特色

* P1finger `v0.0.4` 版本开始支持两种指纹识别模式

  1. 基于本地规则库的重点资产指纹识别

  1. 基于Fofa测绘系统的web指纹识别


* 输出到表格或其他格式
* 支持 http / socks代理使用



## Version

当前最新 `beta_version v0.0.5` （2025/2/19更新）[更新日志参见](https://github.com/P001water/P1finger/blob/master/更新日志.md)

---

# 基本使用

## 配置Fofa key

P1finger在命令行下首次运行生成 `p1fingerConf.yaml` 配置文件，在配置文件中填上 `email` 和 `key` 即可。

文件内容参考

```
FofaCredentials:
    Email: P001water@163.com
    ApiKey: xxxx
```



## 开始使用

`-m` 参数切换模式，

1. `-m rule` 基于本地规则库模式，（默认模式）
2. `-m fofa` 基于fofa的采集模式，（手动开启）

基于本地规则库模式使用

1.  单个目标探测

```
P1finger -u [target]
P1finger -uf [target file] //-uf 指定url文件
```

基于fofa的采集模式

```
P1finger -m fofa -u [target]
P1finger -m fofa -uf [target file]
```

![image-20250212212234854](./img/image-20250212212234854.png)



3. 输出到 Excel 表格文件或者json格式保存到文件

```
P1finger [-m fofa] -uf [target file] -o xlsx
```

![image-20250219152847788](./img/image-20250219152847788.png)

4. socks5 代理

```
P1finger.exe -uf urls.txt -socks 127.0.0.1:4781
```

5. http 代理

```
P1finger.exe -uf urls.txt -httpproxy 127.0.0.1:4781
```

