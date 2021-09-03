var websocket
var createWebSocketClient = function(connectHandle,messageHandle,closeHandle){
    websocket = new WebSocket("ws://"+window.location.host+"/ws")
    websocket.onopen=connectHandle
    websocket.onmessage = messageHandle
    websocket.onclose=closeHandle
    websocket.binaryType = "arraybuffer"
    return websocket
}
var sendMessage = function(msg){
    websocket.send(msg)
}

function encodeUtf8(text) {
    const code = encodeURIComponent(text);
    const bytes = [];
    for (var i = 0; i < code.length; i++) {
        const c = code.charAt(i);
        if (c === '%') {
            const hex = code.charAt(i + 1) + code.charAt(i + 2);
            const hexVal = parseInt(hex, 16);
            bytes.push(hexVal);
            i += 2;
        } else bytes.push(c.charCodeAt(0));
    }
    return bytes;
}


function byteToString(arr) {
    if(typeof arr === 'string') {
        return arr;
    }
    var str = '',
        _arr = arr;
    for(var i = 0; i < _arr.length; i++) {
        var one = _arr[i].toString(2),
            v = one.match(/^1+?(?=0)/);
        if(v && one.length == 8) {
            var bytesLength = v[0].length;
            var store = _arr[i].toString(2).slice(7 - bytesLength);
            for(var st = 1; st < bytesLength; st++) {
                store += _arr[st + i].toString(2).slice(2);
            }
            str += String.fromCharCode(parseInt(store, 2));
            i += bytesLength - 1;
        } else {
            str += String.fromCharCode(_arr[i]);
        }
    }
    return str;
}

function doReceive(buffer) {
    var receive = [];
    receive = receive.concat(Array.from(new Uint8Array(buffer)));
    if (receive.length < 4) {
        return;
    }
    return receive
}



// init websocket client

var connectHandle = function(e){
    console.log("websocket connected")
    data = '{"type":"container_log","value":"123456789","interval":1}'
    console.log('get containerId log:'+data)
    websocket.send(data)
}
var messageHandle = function(e){
    var res = doReceive(e.data)
    var str = byteToString(res)
    console.log("recv message: ", str)
} 
var closeHandle = function(e){
    console.log("websocket closed")
}

websocket = createWebSocketClient(connectHandle,messageHandle,closeHandle)