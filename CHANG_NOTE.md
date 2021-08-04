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