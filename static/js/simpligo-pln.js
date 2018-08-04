function logout() {
    sessionStorage.setItem('simpligo.pln.jtw.key', null);
    window.location = '/';
}


function loadMenu() {
    $.ajax({
        type: 'GET',
        url: '/',
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jtw.key')
        }
    }).done(function(data) {
        document.write(data);
        document.close();
    }).fail(function(error) {
        alert( "Erro" );
    });

}

function abrirSenter() {
    abrirPagina("senter");
}

function abrirCloze() {
    abrirPagina("cloze");
}

function abrirPalavras() {
    abrirPagina("palavras");
}

function abrirAnotador() {
    abrirPagina("anotador");
}

function abrirPagina(pagina) {
    $.ajax({
        type: 'GET',
        url: '/' + pagina,
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jtw.key')
        }
    }).done(function(data) {
        document.write(data);
        document.close();
    }).fail(function(error) {
        alert( "Erro" );
    });

}
