function eval() {

    $.ajax({
        type: 'POST',
        url: '/ranker/eval',
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        },
        data: {
            content: $('#content').val()
        }
    }).done(function(data) {
        $('#results').show();
        $('#output').val(data);        
    }).fail(function(error) {
        console.log(error);
        alert( "Erro" );
    });

  }