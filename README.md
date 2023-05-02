# 哔哩哔哩（Bilibili）番剧名称扫描工具（Go）

## 简介
本工具为 [Bilibili-Bangumi-Scanner](https://github.com/Haoyi-Han/Bilibili-Bangumi-Scanner) 项目的 Golang 重构和优化加强版本，用于遍历 md 号爬取 B 站番剧名称，将爬取到的 md 号、番剧名称、番剧主页地址保存在 csv 文档中。

## 用法
用以下命令克隆本项目：
```shell
git clone https://github.com/Haoyi-Han/Bilibili-Bangumi-Scanner-Go.git
cd Bilibili-Bangumi-Scanner-Go
```

您需要安装 Go 的运行环境，之后请执行以下命令以构建和启动此项目：

```shell
go build
./Bilibili-Bangumi-Scanner-Go -h
```

程序将显示主要参数如下：

```shell
Usage of Bilibili-Bangumi-Scanner-Go:
  -begin int
        起始 media ID (default 28221450)
  -delimiter string
        分隔号 (default ";")
  -end int
        终止 media ID (default 28222450)
  -log
        是否启用日志打印 (default false)
  -output string
        输出文件名 (default "output.csv")
  -retry int
        最大重试次数 (default 5)
  -thread int
        线程数 (default 20)
```

您可以修改相应参数以适应自身需要。程序将以多线程运行，爬取的番剧名称信息将保存到同目录下的 csv 文件中。