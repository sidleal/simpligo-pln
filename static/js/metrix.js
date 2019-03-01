function parse() {

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
        $('#results').show();
        $('#output').val(data); 
        document.getElementById("btnparse").style = "cursor:pointer;";      
    }).fail(function(error) {
        console.log(error);
        document.getElementById("btnparse").style = "cursor:pointer;";      
        alert( "Erro" );
    });

  }