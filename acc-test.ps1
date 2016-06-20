$go = Get-Command go

$oldTF_ACC = $env:TF_ACC
$env:TF_ACC="1"

try {
	& $go test -v -timeout 120m
}
finally {
	$env:TF_ACC = $oldTF_ACC
}
