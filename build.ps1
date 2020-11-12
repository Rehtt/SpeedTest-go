$version = '0.0.0';
$tmpContent = Get-Content './main.go';
for ($i = 0; $i -le $tmpContent.length; $i++) {
    if ($tmpContent[$i] -like '*speedtest *'){
        $version = [regex]::Matches($tmpContent[$i], '(([0-9]|([1-9]([0-9]*))).){2}([0-9]|([1-9]([0-9]*)))').Value;
    }
}

$windows = 'amd64', '386';
$linux = 'amd64', '386', 'arm', 'mipsle';

function build($os, $arch) {
    foreach ($a in $arch) {
        $env:GOOS = $os;
        $env:GOARCH = $a;
        $name = "speedtest_" + $version + "_" + $os + "_" + $a;
        if ($os -eq 'windows') {
            $name += '.exe';
        }
        "[" + $os + "_" + $a + "]"
        go build -ldflags '-w -s' -o $name;
        upx -9 $name;
    }
}

build 'windows' $windows;
build 'linux' $linux;