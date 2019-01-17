var wsUri;
var output;
var msgInput;
var ws;

window.addEventListener("load", function(evt) {
    wsUri  = "ws://" + "localhost:9070" + "/stream"
    output = document.getElementById("output");
    msgInput = "udin"

  var print = function(message) {
    var d       = document.createElement("div");
    d.innerHTML = message;
    output.appendChild(d);
  };

  var parseData = function(evt) {
    return JSON.parse(evt.data).msg
  };

  var newSocket = function() {
    ws           = new WebSocket(wsUri);
    ws.onopen = function(evt) {
      print('<span style="color: green;">Connection Open</span>');
    }
    ws.onclose = function(evt) {
      print('<span style="color: red;">Connection Closed</span>');
      ws = null;
    }
    ws.onmessage = function(evt) {
      print('<span style="color: blue;">Update: </span>' + parseData(evt));
    }
    ws.onerror = function(evt) {
      print('<span style="color: red;">Error: </span>' + parseData(evt));
    }
  };

//   newSocket()

  document.getElementById("send").onclick = function(evt) {
    msgInput = document.getElementById("msgParam").value;

    msgInput = (msgInput == "") ? "udin" : msgInput
    var msgParam = { Msg: msgInput }

    req = JSON.stringify(msgParam)

    if (!ws) {
        console.log("failed", ws)
      return false
    }
    
    ws.send(req);

    return false;
  };

  document.getElementById("cancel").onclick = function(evt) {
    if (!ws) {
      return false;
    }
    ws.close();
    print('<span style="color: red;">Request Canceled</span>');
    return false;
  };

  document.getElementById("open").onclick = function(evt) {
    if (!ws) {
      newSocket()
    }
    return false;
  };

  document.getElementById("close-conn").onclick = function(evt) {
   if (ws) {
       ws.close()
   }
  };
})
