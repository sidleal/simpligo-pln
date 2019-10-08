function newAnalysis() {
    this.stage = "newAnalysis";
    this.breadcrumb = "Analisador de parágrafos > nova análise";
    $('#analysisTitle').val("");
    $('#textContent').val("")
    this.refresh();
}


function refresh() {

    $("#menu").hide();
    $("#newAnalysis").hide();
    $("#analysisList").hide();
    $("#paragraphs").hide();

    $("#" + this.stage).show();

    $("#breadcrumb").text(": : " + this.breadcrumb);

}

function back() {
    switch(this.stage) {
      case "newAnalysis":
        this.showMenu();
        break;
      case "analysisList":
        this.showMenu();
        break;
      case "paragraphs":
        this.showMenu();
        break;
      default:
        loadMenu("/");
    }
}

function saveAnalysis() {

    $.ajax({
        type: 'POST',
        url: '/analysis/new',
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        },
        data: JSON.stringify({
            title: $('#analysisTitle').val(),
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
    this.breadcrumb = "Analisador de parágrafos > menu";
    this.refresh();
}


function deleteCloze(clozeID) {
    this.confirmDialog('Confirma a exclusão?', ret => {
      if (ret) {

        $.ajax({
            type: 'DELETE',
            url: '/cloze/' + clozeID,
            headers: {
                "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
            },
        }).done(function(data) {
            listCloze();
        }).fail(function(error) {
            alert( "Erro" );
        });
         
      }
    });

}

function confirmDialog(message, callback) {
    $('<div></div>').appendTo('body')
    .html('<div><h6>'+message+'?</h6></div>')
    .dialog({
        modal: true, title: 'Atenção', zIndex: 10000, autoOpen: true,
        width: 'auto', resizable: false,
        buttons: {
            Yes: function () {
                $(this).dialog("close");
                callback(true);
            },
            No: function () {
                $(this).dialog("close");
                callback(false);
            }
        },
        close: function (event, ui) {
            $(this).remove();
        },
        open: function() {
            $('.ui-dialog :button').blur();
        }
    });
}


function listAnalysis() {

    $.ajax({
        type: 'GET',
        url: '/analysis/list',
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        }
    }).done(function(data) {
        var result = JSON.parse(data);
        var lista = "";
        result.list.forEach(item => {
            lista += "<a onclick=\"analysisDetails('" + item.id + "','" + item.title + "')\">";
            lista += item.title + "<br/><p>" + item.id + "</p>";
            lista += "<i class=\"fa fa-trash-o inner-button\" data-toggle=\"tooltip\" title=\"Excluir\" onclick=\"deleteAnalysis('" + item.id + "');event.stopPropagation();\"></i>";
            // lista += "<i style=\"margin-right:15px;\" class=\"fa fa-save inner-button\" data-toggle=\"tooltip\" title=\"Exportar resultado\" onclick=\"exportAnalysis('" + item.id + "');event.stopPropagation();\"></i>";
            lista += "</a>";        
        })
        $('#details-container-analysis').html(lista);
        
    }).fail(function(error) {
        alert( "Erro" );
    });

    this.stage = "analysisList";
    this.breadcrumb = "analysis > lista";
    this.refresh();
}


function analysisDetails(analysisId, analysisTitle) {

    $.ajax({
        type: 'GET',
        url: '/analysis/' + analysisId,
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        }
    }).done(function(data) {
        // console.log(data);
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
    this.breadcrumb = "Analisador de parágrafos > " + analysisTitle + " > detalhes";
    this.refresh();
}



// function exportAnalysis(analysisId) {

//     $.ajax({
//         type: 'GET',
//         url: '/analysis/export/' + analysisId,
//         headers: {
//             "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
//         }
//     }).done(function(data) {
//         var link = document.createElement("a");
//         link.href = 'data:text/csv,' + encodeURIComponent(data);
//         link.download = "analysis_" + clozeId + "_data.csv";
//         link.click();

//     }).fail(function(error) {
//         alert( "Erro ao exportar." );
//     });

// }
