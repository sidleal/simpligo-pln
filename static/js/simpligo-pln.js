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
    $.ajax({
        type: 'GET',
        url: '/senter',
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