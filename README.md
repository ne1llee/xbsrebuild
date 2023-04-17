# 香色闺阁xbs书源加解密工具

## linux、macOS在终端运行需要对程序添加可执行权限
```
chmod +x 程序路径
```

## 启动转换web服务
```
// 默认监听0.0.0.0:8282
xbsrebuild server 
// 指定监听地址
xbsrebuild server -s 127.0.0.1 -p 8282
```

## xbs 转 json
```
xbsrebuild xbs2json -i xx.xbs -o xx.json
```
## json 转 xbs
```
xbsrebuild json2xbs -i xx.json -o xx.xbs
```