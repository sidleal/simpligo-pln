function eval() {

    $.ajax({
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
        alert(xhr.responseText);
    });

  }