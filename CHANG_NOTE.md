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
1. 前端上传文件后增加返回值
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
2. 获取自定义请求参数
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