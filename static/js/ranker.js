function eval() {
    $('#output').val(""); 
    ws.send(JSON.stringify({ 
        content: $('#content').val(),
        auth: sessionStorage.getItem('simpligo.pln.jwt.key') 
    }));

/*    $.ajax({
        type: 'POST',
        url: '/ranker/eval',
        timeout: 30000,
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        },
        data: {
            content: $('#content').val()
        }
    }).done(function(data) {
        $('#results').show();
        $('#output').val(data);        
    }).fail(function(xhr, status, error) {
        console.log(xhr.responseText);
        console.log(error);
        alert(error);
    });
*/
  }


var ws;

if (window.WebSocket === undefined) {
    $("#container").append("Your browser does not support WebSockets");
} else {
    ws = initWS();
}

function initWS() {
    var socket = new WebSocket("wss://simpligo.sidle.al:443/ranker/ws");

    socket.onopen = function() {
        console.log("Socket is open.");
    };
    socket.onmessage = function (e) {
        console.log("Got:" + e.data);
        var result = JSON.parse(e.data);

        $('#results').show();
        $('#output').val(result.raw_result);   

    }
    socket.onclose = function () {
        console.log("Socket closed.");
    }
    return socket;
}
