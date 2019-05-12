function parse() {
    $('#output').val(""); 
    $.ajax({
        type: 'POST',
        url: '/metrix/parse',
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        },
        data: {
            options:$('input[name=options]:checked').val(), 
            content: $('#content').val()
        }
    }).done(function(data) {
        resData = JSON.parse(data);
        console.log(resData);

        $('#results').show();
        $('#output').val(resData.raw); 
        document.getElementById("btnparse").style = "cursor:pointer;";      
    }).fail(function(error) {
        console.log(error);
        document.getElementById("btnparse").style = "cursor:pointer;";      
        alert( "Erro" );
    });

  }