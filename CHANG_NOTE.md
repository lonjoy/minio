### MinIO源码修改详细操作步骤
- 参考链接
    - https://gaohaoyang.github.io/2015/04/12/Syncing-a-fork/
    - https://www.zhihu.com/question/28676261

- Fork
```
fork https://github.com/minio/minio
git clone git@github.com:lonjoy/minio.git
```

- 设置upstream
```
git remote add upstream https://github.com/minio/minio
git remote -v
```

- 同步更新源码版本
```
git fetch upstream
git checkout master
git merge upstream/master
```

- Push到Fork后的仓库
```
git push origin master
```

### 本地多仓库同步Push详细操作步骤
- 添加多个仓库
```
git remote add gitlab ssh://git@git.ynedut.cn:2222/wjs/wjs-minio.git
```

- 本地仓库推送到多个远端仓库
```
git push -u origin
git push -u gitlab
```

- OR - 同时push到多个源
```
git remote set-url --add --push origin ssh://git@git.ynedut.cn:2222/wjs/wjs-minio.git
git push origin master
```


### 多仓库分支合并

- pull upstream/master分支
- 合并upstream/mashter到origin/master、origin/develop
- push origin/master、origin/develop
- 合并origin/develop到gitlab/develop
- push gitlab


### MinIO源码仓库地址
- https://github.com/minio/minio

### 修改功能
- 前端上传文件后增加返回值
```
// 源码 cmd\object-handlers.go:1762
// writeSuccessResponseHeadersOnly(w)
// Marshal API response
resJsonBytes, err := json.Marshal(objInfo)
if err != nil {
	writeErrorResponseJSON(ctx, w, toAPIError(ctx, err), r.URL)
	return
}
writeSuccessResponseJSON(w, resJsonBytes)

```
- 获取自定义请求参数
```
// 源码 cmd\handler-utils.go:249
// http://{host}/upload/7ca94336-1d08-4439-b308-0ad26599f70b.png?X-Wjs-OriginalFilename=Screenshot_2021_0305_080343.png&response-content-type=application/json&originalFilename=Screenshot_2021_0305_080343.png&X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=admin/20210810/us-east-1/s3/aws4_request&X-Amz-Date=20210810T140139Z&X-Amz-Expires=7200&X-Amz-SignedHeaders=host&X-Amz-Signature=774fa07b412017dbf56ff356546f7e1dfccada0513ca5a4e5ddaaa270b6889f3
// r.URL.RawQuery = X-Wjs-OriginalFilename=default_bg.png&X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=admin/20210810/us-east-1/s3/aws4_request&X-Amz-Date=20210810T142138Z&X-Amz-Expires=7200&X-Amz-SignedHeaders=host&X-Amz-Signature=ed29139cd882ec952b13382e724f6633629acf1409ecfd98986bd57ae22ebd2d
// 获取自定义请求参数
rawQuery := strings.Split(r.URL.RawQuery, "&")
var temp []string
for _, v := range rawQuery{
    temp = strings.Split(v, "=")
    if strings.HasPrefix(temp[0], "X-Wjs-"){
        m[temp[0]] = temp[1]
    }
}
```
- 前端上传时将事件回调由异步变为同步，增加同步发送事件方法
```
// 源码 cmd\object-handlers.go:1746
// 同步发送事件通知
// 申明带1个缓冲的通道，用于接收事件发送结束的通知
done := make(chan bool, 1)
// Sync Notify object created event.
fmt.Printf("api handler: %d\n", time.Now().Unix())
sendEventSync(eventArgs{
    EventName:    event.ObjectCreatedPut,
    BucketName:   bucket,
    Object:       objInfo,
    ReqParams:    extractReqParams(r),
    RespElements: extractRespElements(w),
    UserAgent:    r.UserAgent(),
    Host:         handlers.GetSourceIP(r),
}, done)

// 等待事件发送goroutine, 会阻塞请求
<-done

// 源码 cmd\notification.go:1780
// 同步发送事件
func sendEventSync(args eventArgs, done chan bool) {
	args.Object.Size, _ = args.Object.GetActualSize()

	// avoid generating a notification for REPLICA creation event.
	if _, ok := args.ReqParams[xhttp.MinIOSourceReplicationRequest]; ok {
		return
	}
	// remove sensitive encryption entries in metadata.
	crypto.RemoveSensitiveEntries(args.Object.UserDefined)
	crypto.RemoveInternalEntries(args.Object.UserDefined)

	// globalNotificationSys is not initialized in gateway mode.
	if globalNotificationSys == nil {
		return
	}

	if globalHTTPListen.NumSubscribers() > 0 {
		globalHTTPListen.Publish(args.ToEvent(false))
	}
	fmt.Printf("sendEventSync: %d\n", time.Now().Unix())
	globalNotificationSys.SendSync(args, done)
}


// 源码 cmd\notification.go:788
// 同步发送通知
// Synchronization Send - sends event data to all matching targets.
func (sys *NotificationSys) SendSync(args eventArgs, done chan bool) {
	sys.RLock()
	targetIDSet := sys.bucketRulesMap[args.BucketName].Match(args.EventName, args.Object.Name)
	sys.RUnlock()

	if len(targetIDSet) == 0 {
		return
	}
	fmt.Printf("SendSync: %d\n", time.Now().Unix())
	sys.targetList.SendSync(args.ToEvent(true), targetIDSet, sys.targetResCh, done)
}


// 源码 internal\event\targetlist.go:152
// 同步发送事件到指定目标
// done - 通过channel返回事件是否发送完成
// Synchronization Send - sends events to targets identified by target IDs.
func (list *TargetList) SendSync(event Event, targetIDset TargetIDSet, resCh chan<- TargetIDResult, done chan bool) {
	// 声明一个等待组
	var wg sync.WaitGroup
	for id := range targetIDset {
		fmt.Printf("SendSync: id=%s, %d\n", id, time.Now().Unix())
		list.RLock()
		target, ok := list.targets[id]
		list.RUnlock()
		if ok {
			// 每个任务开始时，将等待组加1
			wg.Add(1)
			go func(id TargetID, target Target) {
				// 使用defer, 使方法完成时将等待组减1
				defer wg.Done()
				tgtRes := TargetIDResult{ID: id}
				if err := target.Save(event); err != nil {
					tgtRes.Err = err
				}
				resCh <- tgtRes
			}(id, target)
		} else {
			resCh <- TargetIDResult{ID: id}
		}
	}
	// 等待所有任务完成
	wg.Wait()
	// 向channel中写入完成标志
	done <- true
}


```