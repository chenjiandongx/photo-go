# 美女写真图爬虫 Golang 版

*其他版本*

* [photo-asyncio 版](https://github.com/chenjiandongx/photo-asyncio)
* [photo-gevent 版](https://github.com/chenjiandongx/photo-gevent)

### goroutine
最近发现 Golang 其实也是一门很有趣的语言，goroutine 并发模式很特别，所以决定拿个爬虫来试试。觉得还是拿妹子爬虫来试可能比较好一点，毕竟兴趣是最好的老师...

goroutine pool 使用了 [workerpool](https://github.com/gammazero/workerpool) 第三方库，该库介绍如下。

> Concurrency limiting goroutine pool. Limits the concurrency of task execution, not the number of tasks queued.


### 构建运行

**运行项目**

```bash
$ git clone https://github.com/chenjiandongx/photo-go.git
$ cd photo-go
$ go get
$ go run core.go
```

**图片数据**

图片地址数据保存在了 `data.txt`，共 17w+ 张照片，图片的数据是我从 [mmjpg](https://github.com/chenjiandongx/mmjpg) 和 [mzitu](https://github.com/chenjiandongx/mzitu) 里提取出来的。
```bash
$ wc -l data.txt
178075 data.txt
```

**运行效果**

![效果图](https://user-images.githubusercontent.com/19553554/41359931-28bc5e5a-6f5e-11e8-81ad-0ab5c4f6e26e.gif)


### License

MIT [©chenjiandongx](https://github.com/chenjiandongx)
