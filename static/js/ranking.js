function eval() {
    $('#output').val(""); 
    ws.send(JSON.stringify({ 
        content: $('#content').val(),
        options:$('input[name=options]:checked').val(), 
        auth: sessionStorage.getItem('simpligo.pln.jwt.key') 
    }));
}


var ws;

if (window.WebSocket === undefined) {
    $("#container").append("Your browser does not support WebSockets");
} else {
    ws = initWS();
}

function initWS() {
    var socket = new WebSocket("wss://simpligo.sidle.al:443/ranking/ws");

    socket.onopen = function() {
        console.log("Socket is open.");
    };
    socket.onmessage = function (e) {
//        console.log("Got:" + e.data);
        var result = JSON.parse(e.data);

        $('#results').show();
        $('#output').val(result.raw_result);   
        document.getElementById("btneval").style = "cursor:pointer;";

    }
    socket.onclose = function () {
        console.log("Socket closed.");
    }
    return socket;
}
