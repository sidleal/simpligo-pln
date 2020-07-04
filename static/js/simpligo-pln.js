function logout() {
    sessionStorage.setItem('simpligo.pln.jwt.key', "logout");
    window.location = '/';
}


function loadMenu(path) {
    $.ajax({
        type: 'GET',
        url: path,
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
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

function abrirRanking() {
    abrirPagina("ranking");
}

function abrirMetrix() {
    abrirPagina("nilcmetrix");
}

function abrirAnalysis() {
    abrirPagina("analysis");
}

function abrirPagina(pagina) {
    $.ajax({
        type: 'GET',
        url: '/' + pagina,
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        }
    }).done(function(data) {
        document.write(data);
        document.close();
    }).fail(function(error) {
        alert( "Erro" );
    });

}
