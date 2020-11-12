// @Title  main.go
// @Description  请填写文件描述（需要改）
// @Author  Rehtt  2020/11/10 下午 4:59
// @Update  Rehtt  2020/11/10 下午 11:58

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

const version = "speedtest v0.1.0"

func main() {
	port := flag.String("p", "80", "端口")
	flag.Parse()
	http.HandleFunc("/", index)
	http.HandleFunc("/garbage", garbage)
	http.HandleFunc("/empty", empty)
	http.HandleFunc("/getIP", getIp)
	http.HandleFunc("/speedtest_worker.min.js", js)
	fmt.Println(":", *port)
	http.ListenAndServe(":"+*port, nil)
}

func js(writer http.ResponseWriter, request *http.Request) {
	const name = `function clearRequests(){if(xhr){for(var i=0;i<xhr.length;i++){if(useFetchAPI)try{xhr[i].cancelRequested=!0}catch(e){}try{xhr[i].onprogress=null,xhr[i].onload=null,xhr[i].onerror=null}catch(e){}try{xhr[i].upload.onprogress=null,xhr[i].upload.onload=null,xhr[i].upload.onerror=null}catch(e){}try{xhr[i].abort()}catch(e){}try{delete xhr[i]}catch(e){}}xhr=null}}function getIp(done){xhr=new XMLHttpRequest,xhr.onload=function(){clientIp=xhr.responseText,done()},xhr.onerror=function(){done()},xhr.open("GET",settings.url_getIp+"?r="+Math.random(),!0),xhr.send()}function dlTest(done){if(!dlCalled){dlCalled=!0;var totLoaded=0,startT=(new Date).getTime(),failed=!1;xhr=[];for(var testStream=function(i,delay){setTimeout(function(){if(1===testStatus)if(useFetchAPI)xhr[i]=fetch(settings.url_dl+"?r="+Math.random()+"&ckSize="+settings.garbagePhp_chunkSize).then(function(response){var reader=response.body.getReader(),consume=function(){return reader.read().then(function(result){return result.done?testStream(i):(totLoaded+=result.value.length,xhr[i].cancelRequested&&reader.cancel()),consume()}.bind(this))}.bind(this);return consume()}.bind(this));else{var prevLoaded=0,x=new XMLHttpRequest;xhr[i]=x,xhr[i].onprogress=function(event){if(1!==testStatus)try{x.abort()}catch(e){}var loadDiff=event.loaded<=0?0:event.loaded-prevLoaded;isNaN(loadDiff)||!isFinite(loadDiff)||0>loadDiff||(totLoaded+=loadDiff,prevLoaded=event.loaded)}.bind(this),xhr[i].onload=function(){try{xhr[i].abort()}catch(e){}testStream(i,0)}.bind(this),xhr[i].onerror=function(){failed=!0;try{xhr[i].abort()}catch(e){}delete xhr[i]}.bind(this);try{settings.xhr_dlUseBlob?xhr[i].responseType="blob":xhr[i].responseType="arraybuffer"}catch(e){}xhr[i].open("GET",settings.url_dl+"?r="+Math.random()+"&ckSize="+settings.garbagePhp_chunkSize,!0),xhr[i].send()}}.bind(this),1+delay)}.bind(this),i=0;i<settings.xhr_dlMultistream;i++)testStream(i,100*i);interval=setInterval(function(){var t=(new Date).getTime()-startT;if(!(200>t)){var speed=totLoaded/(t/1e3);dlStatus=(8*speed/925e3).toFixed(2),(t/1e3>settings.time_dl||failed)&&((failed||isNaN(dlStatus))&&(dlStatus="Fail"),clearRequests(),clearInterval(interval),done())}}.bind(this),200)}}function ulTest(done){if(!ulCalled){ulCalled=!0;var totLoaded=0,startT=(new Date).getTime(),failed=!1;xhr=[];for(var testStream=function(i,delay){setTimeout(function(){if(3===testStatus){var prevLoaded=0,x=new XMLHttpRequest;xhr[i]=x;var ie11workaround;try{xhr[i].upload.onprogress,ie11workaround=!1}catch(e){ie11workaround=!0}ie11workaround?(xhr[i].onload=function(){totLoaded+=262144,testStream(i,0)},xhr[i].onerror=function(){failed=!0;try{xhr[i].abort()}catch(e){}delete xhr[i]},xhr[i].open("POST",settings.url_ul+"?r="+Math.random(),!0),xhr[i].setRequestHeader("Content-Encoding","identity"),xhr[i].send(reqsmall)):(xhr[i].upload.onprogress=function(event){if(3!==testStatus)try{x.abort()}catch(e){}var loadDiff=event.loaded<=0?0:event.loaded-prevLoaded;isNaN(loadDiff)||!isFinite(loadDiff)||0>loadDiff||(totLoaded+=loadDiff,prevLoaded=event.loaded)}.bind(this),xhr[i].upload.onload=function(){testStream(i,0)}.bind(this),xhr[i].upload.onerror=function(){failed=!0;try{xhr[i].abort()}catch(e){}delete xhr[i]}.bind(this),xhr[i].open("POST",settings.url_ul+"?r="+Math.random(),!0),xhr[i].setRequestHeader("Content-Encoding","identity"),xhr[i].send(req))}}.bind(this),1)}.bind(this),i=0;i<settings.xhr_ulMultistream;i++)testStream(i,100*i);interval=setInterval(function(){var t=(new Date).getTime()-startT;if(!(200>t)){var speed=totLoaded/(t/1e3);ulStatus=(8*speed/925e3).toFixed(2),(t/1e3>settings.time_ul||failed)&&((failed||isNaN(ulStatus))&&(ulStatus="Fail"),clearRequests(),clearInterval(interval),done())}}.bind(this),200)}}function pingTest(done){if(!ptCalled){ptCalled=!0;var prevT=null,ping=0,jitter=0,i=0,prevInstspd=0;xhr=[];var doPing=function(){prevT=(new Date).getTime(),xhr[0]=new XMLHttpRequest,xhr[0].onload=function(){if(0===i)prevT=(new Date).getTime();else{var instspd=((new Date).getTime()-prevT),instjitter=Math.abs(instspd-prevInstspd);1===i?ping=instspd:(ping=.9*ping+.1*instspd,jitter=instjitter>jitter?.2*jitter+.8*instjitter:.9*jitter+.1*instjitter),prevInstspd=instspd}pingStatus=ping.toFixed(2),jitterStatus=jitter.toFixed(2),i++,i<settings.count_ping?doPing():done()}.bind(this),xhr[0].onerror=function(){pingStatus="Fail",jitterStatus="Fail",clearRequests(),done()}.bind(this),xhr[0].open("GET",settings.url_ping+"?r="+Math.random(),!0),xhr[0].send()}.bind(this);doPing()}}var testStatus=0,dlStatus="",ulStatus="",pingStatus="",jitterStatus="",clientIp="",settings={time_ul:15,time_dl:15,count_ping:35,url_dl:"garbage",url_ul:"empty",url_ping:"empty",url_getIp:"getIP",xhr_dlMultistream:10,xhr_ulMultistream:3,xhr_dlUseBlob:!1,garbagePhp_chunkSize:20,enable_quirks:!0,allow_fetchAPI:!1,force_fetchAPI:!1},xhr=null,interval=null,useFetchAPI=!1;this.addEventListener("message",function(e){var params=e.data.split(" ");if("status"===params[0]&&postMessage(testStatus+";"+dlStatus+";"+ulStatus+";"+pingStatus+";"+clientIp+";"+jitterStatus),"start"===params[0]&&0===testStatus){testStatus=1;try{var s=JSON.parse(e.data.substring(5));if("undefined"!=typeof s.url_dl&&(settings.url_dl=s.url_dl),"undefined"!=typeof s.url_ul&&(settings.url_ul=s.url_ul),"undefined"!=typeof s.url_ping&&(settings.url_ping=s.url_ping),"undefined"!=typeof s.url_getIp&&(settings.url_getIp=s.url_getIp),"undefined"!=typeof s.time_dl&&(settings.time_dl=s.time_dl),"undefined"!=typeof s.time_ul&&(settings.time_ul=s.time_ul),"undefined"!=typeof s.enable_quirks&&(settings.enable_quirks=s.enable_quirks),"undefined"!=typeof s.allow_fetchAPI&&(settings.allow_fetchAPI=s.allow_fetchAPI),settings.enable_quirks){var ua=navigator.userAgent;/Firefox.(\d+\.\d+)/i.test(ua)&&(settings.xhr_ulMultistream=1),/Edge.(\d+\.\d+)/i.test(ua)&&(settings.xhr_dlMultistream=3),/Safari.(\d+)/i.test(ua)&&!/Chrome.(\d+)/i.test(ua)&&(settings.xhr_ulMultistream=10,settings.garbagePhp_chunkSize=5),/Chrome.(\d+)/i.test(ua)&&self.fetch&&(settings.allow_fetchAPI&&(useFetchAPI=!0),settings.xhr_dlMultistream=5)}"undefined"!=typeof s.count_ping&&(settings.count_ping=s.count_ping),"undefined"!=typeof s.xhr_dlMultistream&&(settings.xhr_dlMultistream=s.xhr_dlMultistream),"undefined"!=typeof s.xhr_ulMultistream&&(settings.xhr_ulMultistream=s.xhr_ulMultistream),"undefined"!=typeof s.xhr_dlUseBlob&&(settings.xhr_dlUseBlob=s.xhr_dlUseBlob),"undefined"!=typeof s.garbagePhp_chunkSize&&(settings.garbagePhp_chunkSize=s.garbagePhp_chunkSize),"undefined"!=typeof s.force_fetchAPI&&(settings.force_fetchAPI=s.force_fetchAPI),settings.allow_fetchAPI&&settings.force_fetchAPI&&self.fetch&&(useFetchAPI=!0)}catch(e){}console.log(settings),console.log("Fetch API: "+useFetchAPI),getIp(function(){dlTest(function(){testStatus=2,pingTest(function(){testStatus=3,ulTest(function(){testStatus=4})})})})}"abort"===params[0]&&(clearRequests(),interval&&clearInterval(interval),testStatus=5,dlStatus="",ulStatus="",pingStatus="",jitterStatus="")});var dlCalled=!1,r=new ArrayBuffer(1048576);try{r=new Float32Array(r);for(var i=0;i<r.length;i++)r[i]=Math.random()}catch(e){}for(var req=[],reqsmall=[],i=0;20>i;i++)req.push(r);req=new Blob(req),r=new ArrayBuffer(262144);try{r=new Float32Array(r);for(var i=0;i<r.length;i++)r[i]=Math.random()}catch(e){}reqsmall.push(r),reqsmall=new Blob(reqsmall);var ulCalled=!1,ptCalled=!1;`
	writer.Write([]byte(name))
}

