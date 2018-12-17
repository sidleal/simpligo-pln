function parseFlat() {
    this.parse("flat");
}

function parseTree() {
    this.parse("tree");
}

function parse(returnType) {

    $.ajax({
        type: 'POST',
        url: '/palavras/parse',
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        },
        data: {
            type: returnType, 
            options:$('input[name=options]:checked').val(), 
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