Set-Location config
echo 'Contents of dev.yaml:'
etcdctl --endpoints=http://127.0.0.1:12379 put /webook $(Get-Content dev.yaml -Raw)
Set-Location ..
go mod tidy