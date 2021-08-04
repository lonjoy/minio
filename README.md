### MinIO源码
- https://github.com/minio/minio

### 修改功能
- 预签名后，前端上传文件后增加返回值
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