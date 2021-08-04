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