# LazyDog



通过调用hunter的API接口获取指定域名或IP的资产，再~~抄~~借鉴[Ehole](https://github.com/EdgeSecurityTeam/EHole)进行指纹识别。

- 增加了导出csv格式
- 将[web_fingerprint_v3.json](https://github.com/0x727/FingerprintHub/blob/main/web_fingerprint_v3.json)指纹和原有指纹合并在一起
- 优化了导出功能，将hunter资产导出，并将资产全放在Result里

# 命令使用

```
Usage:
  lazydog hunter [flags]

Flags:
  -d, --domain string   从Hunter提取域名资产进行指纹识别
  -h, --help            help for hunter
  -i, --ip string       从Hunter提取IP资产进行指纹识别，支持ip或者ip段，例如：192.168.1.1 | 192.168.1.0/24
  -l, --local string    从本地文件读取资产，进行指纹识别，支持无协议，例如：192.168.1.1:9090 | http://192.168.1.1:9090
  -o, --output string   输出所有结果，支持csv、json和xlsx后缀的文件。
  -p, --proxy string    指定访问目标时的代理，支持http代理和socks5，例如：http://127.0.0.1:8080 | socks5://127.0.0.1:8080
  -t, --thread int      指纹识别线程大小。 (default 100)
  -u, --url string      识别单个目标。
```

```
./lazydog hunter -d domain.com -o xxx.csv
```
huanter扫描doamin.com(搜索参数是domain.suffix="domain.com")，再进行指纹识别。
```
./lazydog hunter -i xx.xx.xx.xx -o xxx.csv
```
hunter搜索参数是ip="xx.xx.xx.xx"
