function newCloze() {
    this.stage = "newCloze";
    this.breadcrumb = "Cloze > novo teste";
    $('#clozeName').val("");
    $('#clozeCode').val("");
    $('#textContent').val("")
    this.refresh();
}


function refresh() {

    $("#menu").hide();
    $("#newCloze").hide();
    $("#clozeList").hide();
    $("#paragraphs").hide();

    $("#" + this.stage).show();

    $("#breadcrumb").text(": : " + this.breadcrumb);

}

function back() {
    switch(this.stage) {
      case "newCloze":
        this.showMenu();
        break;
      case "clozeList":
        this.showMenu();
        break;
      case "paragraphs":
        this.showMenu();
        break;
      default:
        loadMenu("/");
    }
}

function saveCloze() {

    $.ajax({
        type: 'POST',
        url: '/cloze/new',
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        },
        data: JSON.stringify({
            name: $('#clozeName').val(),
            code: $('#clozeCode').val(),
            content: $('#textContent').val(),
        }),
        contentType: "application/json"
    }).done(function(data) {
        showMenu();
    }).fail(function(error) {
        alert( "Erro" );
    });

}

function showMenu() {
    this.stage = "menu";
    this.breadcrumb = "cloze > menu";
    this.refresh();
}

function listCloze() {

    $.ajax({
        type: 'GET',
        url: '/cloze/list',
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        }
    }).done(function(data) {
        var result = JSON.parse(data);
        var lista = "";
        result.list.forEach(item => {
            lista += "<a onclick=\"clozeDetails('" + item.id + "','" + item.name + "')\">";
            lista += item.name + "<br/><p>" + item.code + "</p>";
            lista += "<i class=\"fa fa-trash-o inner-button\" data-toggle=\"tooltip\" title=\"Excluir\" onclick=\"deleteCloze('" + item.id + "');event.stopPropagation();\"></i>";
            lista += "</a>";        
        })
        $('#details-container-cloze').html(lista);
        
    }).fail(function(error) {
        alert( "Erro" );
    });

    this.stage = "clozeList";
    this.breadcrumb = "cloze > lista";
    this.refresh();
}


function clozeDetails(clozeId, clozeName) {

    $.ajax({
        type: 'GET',
        url: '/cloze/' + clozeId,
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        }
    }).done(function(data) {
        var result = JSON.parse(data);
        var lista = "";
        details = "<a>Total: " + result.parsed.totp + " Parágrafos, "  + result.parsed.tots + " Sentenças, " + result.parsed.totw + " Palavras" + "</a>";
        lista += details;

        result.parsed.paragraphs.forEach(item => {
            lista += "<a>";
            lista += "Parágrafo " + item.idx + " - " + item.txt;
            lista += "<br/><p> "       
            item.sentences.forEach(s => {
                lista += "Sentença " + s.idx + ": <b>" + s.qtw + "</b> palavras --> " + s.txt + "<br/>";
            });
            lista += "</p></a>";
        })
        $('#details-container-paragraphs').html(lista);
       
    }).fail(function(error) {
        alert( "Erro" );
    });

    this.stage = "paragraphs";
    this.breadcrumb = "cloze > " + clozeName + " > detalhes";
    this.refresh();
}
