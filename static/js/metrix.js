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
        // console.log(resData);

        $('#results').show();
        $('#output').val(resData.raw); 

        table = "<table style='width:80%;margin-left:auto;margin-right:auto;'><tr><td><b>Nome</b></td><td><b>Valor</b></td></tr>";
        resData.list.forEach(item => {
            table += "<tr><td><a href='/nilcmetrixdoc#" + item.name + "' target='_blank'>"+item.name+"</a></td><td>"+item.val+"</td></tr>";
        })
        table += "</table>";
        $('#table_results').html(table); 

        document.getElementById("btnparse").style = "cursor:pointer;";      
    }).fail(function(error) {
        console.log(error);
        document.getElementById("btnparse").style = "cursor:pointer;";      
        alert( "Erro" );
    });

  }