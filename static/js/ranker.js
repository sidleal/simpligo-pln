function eval() {

    ws.send(JSON.stringify({ content: $('#content').val() }));

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
    return;
} else {
    ws = initWS();
}

function initWS() {
    var socket = new WebSocket("ws://simpligo.sidle.al:443/ranker/ws");

    socket.onopen = function() {
        console.log("<p>Socket is open</p>");
    };
    socket.onmessage = function (e) {
        console.log("<p> Got some shit:" + e.data + "</p>");
    }
    socket.onclose = function () {
        console.log("<p>Socket closed</p>");
    }
    return socket;
}