func getIp(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte(strings.Split(request.RemoteAddr, ":")[0]))
}

func empty(writer http.ResponseWriter, request *http.Request) {
	ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	writer.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0, post-check=0, pre-check=0")
	writer.Header().Set("Pragma", "no-cache")
	writer.WriteHeader(200)
	writer.Write([]byte(""))
}
func garbage(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Description", "File Transfer")
	writer.Header().Set("Content-Type", "application/octet-stream")
	writer.Header().Set("Content-Disposition", "attachment; filename=random.dat")
	writer.Header().Set("Content-Transfer-Encoding", "binary")

	writer.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0, post-check=0, pre-check=0")
	writer.Header().Set("Pragma", "no-cache")
	writer.WriteHeader(200)
	n, _ := strconv.Atoi(request.URL.Query()["ckSize"][0])
	b := make([]byte, 1048576)
	for i := range b {
		b[i] = byte(rand.Intn(255))
	}
	for i := 0; i < n; i++ {
		writer.Write(b)
	}

}
func index(writer http.ResponseWriter, request *http.Request) {
	const name = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>服务器测速</title>
    <style type="text/css">
        html,
        body {
            margin: 0;
            padding: 0;
            border: none;
            text-align: center;
            background-color:#141526;
            color: #FFF;
			font-family: 微软雅黑;
        }

        div.test {
            display: inline-block;
            margin: 1em;
            font-size: 2vw;
            min-width: 10vw;
            text-align: center;
        }

        div.testName,
        div.meterUnit {
            font-size: 1em;
        }

        div.meter {
            font-size: 1.5em;
            line-height: 2em;
            height: 2em !important;
        }

        .flash {
            animation: flash 0.6s linear infinite;
        }

        @keyframes flash {
            0% { opacity: 0.6; }
            50% { opacity: 1; }
        }

        a {
            display: inline-block;
            border: 0.15em solid #FFF;
            padding: 0.3em 0.5em;
            margin: 0.6em;
            color: #FFF;
            text-decoration: none;
        }

        #ip {
            margin: 0.8em 0;
            font-size: 1.2em;
        }

        @media all and (max-width: 50em) {
            div.test {
                font-size: 2em;
            }
        }
      .lite-text{display:inline-block;background-color:#fff38e;color:#0b0c1b;border-radius:.2em;padding:.2em .2em .2em;margin-left:.4em;font-weight:700;line-height:1;font-size:11px;font-size:1.1rem;position:relative;top:-.4em}
    </style>
    <script type="text/javascript">
        var w = null
        function runTest() {
            document.getElementById('startBtn').style.display = 'none'
            document.getElementById('testArea').style.display = ''
            document.getElementById('abortBtn').style.display = ''
            document.getElementById('ip').style.display = ''
            document.getElementById('intro').innerHTML = "当前客户端IP："
            w = new Worker('./speedtest_worker.min.js')
            var interval = setInterval(function () { w.postMessage('status') }, 100)
            w.onmessage = function (event) {
                var data = event.data.split(';')
                var status = Number(data[0])
                var dl = document.getElementById('download')
                var ul = document.getElementById('upload')
                var ping = document.getElementById('ping')
                var ip = document.getElementById('ip')
                var jitter = document.getElementById('jitter')
                dl.className = status === 1 ? 'flash' : ''
                ping.className = status === 2 ? 'flash' : ''
                jitter.className = ul.className = status === 3 ? 'flash' : ''
                if (status === 4) {
                    clearInterval(interval)
                    document.getElementById('abortBtn').style.display = 'none'
                    document.getElementById('startBtn').style.display = ''
                    document.getElementById('startBtn').innerHTML = "重新测试"
                    document.getElementById('intro').innerHTML = "当前客户端IP："
                    w = null
                }
                if (status === 5) {
                    clearInterval(interval)
                    document.getElementById('testArea').style.display = 'none'
                    document.getElementById('abortBtn').style.display = 'none'
                    document.getElementById('startBtn').style.display = ''
                    document.getElementById('startBtn').innerHTML = "开始测试"
                    document.getElementById('intro').innerHTML = "客户端对服务器网络测试"
                    document.getElementById('ip').style.display = 'none'
                }
                dl.textContent = data[1]
                ul.textContent = data[2]
                ping.textContent = data[3]
                jitter.textContent = data[5]
                ip.textContent = data[4]
            }
            w.postMessage('start')
        }
        function abortTest() {
            if (w) w.postMessage('abort')
        }
    </script>
</head>

<body>
    <br />
    <h1>Speedtest测速工具<div class="lite-text">Lite</div></h1>
  <div id="intro">客户端对服务器网络测试</div>
  <div id="ip" style="display:none">None</div><br />
  <a href="javascript:runTest()" id="startBtn">开始测试</a>
  <a href="javascript:abortTest()" style="display:none" id="abortBtn">取消测试</a>
    <div id="testArea" style="display:none">
        
        <div class="test">
            <div class="testName">下载速度</div>
            <div class="meter">&nbsp;<span id="download"></span>&nbsp;</div>
            <div class="meterUnit">Mbit/s</div>
        </div>
        <div class="test">
            <div class="testName">上传速度</div>
            <div class="meter">&nbsp;<span id="upload"></span>&nbsp;</div>
            <div class="meterUnit">Mbit/s</div>
        </div>
        <div class="test">
            <div class="testName">平均延迟</div>
            <div class="meter">&nbsp;<span id="ping"></span>&nbsp;</div>
            <div class="meterUnit">ms</div>
        </div>
        <div class="test">
            <div class="testName">延迟波动</div>
            <div class="meter">&nbsp;<span id="jitter"></span>&nbsp;</div>
            <div class="meterUnit">ms</div>
        </div>
        <br/>
        
    </div>
<!--
Localization and optimization by Jonvi
Thanks:https://github.com/adolfintel/speedtest
-->
</body>
</html>
`
	writer.Write([]byte(name))
}
