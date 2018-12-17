
function overSentence(sentence) {
    toggleObject(sentence, "background: #b0cfff;");
}

function outSentence(sentence) {
    toggleObject(sentence, "");
}

function overToken(token) {
    toggleObject(token, "font-weight:bold;text-decoration:underline;");
}

function outToken(token) {
    toggleObject(token, "");
}

function toggleObject(obj, style) {
    pairObjList = obj.getAttribute('data-pair').split(',');
    pairObjList.forEach( pairObjId => {
        pairObj = document.getElementById(pairObjId);
        if (pairObj != null) {
            var selected = pairObj.getAttribute('data-selected');
            if (selected == 'false') {
                pairObj.style = style;
            }
        }    
    });
}

var selectedWords = [];
var selectedSentences = [];

var markingWords = false;

function markWords(newValue) {

    markingWords = newValue;

    if (newValue) {
        document.getElementById('markWords').style='';
        document.getElementById('markSentences').style='display:none;';
    } else {
        document.getElementById('markWords').style='display:none;';
        document.getElementById('markSentences').style='';
    }
}

var operationsMap = {
    union: 'União de Sentença',
    division: 'Divisão de Sentença',
    remotion: 'Remoção de Sentença',
    inclusion: 'Inclusão de Sentença',
    rewrite: 'Reescrita de Sentença',
    lexicalSubst: 'Substituição Lexical',
    synonymListElab: 'Elaboração com lista de sinônimos',
    explainPhraseElab: 'Elaboração com oração explicativa',
    verbalTenseSubst: 'Substituição de tempo verbal',
    numericExprSimpl: 'Simplificação de expressão numérica',
    partRemotion: 'Remoção de parte da sentença',
    passiveVoiceChange: 'Mudança de voz passiva para voz ativa',
    phraseOrderChange: 'Mudança da ordem das orações',
    svoChange: 'Mudança da ordem dos constituintes para SVO',
    advAdjOrderChange: 'Mudança da ordem de adjunto adverbial',
    pronounToNoun: 'Substituição de pronome por nome',
    nounSintReduc: 'Redução de sintagma nominal',
    discMarkerChange: 'Substituição de marcador discursivo',
    definitionElab: 'Elaboração léxica com definição',
    notMapped: 'Operação Não Mapeada'
}

function sentenceClick(sentence) {
    
    if (!markingWords) {
        var selected = sentence.getAttribute('data-selected');
        
        clearSelection();

        if (selected == 'false') {
            selectSentence(sentence, 'background: #EDE981;', 'true');
            selectedSentences.push(sentence.id);
        
            var qtTokensPair = [];
            pairObjList = sentence.getAttribute('data-pair').split(',');
            pairObjList.forEach( pairObjId => {
                pairObj = document.getElementById(pairObjId);
                if (pairObj != null) {
                    qtTokensPair.push(pairObj.getAttribute("data-qtw"));
                }    
            });
            document.getElementById("qtSelectedTokens").innerHTML = sentence.getAttribute("data-qtw") + " -> " + qtTokensPair.toString();
            document.getElementById("qtSelectedTokens").title = sentence.getAttribute("data-qtw") + " palavras ( e " + sentence.getAttribute("data-qtt") + " tokens) na sentença de origem, destino: " + qtTokensPair.toString();
            refreshOperationsList(sentence);
        }
        $("#selectedSentences").val(selectedSentences.toString());
    }    
}

function clearSelection() {

    clearSentenceSelection();
    clearWordSelection();

}

function clearWordSelection() {
    
    selectedWords.forEach(w => {
        word = document.getElementById(w);
        word.style = '';
        word.setAttribute('data-selected', 'false');
    });
    selectedWords = [];
    $("#selectedWords").val('');
}

function clearSentenceSelection() {
    selectedSentences.forEach( s => {
        selectSentence(document.getElementById(s), '', 'false');
        document.getElementById("qtSelectedTokens").innerHTML = '';
        document.getElementById("qtSelectedTokens").title = "Quantidade de palavras da sentença";
        refreshOperationsList(null);
    });
    selectedSentences = []
    $("#selectedSentences").val('');
    $("#sentenceOperations").html('');


    var textToHTML = document.getElementById("divTextTo").innerHTML;
    textToHTML = textToHTML.substring(textToHTML.indexOf("<p "), textToHTML.lastIndexOf("</p>")+4);
    textToHTML = textToHTML.replace(/(<\/p>)(<p)/g, "$1|||$2");
    var paragraphs = textToHTML.split("|||");

    paragraphs.forEach(p => {
        p = p.substring(p.indexOf("<span "), p.lastIndexOf("</span>")+7);
        p = p.replace(/(<\/span>)(<span)/g, "$1|||$2");
        var sentences = p.split("|||");

        sentences.forEach(s => {
            var newS = s.replace(/style=".*?"/g, 'style=""')
            $("#divTextTo").html($("#divTextTo").html().replace(s, newS));
        });
      
    });


}

function refreshOperationsList(sentence) {
    var operationsHtml = '';
    if (sentence != null) {    
        var operations = sentence.getAttribute('data-operations');
        if (operations != '') {
            var operationsList = operations.split(";");
            operationsList.forEach( op => {
                if (op != '') {

                    var opTokens = op.split('(');
                    var opKey = opTokens[0];
                    
                    var opDesc = this.operationsMap[opKey];
                    var details = '';
        
                    if (this.substOps.indexOf(opKey) >= 0) {
                      var match = /\((.*)\|(.*)\|(.*)\|(.*)\)/g.exec(op);
                      if (match) {
                        details = match[2] + ' --> ' + match[3];
                      }
                    } else if (opKey == 'partRemotion') {
                      var match = /\((.*)\|(.*)\|(.*)\)/g.exec(op);
                      if (match) {
                        details = match[2];
                      }              
                    } else if (opKey == 'notMapped') {
                      var match = /\((.*)\|(.*)\|(.*)\|(.*)\|(.*)\)/g.exec(op);
                      if (match) {
                        details = match[4] + ': ' + match[2] + ' --> ' + match[3];
                      }              
                    }
        
                    operationsHtml += "<li data-toggle=\"tooltip\" title=\"" + details + "\">" + opDesc + "(<i>" + details + "</i>)" + " <i class=\"fa fa-trash-o \" data-toggle=\"tooltip\" title=\"Excluir\" onclick=\"undoOperation('" + op + "');\" onMouseOver=\"this.style='cursor:pointer;color:red;';\" onMouseOut=\"this.style='cursor:pointer;';\"></i>"

                }
            });
        }
    }
    $("#sentenceOperations").html(operationsHtml);
}


function updateOperationsList(sentenceId, type) {
    
    var sentence = document.getElementById(sentenceId);

    var operations = sentence.getAttribute('data-operations');
    if (type != null) {
      operations += type;
    }
    sentence.setAttribute('data-operations', operations);

    refreshOperationsList(sentence);

}


function selectSentence(sentence, style, selected) {
    sentence.style = style;
    sentence.setAttribute('data-selected', selected);

    pairObjList = sentence.getAttribute('data-pair').split(',');
    pairObjList.forEach( pairObj => {
        if (document.getElementById(pairObj) != null) {
            document.getElementById(pairObj).style = style;
            document.getElementById(pairObj).setAttribute('data-selected', selected);
        }    
    });

}

function wordClick(word, right) {
    
    if (markingWords || right) {

        //clearSentenceSelection();
        
        var selected = word.getAttribute('data-selected');
        if (selected == 'true') {
            selectedWords.forEach(w => {
                selectWord(document.getElementById(w), '', 'false');
            });
            selectedWords = [];
            $("#selectedWords").val('');
        } else {
            if (selectedWords.length > 0) {
                var firstWordIdx = document.getElementById(selectedWords[0]).getAttribute('data-idx');
                var lastWordIdx = word.getAttribute('data-idx');
                for (var i = parseInt(firstWordIdx); i < parseInt(lastWordIdx);i++) {
                    var midWordId = 'f.t.' + i;
                    if (selectedWords.indexOf(midWordId) < 0) {
                        var midWord = document.getElementById(midWordId);
                        selectWord(midWord, 'background: #edb65d;font-weight: bold;', 'true');    
                    }
                }
            }
            selectWord(word, 'background: #edb65d;font-weight: bold;', 'true');
            
        }
    }
}

function selectWord(word, style, selected) {
    word.style = style;
    word.setAttribute('data-selected', selected);
    if (selected == 'true') {
        selectedWords.push(word.id);
    }
    $("#selectedWords").val(selectedWords.toString());
}


/*----------------------------------------------*/

var corpora = [];
var simplifications = [];
var texts = [];
var simplification;
var simplificationTextFrom;

var stage;

var selectedCorpusId;
var selectedCorpusName;
 
var selectedSimplificationId;
var selectedSimplificationName;
var selectedSimplificationTextTo;

var selectedTextId;
var selectedTextTitle;
var selectedTextName;
var selectedTextLevel;

var textFrom = '';
var context = this;
var textTo = '';

var searchText;

var loggedUser;

var simplificationParsedText;


var substOps = ['lexicalSubst', 'synonymListElab', 'explainPhraseElab', 'verbalTenseSubst', 'numericExprSimpl', 'pronounToNoun', 'passiveVoiceChange', 'phraseOrderChange', 'svoChange', 'advAdjOrderChange', 'discMarkerChange', 'nounSintReduc', 'definitionElab'];

$("#operations").draggable();  
$("#selected-sentence").draggable();


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
};

function editSentenceDialog(context, operation, sentence, callback) {
    $('<div></div>').appendTo('body')
    .html('<div><textarea rows="5" cols="80" id="editSentenceText">' + sentence + '</textarea></div>')
    .dialog({
        modal: true, title: operation, zIndex: 10000, autoOpen: true,
        width: 'auto', resizable: true,
        buttons: {
            Confirmar: function () {
                var text = $('#editSentenceText').val();
                callback(context, true, text);
                $(this).dialog("close");
            },
            Cancelar: function () {
                $(this).dialog("close");
                callback(context, false, '');
            }
        },
        close: function (event, ui) {
            $(this).remove();
            callback(context, false, '');
        }
    });
};

function editSentenceDialogInclusion(context, operation, sentence, callback) {
    $('<div></div>').appendTo('body')
    .html('<div><textarea rows="5" cols="80" id="editSentenceText">' + sentence + '</textarea></div>')
    .dialog({
        modal: true, title: operation, zIndex: 10000, autoOpen: true,
        width: 'auto', resizable: true,
        buttons: {
            Antes: function () {
                var text = $('#editSentenceText').val();
                callback(context, 1, text);
                $(this).dialog("close");
            },
            Depois: function () {
              var text = $('#editSentenceText').val();
              callback(context, 2, text);
              $(this).dialog("close");
            },
            Cancelar: function () {
                $(this).dialog("close");
                callback(context, -1, '');
            }
        },
        close: function (event, ui) {
            $(this).remove();
            callback(context, false, '');
        }
    });
};

function editSentenceDialogNotMapped(context, operation, sentence, callback) {
    $('<div></div>').appendTo('body')
    .html('<div>Operação: <input type="text" id="opDesc"/><br/><textarea rows="5" cols="80" id="editSentenceText">' + sentence + '</textarea></div>')
    .dialog({
        modal: true, title: operation, zIndex: 10000, autoOpen: true,
        width: 'auto', resizable: true,
        buttons: {
            Confirmar: function () {
                var opDesc = $('#opDesc').val();
                var text = $('#editSentenceText').val();
                callback(context, true, opDesc, text);
                $(this).dialog("close");
            },
            Cancelar: function () {
                $(this).dialog("close");
                callback(context, false, '', '');
            }
        },
        close: function (event, ui) {
            $(this).remove();
            callback(context, false, '');
        }
    });
};


function filterText() {
    var title = this.searchText;
    if (title.length > 3) {
      this.texts = this.af.list('/corpora/' + this.selectedCorpusId + "/texts", {
        query: {
          orderByChild: 'title',
          startAt: title,
          endAt: title + '\uf8ff',
          limitToLast: 50
        }
      });
    }
}

function back() {
    switch(this.stage) {
      case "corpora":
        this.showMenu();
        break;
      case "newCorpus":
        this.showMenu();
        break;
      case "textMenu":
        this.listCorpora();
        break;
      case "newText":
        this.showTextMenu();
        break;
      case "texts":
        this.showTextMenu();
        break;
      case "doSimplification":
        if (this.selectedSimplificationId == null) {
          this.listTexts();          
        } else {
          listSimplifications();
        }
        $("#operations").hide();
        $("#selected-sentence").hide();
        break;        
      case "simplifications":
        this.showTextMenu();
        break;        
      default:
        loadMenu("/");
    }
}

function refresh() {

    $("#menu").hide();
    $("#corpora").hide();
    $("#newCorpus").hide();
    $("#textMenu").hide();
    $("#newText").hide();
    $("#texts").hide();
    $("#doSimplification").hide();
    $("#simplifications").hide();

    $("#" + this.stage).show();

    $("#breadcrumb").text(": : " + this.breadcrumb);

}

function showMenu() {
    this.stage = "menu";
    this.breadcrumb = "editor > menu";
    this.refresh();
}

function listCorpora() {

    $.ajax({
        type: 'GET',
        url: '/anotador/corpus/list',
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        }
    }).done(function(data) {
        var result = JSON.parse(data);
        var lista = "";
        result.list.forEach(item => {
            lista += "<a onclick=\"selectCorpus('" + item.id + "','" + item.name+"')\">";
            lista += item.name + "<br/><p>" + item.source + "</p>";
            lista += "<i class=\"fa fa-trash-o inner-button\" data-toggle=\"tooltip\" title=\"Excluir\" onclick=\"deleteCorpus('" + item.id + "');event.stopPropagation();\"></i>";
            lista += "</a>";        
        })
        $('#details-container-corpora').html(lista);
        
    }).fail(function(error) {
        alert( "Erro" );
    });


    this.stage = "corpora";
    this.breadcrumb = "editor > meus corpora";
    this.refresh();
}

function selectCorpus(corpusId, corpusName) {
    this.selectedCorpusId = corpusId;
    this.selectedCorpusName = corpusName;
    this.showTextMenu();
}

function newCorpus() {
    this.stage = "newCorpus";
    this.breadcrumb = "editor > novo córpus";
    $('#corpusName').val("");
    $('#corpusSource').val("")
    $('#corpusGenre').val("");
    this.refresh();
}

function saveCorpus() {

    $.ajax({
        type: 'POST',
        url: '/anotador/corpus/new',
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        },
        data: JSON.stringify({
            name: $('#corpusName').val(),
            source: $('#corpusSource').val(),
            genre: $('#corpusGenre').val(),
        }),
        contentType: "application/json"
    }).done(function(data) {
        showMenu();
    }).fail(function(error) {
        alert( "Erro" );
    });

}
  
function deleteCorpus(corpusId) {
    this.confirmDialog('Confirma a exclusão?', ret => {
      if (ret) {

        $.ajax({
            type: 'DELETE',
            url: '/anotador/corpus/' + corpusId,
            headers: {
                "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
            },
        }).done(function(data) {
            listCorpora();
        }).fail(function(error) {
            alert( "Erro" );
        });
         
      }
    });

}

function backToMenu() {
    this.router.navigate(['']);
}

function showTextMenu() {
    this.stage = "textMenu";
    this.breadcrumb = "editor > meus corpora > " + this.selectedCorpusName + " > textos";
    this.refresh();
}

function listTexts() {


    $.ajax({
        type: 'GET',
        url: '/anotador/corpus/' + selectedCorpusId + '/text/list',
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        }
    }).done(function(data) {
        var result = JSON.parse(data);
        var lista = "";
        result.list.forEach(item => {
            lista += "<a onclick=\"selectText('" + item.id + "','" + item.title + "','" + item.name + "'," + item.level + ")\" onmousedown=\"$('#waiting').toggle();\">";
            lista += item.name;
            lista += "<i class=\"fa fa-trash-o inner-button\" data-toggle=\"tooltip\" title=\"Excluir\" onclick=\"deleteText('" + item.id + "');event.stopPropagation();\" onmousedown=\"event.stopPropagation();\"></i>";
            lista += "<i class=\"fa fa-spinner fa-pulse fa-fw\" id='waiting' style=\"float:right;display:none;\"></i>";
            lista += "<br/><p>" + item.title + " - " + item.source + " - Nível " + item.level + "(" + item.published + ")</p>";
            lista += "</a>";
        })
        $('#details-container-texts').html(lista);
       
    }).fail(function(error) {
        alert( "Erro" );
    });


    this.stage = "texts";
    this.breadcrumb = "editor > meus corpora > " + this.selectedCorpusName + " > textos";
    this.refresh();
}

function newText() {
    this.stage = "newText";
    this.breadcrumb = "editor > meus corpora > " + this.selectedCorpusName + " > importar novo texto";
    this.textName = '';
    this.textTitle = '';
    this.textSubTitle = '';
    this.textAuthor = '';
    this.textPublished = '';
    this.textSource = '';
    this.textContent = '';
    this.textRawContent = '';
    this.refresh();
}

function deleteText(textId) {
    this.confirmDialog('Confirma a exclusão?', ret => {
      if (ret) {

        $.ajax({
            type: 'DELETE',
            url: '/anotador/corpus/' + selectedCorpusId + "/text/" + textId,
            headers: {
                "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
            },
        }).done(function(data) {
            listTexts();
        }).fail(function(error) {
            alert( "Erro" );
        });

      }
    });
    this.refresh();
}

function saveText() {

    // if (this.textRawContent == null) {
    //   this.textRawContent = $("#textContent").val();
    // }

    this.textTitle = $("#textTitle").val();
    this.textSubTitle = $("#textSubTitle").val();
    this.textContent = $("#textContent").val();

    var textContentFull = "";
    textContentFull += "# " + this.textTitle + "\n";
    textContentFull += "## " + this.textSubTitle + "\n";
    textContentFull += this.textContent;     
    var parsedText = splitText(textContentFull);

    $.ajax({
        type: 'POST',
        url: '/anotador/corpus/' + selectedCorpusId + '/text/new',
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        },
        data: JSON.stringify({
            corpusId: selectedCorpusId,
            name: $("#textName").val(),
            title: $("#textTitle").val(),
            subTitle: $("#textSubTitle").val(),
            author: $("#textAuthor").val(),
            published: $("#textPublished").val(),
            source: $("#textSource").val(),
            content: $("#textContent").val(),
            parsed: parsedText,
            level: 0
        }),
        contentType: "application/json"
    }).done(function(data) {
        showTextMenu();
    }).fail(function(error) {
        alert( "Erro" );
    });

}


function selectText(textId, textTitle, textName, textLevel) {
    this.selectedTextId = textId;
    this.selectedTextTitle = textTitle;
    this.selectedTextName = textName;
    this.selectedTextLevel = textLevel;    
    this.selectedSimplificationId = null;
    this.doSimplification();
}

function listSimplifications() {


    $.ajax({
        type: 'GET',
        url: '/anotador/corpus/' + selectedCorpusId + '/simpl/list',
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        }
    }).done(function(data) {
        var result = JSON.parse(data);
        var lista = "";
        result.list.forEach(item => {
            lista += "<a onclick=\"selectSimplification('" + item.id + "','" + item.name + "','" + item.to + "')\" onmousedown=\"$('#waiting').toggle();\">";
            lista += item.title + ' - ' + item.name;
            lista += "<i class=\"fa fa-trash-o inner-button\" data-toggle=\"tooltip\" title=\"Excluir\" onclick=\"deleteSimplification('" + item.id + "','" + item.to + "');event.stopPropagation();\" onmousedown=\"event.stopPropagation();\"></i>";
            lista += "<i class=\"fa fa-spinner fa-pulse fa-fw\" id='waiting' style=\"float:right;display:none;\"></i>";
            lista += "<br/><p>" + item.tags + " - " + item.updated + ")</p>";
            lista += "</a>";
        })
        $('#details-container-simpl').html(lista);
       
    }).fail(function(error) {
        alert( "Erro" );
    });
    
    this.stage = "simplifications";
    this.breadcrumb = "editor > meus corpora > " + this.selectedCorpusName + " > Simplificações";
    this.refresh();
}
   
function deleteSimplification(simpId, textToId) {
    this.confirmDialog('Confirma a exclusão?', ret => {
      if (ret) {

        $.ajax({
            type: 'DELETE',
            url: '/anotador/corpus/' + selectedCorpusId + "/text/" + textToId,
            headers: {
                "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
            },
        }).done(function(data) {
            //ok          
        }).fail(function(error) {
            alert( "Erro" );
        });

        $.ajax({
            type: 'DELETE',
            url: '/anotador/corpus/' + selectedCorpusId + "/simpl/" + simpId,
            headers: {
                "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
            },
        }).done(function(data) {
            
            listSimplifications();

        }).fail(function(error) {
            alert( "Erro" );
        });

      }
    });

}

function saveSimplification() {
    
    $("#waiting").show();

    var regexp = /> #\s(.*?)<\/span>/g
    var match = regexp.exec($("#divTextTo").html());
    var textToTitle = match[1];

    regexp = /> ##\s(.*?)<\/span>/g
    match = regexp.exec($("#divTextTo").html());
    var textToSubTitle = match[1];

    var parsedParagraphs = [];
    var idxParagraphs = 0;
    var idxSentences = 0;
    var idxTokens = 0;
    var idxWords = 0;
    var newTextcontent = '';
    var qtRemoved = 0;

    var textToHTML = document.getElementById("divTextTo").innerHTML;
    textToHTML = textToHTML.substring(textToHTML.indexOf("<p "), textToHTML.lastIndexOf("</p>")+4);
    textToHTML = textToHTML.replace(/(<\/p>)(<p)/g, "$1|||$2");
    var paragraphs = textToHTML.split("|||");
    paragraphs.forEach(p => {
      var pContent = '';
      idxParagraphs++;
      var parsedParagraph = {"idx": idxParagraphs, "sentences": []};

      p = p.substring(p.indexOf("<span "), p.lastIndexOf("</span>")+7);
      p = p.replace(/(<\/span>)(<span)/g, "$1|||$2");
      var sentences = p.split("|||");
      sentences.forEach(s => {
        idxSentences++;
        var regexp = /id="(.+?)"/g
        var match = regexp.exec(s);
        var sId = match[1];

        regexp = /<span.+?>(.+?)<\/span>/g
        match = regexp.exec(s);
        var sContent = match[1];
        if (sContent != '#rem#') {
          newTextcontent += sContent;
          pContent += sContent;
        } else {
          qtRemoved++;
        }
        regexp = /data-pair="(.+?)"/g
        match = regexp.exec(s);
        var sPair = match[1];

        var sPairList = sPair.split(',');
        var operations = '';
        sPairList.forEach(pair => {
          if (document.getElementById(pair) != null) {
            var newOperations = document.getElementById(pair).getAttribute("data-operations").split(';');
            for (var i in newOperations) {
              var newOperation = newOperations[i];
              if (newOperation != '' && operations.indexOf(newOperation) < 0) {
                operations += newOperation + ';';
              }
            }            
          }
        });
        
        var parsedSentence = {"idx": idxSentences, "text": sContent, "id": sId, "pair": sPair, "operations": operations, "tokens": []};
        
        var parsedText = splitText(sContent);
        var parsedS = parsedText['paragraphs'][0]['sentences'][0];
        parsedSentence["qtt"] = parsedS["qtt"];
        parsedSentence["qtw"] = parsedS["qtw"];
        parsedS["tokens"].forEach(t => {
          idxTokens++;
          parsedSentence['tokens'].push({"idx": idxTokens, "token": t["token"]});
          if (t["token"].length > 1 || '{[()]}.,"?!;:-\'#'.indexOf(t["token"]) < 0) {
            idxWords++;
          }
        });

        parsedParagraph.sentences.push(parsedSentence);
      });
      parsedParagraph['text'] = pContent;
      parsedParagraphs.push(parsedParagraph);
      newTextcontent += '\n';
    });

    this.simplificationParsedText = {"totP": idxParagraphs, "totS": idxSentences - qtRemoved, "totT": idxTokens, "totW": idxWords, "paragraphs": parsedParagraphs};


    var newName = selectedTextName;
    newName = newName.replace(/nível_[0-9]+/g, "nível_" + (selectedTextLevel + 1));

    if (this.selectedSimplificationId != null) {

        $.ajax({
            type: 'DELETE',
            url: '/anotador/corpus/' + selectedCorpusId + "/text/" + selectedSimplificationTextTo,
            headers: {
                "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
            },
        }).done(function(data) {
            //ok
        }).fail(function(error) {
            alert( "Erro" );
        });


        $.ajax({
            type: 'DELETE',
            url: '/anotador/corpus/' + selectedCorpusId + "/simpl/" + selectedSimplificationId,
            headers: {
                "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
            },
        }).done(function(data) {
            //ok
        }).fail(function(error) {
            alert( "Erro" );
        });
        
    } 

    var parsedText = splitText(newTextcontent);

    $.ajax({
        type: 'POST',
        url: '/anotador/corpus/' + selectedCorpusId + '/text/new',
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        },
        data: JSON.stringify({
            corpusId: selectedCorpusId,
            name: newName,
            title: textToTitle,
            subTitle: textToSubTitle,
            author: loggedUser,
            published: moment().format("YYYY-MM-DD"),
            updated: moment().format("YYYY-MM-DD"),
            source: 'Simplificação Nível ' + (selectedTextLevel + 1),
            content: newTextcontent,
            parsed: parsedText,
            level: selectedTextLevel + 1
        }),
        contentType: "application/json"
    }).done(function(data) {
        saveOperationList(data)
    }).fail(function(error) {
        alert( "Erro" );
    });


}

function saveOperationList(textToId) {

    var simplSentences = [];

    this.simplificationParsedText.paragraphs.forEach(p => {
      p.sentences.forEach(s => {

        var from = s.pair.replace(/f\.s\./g, '');
        var operations = s.operations.replace(/f\.s\./g, '');
        if (operations == '') {
          operations = 'none();'
        }
        
        simplSentences.push(
          {
            from: from,
            to: s.idx.toString(),
            operations: operations
          }
        );
      });
    });

    $.ajax({
        type: 'POST',
        url: '/anotador/corpus/' + selectedCorpusId + '/simpl/new',
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        },
        data: JSON.stringify({
            corpusId: selectedCorpusId,
            name: $('#simplificationName').val(),
            title: selectedTextTitle,
            from: selectedTextId,
            to: textToId,
            tags: $('#simplificationTag').val(),
            updated: moment().format("YYYY-MM-DD"),
            sentences: simplSentences

        }),
        contentType: "application/json"
    }).done(function(data) {
        
        console.log(data);
        selectSimplification(data, $('#simplificationName').val(), textToId);

    }).fail(function(error) {
        alert( "Erro" );
    });

}

function selectSimplification(simplId, simplName, simplTextTo) {
    this.selectedSimplificationId = simplId;
    this.selectedSimplificationName = simplName;
    this.selectedSimplificationTextTo = simplTextTo;
    
    editSimplification();
}

function doSimplification() {
    
    this.stage = "doSimplification";
    this.breadcrumb = "editor > meus corpora > " + this.selectedCorpusName + " > textos > " + this.selectedTextTitle + " > Nova Simplificação";

    this.refresh();


    $.ajax({
        type: 'GET',
        url: '/anotador/corpus/' + selectedCorpusId + '/text/' + selectedTextId,
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        }
    }).done((data) => {
        var text = JSON.parse(data);

        simplificationName = "Natural " + (text.level) + ' -> ' + (text.level + 1);
        simplificationTag = "Nível " + (text.level + 1);
    
        $('#simplificationName').val(simplificationName);
        $('#simplificationTag').val(simplificationTag);

        $('#qtSentencesFrom').text(text.parsed.totP + ' - ' + text.parsed.totS + ' - ' + text.parsed.totW );
        $('#qtSentencesTo').text(text.parsed.totP + ' - ' + text.parsed.totS + ' - ' + text.parsed.totW );

        textFrom = parseTextFromOut(text, null);
        textTo = parseTextToOut(text, null);
    
        $('#divTextFrom').html(textFrom);
        $('#divTextTo').html(textTo);

        $('#operations').show();
        $('#selected-sentence').show();
  

    }).fail(function(error) {
        alert( "Erro" );
    });

}

function editSimplification() {
    
    $.ajax({
        type: 'GET',
        url: '/anotador/corpus/' + selectedCorpusId + '/simpl/' + selectedSimplificationId,
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        }
    }).done((data) => {
        var simpl = JSON.parse(data);


        $.ajax({
            type: 'GET',
            url: '/anotador/corpus/' + selectedCorpusId + '/text/' + simpl.from,
            headers: {
                "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
            }
        }).done((data) => {
            var textFrom = JSON.parse(data);
    

            $.ajax({
                type: 'GET',
                url: '/anotador/corpus/' + selectedCorpusId + '/text/' + simpl.to,
                headers: {
                    "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
                }
            }).done((data) => {
                var textTo = JSON.parse(data);
        
                
                selectedTextId = simpl.from;
                selectedTextTitle = textFrom.title;
                selectedTextName = textFrom.name;
                selectedTextLevel = textFrom.level;
                
                editSimplificationText(textFrom, textTo, simpl);
          

            }).fail(function(error) {
                alert( "Erro" );
            });
    
    
        }).fail(function(error) {
            alert( "Erro" );
        });


    }).fail(function(error) {
        alert( "Erro" );
    });





}

function editSimplificationText(textFrom, textTo, simp) {
    
    this.stage = "doSimplification";
    this.breadcrumb = "editor > meus corpora > " + this.selectedCorpusName + " > textos > " + this.selectedTextTitle + " > Editar Simplificação";

    this.refresh();

    $("#sentenceOperations").html('');

    $('#simplificationName').val(simp.name);
    $('#simplificationTag').val(simp.tags);

    $('#qtSentencesFrom').text(textFrom.parsed.totP + ' - ' + textFrom.parsed.totS + ' - ' + textFrom.parsed.totW );
    $('#qtSentencesTo').text(textTo.parsed.totP + ' - ' + textTo.parsed.totS + ' - ' + textTo.parsed.totW );

    parsedTextFrom = parseTextFromOut(textFrom, simp);
    parsedTextTo = parseTextToOut(textTo, simp);

    $('#divTextFrom').html(parsedTextFrom);
    $('#divTextTo').html(parsedTextTo);

    $('#operations').show();
    $('#selected-sentence').show();

    $("#waiting").hide();  

  }
 
function getSimplificationSentences(simp, sentenceId, source) {
    var ret = [];
    for (var s in simp.sentences) {
      var simpSentence = simp.sentences[s];
      listIds = simpSentence[source].split(",");
      if (listIds.indexOf(sentenceId.toString()) >= 0) {
        ret.push(simpSentence);
      }
    }
    return ret;
}

function parseTextFromOut(textFrom, simp) {
    var out = '';
    out += "<style type='text/css'>";
    out += " p span:hover {background:#cdff84;cursor:pointer;}";
    out += " p span span {display:inline-block; margin:0px;}";
    out += " p span span:hover {font-weight:bold;text-decoration:underline;cursor:pointer;}";
    out += "</style>";
    var openQuotes = false;
    var lastToken = '';

    for(var p in textFrom.parsed.paragraphs) {
      var par = textFrom.parsed.paragraphs[p];
      out += '<p id=\'f.p.' + par.idx + '\'>';
      for(var s in par.sentences) {
        var sObj = par.sentences[s];

        if (sObj.text != '#rem#') {

          var po = getPairAndOperations(simp, sObj.idx, 'to'); 
          
          out += '<span id=\'f.s.' + sObj.idx + '\' data-selected=\'false\' data-pair=\'' + po['pair'] + '\'';
          out += ' data-qtt=\'' + sObj.qtt + '\' data-qtw=\'' + sObj.qtw + '\'';
          out += ' data-operations=\'' + po['operations'] + '\'';
          out += ' onclick=\'sentenceClick(this)\'';
          out += ' onmouseover=\'overSentence(this);\' onmouseout=\'outSentence(this);\'>'
          for(var t in sObj.tokens) {
            var token = sObj.tokens[t].token;
            var idx = sObj.tokens[t].idx;
  
            if ('\"\''.indexOf(token) >= 0) {
              openQuotes = !openQuotes;
              if (openQuotes) {
                out += ' ';
              }
            } else if (openQuotes && '\"\''.indexOf(lastToken) >= 0) {
              //nothing
            } else if ('.,)]}!?:'.indexOf(token) < 0 && '([{'.indexOf(lastToken) < 0) {
              out += ' ';
            }
  
            out += '<span id=\'f.t.' + idx + '\' data-selected=\'false\' data-pair=\'t.t.' + idx + '\'';
            out += ' data-idx=\'' + idx + '\'';            
            out += ' onclick=\'wordClick(this, false)\'';
            out += ' oncontextmenu=\'wordClick(this, true); return false;\'';
            out += ' onmouseover=\'overToken(this);\' onmouseout=\'outToken(this);\'>' + token + '</span>';
            lastToken = token;
           }
          out += ' </span>';
          
        }

      }
      out += "</p>"
    }

    return out; 
}

function getPairAndOperations(simp, s, source) {
    var ret = {};
    
    var prefix = source[0] + '.s.';

    var inverseSource = '';
    if (source == 'to') {
      inverseSource = 'from';
    } else {
      inverseSource = 'to';
    }
    
    var pair = '';
    var operations = '';
    
    if (simp != null) {
      var newPairList = [];
      var simpSentences = getSimplificationSentences(simp, s, inverseSource);
      for (var j in simpSentences) {
        var pairList = simpSentences[j][source].split(','); 
        for (var i in pairList) {
          if(pairList[i] != '') {
            newPairList.push(prefix + pairList[i]);                
          }
        }
        var oper = simpSentences[j].operations;
        if (!oper.startsWith('none') && operations.indexOf(oper) < 0 && oper.indexOf(s) > 0) {
            operations += oper;  
        }  
      }
      pair = newPairList.toString();      
    } else {
      pair = prefix + s;
    }
    ret['pair'] = pair;
    ret['operations'] = operations;

    return ret;
}

function parseTextToOut(textTo, simp) {
    var out = '';
    out += "<style type='text/css'>";
    out += " p span:hover {background:#cdff84;cursor:text;}";
    out += "</style>";
    var openQuotes = false;
    for(var p in textTo.parsed.paragraphs) {
      var par = textTo.parsed.paragraphs[p];
      out += '<p id=\'t.p.' + par.idx + '\'>';
      for(var s in par.sentences) {
        var sObj = par.sentences[s];
        
        if (simp == null && sObj.text == '#rem#') {
          //ignore
        } else {

            var po = this.getPairAndOperations(simp, sObj.idx, 'from'); 

            out += '<span id=\'t.s.' + sObj.idx + '\'  data-selected=\'false\' data-pair=\'' + po['pair'] + '\'';
            out += ' data-qtt=\'' + sObj.qtt + '\' data-qtw=\'' + sObj.qtw + '\'';
            out += ' data-operations=\'' + po['operations'] + '\'';
            out += ' onmouseover=\'overSentence(this);\' onmouseout=\'outSentence(this);\'>';
            // TODO: alterar pela function criada.
            for(var t in sObj.tokens) {
              var token = sObj.tokens[t].token;
              if ('\"\''.indexOf(token) >= 0) {
                openQuotes = !openQuotes;
                if (openQuotes) {
                  out += ' ';
                }
              } else if (openQuotes && '\"\''.indexOf(out.substr(-1)) >= 0) {
                //nothing
              } else if ('.,)]}!?:'.indexOf(token) < 0 && '([{'.indexOf(out.substr(-1)) < 0) {
                out += ' ';
              }
              out += token;
            }
            out += ' </span>';

        }
      }
      out += "</p>"
    }
    return out;
}


function changeListener(event) {
    
    var reader = new FileReader();
    var fileTokens = event.target.value.split('\\')[2].split('.');
    var fileName = fileTokens[0] + ' nível_0';

    this.textName = fileName;
    this.textTitle = '';
    this.textSubTitle = '';
    this.textAuthor = '';
    this.textPublished = '';
    this.textSource = '';
    
    reader.onloadend = (e) => {
        this.textRawContent = reader.result;

        var linhas = reader.result.split("\n");
        var parsedText = '';

        linhas.forEach(element => {
          var meta = false;
          
          var match;
          
          match = /<title>(.*)<\/title>/g.exec(element);
          if (match) {
            this.textTitle = match[1];
            meta = true;
          }

          match = /<subtitle>(.*)<\/subtitle>/g.exec(element);
          if (match) {
            if (!this.textSubTitle) {
              this.textSubTitle = match[1];
            } else {
              parsedText += "\n" + match[1] + "\n\n";
            }
            meta = true;
          }

          match = /<author>(.*)<\/author>/g.exec(element);
          if (match) {
            this.textAuthor = match[1];
            meta = true;
          }

          match = /<date>(.*)<\/date>/g.exec(element);
          if (match) {
            this.textPublished = match[1];
            meta = true;
          }

          match = /<url>(.*)<\/url>/g.exec(element);
          if (match) {
            this.textSource = match[1];
            meta = true;
          }

          if (!meta) {
            parsedText += element + "\n";
          }

        });

        this.textContent = parsedText;

        $("#textName").val(fileName);
        $("#textTitle").val(textTitle);
        $("#textSubTitle").val(textSubTitle);
        $("#textAuthor").val(textAuthor);
        $("#textPublished").val(textPublished);
        $("#textSource").val(textSource);
        $("#textContent").val(textContent);

    }
    reader.readAsText(event.target.files[0], 'UTF-8');

}

// OPERATIONS

function doOperation(type) {

    var selectedSentences = $('#selectedSentences').val().split(',');
    var selectedSentence = selectedSentences[0]; // apenas uma sentença (talvez mais de uma no futuro)

    var selectedWords = $('#selectedWords').val().split(',');
    if (selectedWords.length > 0 && selectedWords[0] != '' ) {
      selectedSentence = document.getElementById(selectedWords[0]).parentNode.attributes['id'].value;
    }

    this.rewriteTextTo(type, selectedSentence, selectedWords);

    $('#operations-sentenciais').hide();
    $('#operations-intra-sentenciais').hide();
    $('#hideOps').hide();
    $('#showOps').show();

}

function rewriteTextTo(type, selectedSentence, selectedWords) {
      
    switch (type) {
      //intra-sentenciais
      case 'lexicalSubst':
          this.doLexicalSubst(selectedSentence, selectedWords); break;
      case 'synonymListElab':
          this.doSynonymListElab(selectedSentence, selectedWords); break;
      case 'explainPhraseElab':
          this.doExplainPhraseElab(selectedSentence, selectedWords); break;
      case 'verbalTenseSubst':
          this.doVerbalTenseSubst(selectedSentence, selectedWords); break;
      case 'numericExprSimpl':
          this.doNumericExprSimpl(selectedSentence, selectedWords); break;
      case 'partRemotion':
          this.doPartRemotion(selectedSentence, selectedWords); break;
      case 'passiveVoiceChange':
          this.doPassiveVoiceChange(selectedSentence, selectedWords); break;
      case 'phraseOrderChange':
          this.doPhraseOrderChange(selectedSentence, selectedWords); break;
      case 'svoChange':
          this.doSVOChange(selectedSentence, selectedWords); break;
      case 'advAdjOrderChange':
          this.doAdvAdjOrderChange(selectedSentence, selectedWords); break;
      case 'pronounToNoun':
          this.doPronounToNoun(selectedSentence, selectedWords); break;
      case 'nounSintReduc':
          this.doNounSintReduc(selectedSentence, selectedWords); break;
      case 'discMarkerChange':
          this.doDiscMarkerChange(selectedSentence, selectedWords); break;
      case 'notMapped':
          this.doNotMapped(selectedSentence, selectedWords); break;
      case 'definitionElab':
          this.doDefinitionElab(selectedSentence, selectedWords); break;
    }
    
    var textToHTML = document.getElementById("divTextTo").innerHTML;
      
    textToHTML = textToHTML.substring(textToHTML.indexOf("<p "), textToHTML.lastIndexOf("</p>")+4);
    textToHTML = textToHTML.replace(/(<\/p>)(<p)/g, "$1|||$2");
    var paragraphs = textToHTML.split("|||");

    paragraphs.forEach(p => {
        p = p.substring(p.indexOf("<span "), p.lastIndexOf("</span>")+7);
        p = p.replace(/(<\/span>)(<span)/g, "$1|||$2");
        var sentences = p.split("|||");
        
        switch (type) {
            case 'union':
                this.doUnion(sentences, selectedSentence); break;
            case 'division':
                this.doDivision(sentences, selectedSentence); break;
            case 'remotion':
                this.doRemotion(sentences, selectedSentence); break;
            case 'inclusion':
                this.doInclusion(sentences, selectedSentence); break;
            case 'rewrite':
                this.doRewrite(sentences, selectedSentence); break;
        }

    });
}



function parseSentence(sentence) {
      var ret = {};

      var regexp = /ngcontent-([^=]+)/g;
      var match = regexp.exec(sentence);
      if (match != null) {
        ret['ngContent'] = match[1];         
      } else {
        ret['ngContent'] = 'aa';
      }

      regexp = /data-pair="(.+?)"/g;
      match = regexp.exec(sentence);
      ret['pair'] = match[1];

      regexp = /data-qtt="(.+?)".*data-qtw="(.+?)"/g;
      match = regexp.exec(sentence);
      ret['qtt'] = parseInt(match[1]);
      ret['qtw'] = parseInt(match[2]);

      regexp = /id="(.+?)"/g;
      match = regexp.exec(sentence);
      ret['id'] = match[1];

      regexp = />(.+?)<\/span>/g;
      match = regexp.exec(sentence);    
      ret['content'] = match[1];

      return ret;
}

function doUnion(sentences, selectedSentence) {
    var previousSentence = '';
    sentences.forEach(s => {
        if (s.indexOf("pair=\"" + selectedSentence + "\"") > 0) {
            var ps = this.parseSentence(s);
            var pps = this.parseSentence(previousSentence);
  
            this.editSentenceDialog(this, "União de Sentença", pps['content'] + ' ' + ps['content'], function (context, ret, text) {
                if (ret) {
                    var parsedText = splitText(text);
                    var parsedSentence = parsedText['paragraphs'][0]['sentences'][0]; 

                    var newId = pps['id'] + '|' + ps['id'];

                    var newSentHtml = "<span data-pair=\"{pair}\" data-qtt=\"{qtt}\" data-qtw=\"{qtw}\" data-selected=\"true\" id=\"{id}\" onmouseout=\"outSentence(this);\" onmouseover=\"overSentence(this);\" style=\"font-weight: bold;background: #EDE981;\"> {content}</span>";
                    newSentHtml = newSentHtml.replace("{id}", newId);
                    newSentHtml = newSentHtml.replace("{pair}", pps['pair'] + ',' + ps['pair']);
                    newSentHtml = newSentHtml.replace("{qtt}", parsedSentence['qtt']);
                    newSentHtml = newSentHtml.replace("{qtw}", parsedSentence['qtw']);
                    newSentHtml = newSentHtml.replace("{content}", parsedSentence['text']);

                    $("#divTextTo").html($("#divTextTo").html().replace(document.getElementById(pps['id']).outerHTML, ''));
                    $("#divTextTo").html($("#divTextTo").html().replace(s, newSentHtml));

                    document.getElementById(ps['pair']).setAttribute('data-pair', newId);
                    document.getElementById(pps['pair']).setAttribute('data-pair', newId);
                  
                    context.updateOperationsList(ps['pair'], 'union(' + ps['pair'] + '|' +  pps['pair'] + ');');
                    context.updateOperationsList(pps['pair'], 'union(' + ps['pair'] + '|' +  pps['pair'] + ');');
                }
  
            });
        }
        previousSentence = s;
    });
}


function doDivision(sentences, selectedSentence) {

    sentences.forEach(s => {
        if (s.indexOf("pair=\"" + selectedSentence + "\"") > 0) {
            var ps = this.parseSentence(s);

           this.editSentenceDialog(this, "Divisão de Sentença", ps['content'], function (context, ret, text) {
                if (ret) {
                    var parsedText = splitText(text);
                    var parsedSentences = parsedText['paragraphs'][0]['sentences']; 

                    var newHtml = "";
                    var newIds = [];
                    parsedSentences.forEach(ns => {
                      var newSentHtml = "<span data-pair=\"{pair}\" data-qtt=\"{qtt}\" data-qtw=\"{qtw}\" data-selected=\"true\" id=\"{id}\" onmouseout=\"outSentence(this);\" onmouseover=\"overSentence(this);\" style=\"font-weight: bold;background: #EDE981;\"> {content}</span>";
                      newSentHtml = newSentHtml.replace("{id}", ps['id'] + '_new_' + ns['idx']);
                      newSentHtml = newSentHtml.replace("{pair}", ps['pair']);
                      newSentHtml = newSentHtml.replace("{qtt}", ns['qtt']);
                      newSentHtml = newSentHtml.replace("{qtw}", ns['qtw']);
                      newSentHtml = newSentHtml.replace("{content}", ns['text']);

                      newIds.push(ps['id'] + '_new_' + ns['idx']);
                      newHtml += newSentHtml;
                    });
                   
                    $("#divTextTo").html($("#divTextTo").html().replace(s, newHtml));
                    document.getElementById(selectedSentence).setAttribute('data-pair', newIds.toString());

                    context.updateOperationsList(selectedSentence, 'division(' + selectedSentence + ');');
                }
  
            });
                
        }
    });
}


function doInclusion(sentences, selectedSentence) {
    sentences.forEach(s => {
        if (s.indexOf("pair=\"" + selectedSentence + "\"") > 0) {
            var ps = this.parseSentence(s);
  
            this.editSentenceDialogInclusion(this, "Inclusão de Sentença", '', function (context, ret, text) {
                if (ret > 0) {
                    var parsedText = splitText(text);
                    var parsedSentence = parsedText['paragraphs'][0]['sentences'][0]; 

                    var newId = ps['id'] + '_new';

                    var newSentHtml = "<span data-pair=\"{pair}\" data-qtt=\"{qtt}\" data-qtw=\"{qtw}\" data-selected=\"true\" id=\"{id}\" onmouseout=\"outSentence(this);\" onmouseover=\"overSentence(this);\" style=\"font-weight: bold;background: #EDE981;\"> {content}</span>";
                    newSentHtml = newSentHtml.replace("{id}", newId);
                    newSentHtml = newSentHtml.replace("{pair}", ps['pair']);
                    newSentHtml = newSentHtml.replace("{qtt}", parsedSentence['qtt']);
                    newSentHtml = newSentHtml.replace("{qtw}", parsedSentence['qtw']);
                    newSentHtml = newSentHtml.replace("{content}", parsedSentence['text']);

                    var newHtml = '';
                    if (ret == 1) {
                      newHtml = newSentHtml + s;
                    } else {
                      newHtml = s + newSentHtml;
                    }
                    
                    $("#divTextTo").html($("#divTextTo").html().replace(s, newHtml));

                    document.getElementById(selectedSentence).setAttribute('data-pair', ps['id'] + ',' + newId);
                    context.updateOperationsList(selectedSentence, 'inclusion(' + selectedSentence + '|' + ret + ');');
                }
  
            });
        }
    });
}

function doRewrite(sentences, selectedSentence) {
    sentences.forEach(s => {
        if (s.indexOf("pair=\"" + selectedSentence + "\"") > 0) {
            var ps = this.parseSentence(s);

            this.editSentenceDialog(this, "Reescrita de Sentença", ps['content'], function (context, ret, text) {
                if (ret) {
                    var parsedText = splitText(text);
                    var ns = parsedText['paragraphs'][0]['sentences'][0]; 

                    var newSentHtml = "<span data-pair=\"{pair}\" data-qtt=\"{qtt}\" data-qtw=\"{qtw}\" data-selected=\"true\" id=\"{id}\" onmouseout=\"outSentence(this);\" onmouseover=\"overSentence(this);\" style=\"font-weight: bold;background: #EDE981;\"> {content}</span>";
                    newSentHtml = newSentHtml.replace("{id}", ps['id']);
                    newSentHtml = newSentHtml.replace("{pair}", ps['pair']);
                    newSentHtml = newSentHtml.replace("{qtt}", ns['qtt']);
                    newSentHtml = newSentHtml.replace("{qtw}", ns['qtw']);
                    newSentHtml = newSentHtml.replace("{content}", ns['text']);
                   
                    $("#divTextTo").html($("#divTextTo").html().replace(s, newSentHtml));

                    context.updateOperationsList(selectedSentence, 'rewrite(' + selectedSentence + ');');
                }
  
            });
                
        }
    });
}

function doRemotion(sentences, selectedSentence) {
    sentences.forEach(s => {
            if (s.indexOf("pair=\"" + selectedSentence + "\"") > 0) {
            var ps = this.parseSentence(s);
  
            $("#divTextTo").html($("#divTextTo").html().replace(s, s.replace(ps['content'], '#rem#')));

            this.updateOperationsList(selectedSentence, 'remotion(' + selectedSentence + ');');
        }
    });


}


function parseWordTokens(tokens) {
    var openQuotes = false;
    var out = '';
    for(var i in tokens) {

      var token = tokens[i];
      if ('\"\''.indexOf(token) >= 0) {
        openQuotes = !openQuotes;
        if (openQuotes) {
          out += ' ';
        }
      } else if (openQuotes && '\"\''.indexOf(out.substr(-1)) >= 0) {
        //nothing
      } else if ('.,)]}!?:'.indexOf(token) < 0 && '([{'.indexOf(out.substr(-1)) < 0) {
        out += ' ';
      }
      out += token;
    }
    // tokens = tokens.substring(0, tokens.length - 1);    
    return out;
}

function parseSelectedWordTokens(selectedWords) {
    var tokens = [];
    for (var i in selectedWords) {
      var word = document.getElementById(selectedWords[i]).innerText;
      tokens.push(word);
    }
    return this.parseWordTokens(tokens);
}

function doIntraSentenceSubst(selectedSentence, selectedWords, operation) {
    var opDesc = this.operationsMap[operation];
    var parsedTokens = this.parseSelectedWordTokens(selectedWords);
    var firstTokenIdx = document.getElementById(selectedWords[0]).getAttribute('data-idx');

    this.editSentenceDialog(this, opDesc, parsedTokens, function (context, ret, text) {
        if (ret) {
            context.updateOperationsList(selectedSentence, operation + '(' + selectedSentence + '|' + parsedTokens + '|' + text + '|' + firstTokenIdx + ');');
        }

    });

}

function doLexicalSubst(selectedSentence, selectedWords) {
    this.doIntraSentenceSubst(selectedSentence, selectedWords, 'lexicalSubst');
}

function doSynonymListElab(selectedSentence, selectedWords) {
    this.doIntraSentenceSubst(selectedSentence, selectedWords, 'synonymListElab');
}

function doExplainPhraseElab(selectedSentence, selectedWords) {
    this.doIntraSentenceSubst(selectedSentence, selectedWords, 'explainPhraseElab');
}

function doVerbalTenseSubst(selectedSentence, selectedWords) {
    this.doIntraSentenceSubst(selectedSentence, selectedWords, 'verbalTenseSubst');
}

function doNumericExprSimpl(selectedSentence, selectedWords) {
    this.doIntraSentenceSubst(selectedSentence, selectedWords, 'numericExprSimpl');
}

function doPassiveVoiceChange(selectedSentence, selectedWords) {
    this.doIntraSentenceSubst(selectedSentence, selectedWords, 'passiveVoiceChange');
}

function doPhraseOrderChange(selectedSentence, selectedWords) {
    this.doIntraSentenceSubst(selectedSentence, selectedWords, 'phraseOrderChange');
}

function doSVOChange(selectedSentence, selectedWords) {
    this.doIntraSentenceSubst(selectedSentence, selectedWords, 'svoChange');
}

function doAdvAdjOrderChange(selectedSentence, selectedWords) {
    this.doIntraSentenceSubst(selectedSentence, selectedWords, 'advAdjOrderChange');
}

function doPronounToNoun(selectedSentence, selectedWords) {
    this.doIntraSentenceSubst(selectedSentence, selectedWords, 'pronounToNoun');
}

function doNounSintReduc(selectedSentence, selectedWords) {
    this.doIntraSentenceSubst(selectedSentence, selectedWords, 'nounSintReduc');
}

function doDiscMarkerChange(selectedSentence, selectedWords) {
    this.doIntraSentenceSubst(selectedSentence, selectedWords, 'discMarkerChange');
}

function doPartRemotion(selectedSentence, selectedWords) {
    var parsedTokens = this.parseSelectedWordTokens(selectedWords);
    var firstTokenIdx = document.getElementById(selectedWords[0]).getAttribute('data-idx');
    
    this.updateOperationsList(selectedSentence, 'partRemotion(' + selectedSentence + '|' + parsedTokens + '|' + firstTokenIdx+ ');');
}

function doNotMapped(selectedSentence, selectedWords) {
    var parsedTokens = this.parseSelectedWordTokens(selectedWords);
    var firstTokenIdx = document.getElementById(selectedWords[0]).getAttribute('data-idx');

    this.editSentenceDialogNotMapped(this, this.operationsMap["notMapped"], parsedTokens, function (context, ret, opDesc, text) {
        if (ret) {
            context.updateOperationsList(selectedSentence, 'notMapped(' + selectedSentence + '|' + parsedTokens + '|' + text + '|' + opDesc + '|' + firstTokenIdx + ');');
        }

    });

}

function doDefinitionElab(selectedSentence, selectedWords) {
    this.doIntraSentenceSubst(selectedSentence, selectedWords, 'definitionElab');
  }


function undoOperation(operation) {

    var opTokens = operation.split('(');
    var opKey = opTokens[0];
    var sentenceId = opTokens[1].split('|')[0];

    if (!sentenceId.startsWith('f.s.')) {
      sentenceId = 'f.s.' + sentenceId;
    }

    if (sentenceId.indexOf(')') > 0) {
      sentenceId = sentenceId.substring(0, sentenceId.indexOf(')'));
    }

    var textToHTML = document.getElementById("divTextTo").innerHTML;
    
    textToHTML = textToHTML.substring(textToHTML.indexOf("<p "), textToHTML.lastIndexOf("</p>")+4);
    textToHTML = textToHTML.replace(/(<\/p>)(<p)/g, "$1|||$2");
    var paragraphs = textToHTML.split("|||");

    if (opKey == 'rewrite') {
      this.undoRewrite(sentenceId, operation);

    } else if (this.substOps.indexOf(opKey) >= 0 || 'partRemotion,notMapped'.indexOf(opKey) >= 0 ) {
      this.undoIntraSentence(sentenceId, operation);

    } else {

      paragraphs.forEach(p => {
          p = p.substring(p.indexOf("<span "), p.lastIndexOf("</span>")+7);
          p = p.replace(/(<\/span>)(<span)/g, "$1|||$2");
          var sentences = p.split("|||");

          if (opKey == 'inclusion') {
            this.undoInclusion(sentences, sentenceId, operation);
          } else if (opKey == 'division') {
            this.undoDivision(sentences, sentenceId, operation);
          } else if (opKey == 'remotion') {
            this.undoRemotion(sentences, sentenceId, operation);
          } else if (opKey == 'union') {
            this.undoUnion(sentences, sentenceId, operation);
          }

      });

    }

}

function undoIntraSentence(sentenceId, operation) {
    var fromSentence = document.getElementById(sentenceId);
    var operations = fromSentence.getAttribute('data-operations');
    fromSentence.setAttribute('data-operations', operations.replace(operation, ''));
    this.updateOperationsList(sentenceId, null);

  }

function undoInclusion(sentences, sentenceId, operation) {
    var fromSentence = document.getElementById(sentenceId);
    var fromPair = fromSentence.getAttribute("data-pair");
    var includedSentence = fromPair.split(',')[1];

    sentences.forEach(s => {
      if (s.indexOf(includedSentence) > 0) {
          $("#divTextTo").html($("#divTextTo").html().replace(s, ''));
          
          var operations = fromSentence.getAttribute('data-operations');
          fromSentence.setAttribute('data-operations', operations.replace(operation, ''));
          this.updateOperationsList(sentenceId, null);

      }
     
    });

}

function undoDivision(sentences, sentenceId, operation) {
    var fromSentence = document.getElementById(sentenceId);
    var fromPair = fromSentence.getAttribute("data-pair");
    var pairSentences = fromPair.split(',');

    var unitedSentence = "";
    sentences.forEach(s => {
      pairSentences.forEach(pair => {
        if (s.indexOf("id=\"" + pair + "\"") > 0) {    
          var ps = this.parseSentence(s);
          unitedSentence += ps['content'];
          pairSentences.splice(pairSentences.indexOf(pair), 1);

          if (pairSentences.length > 0) {
            $("#divTextTo").html($("#divTextTo").html().replace(s, ''));        
          } else {
            var newId = sentenceId.replace('f.s', 't.s');

            var newSentHtml = "<span data-pair=\"{pair}\" data-qtt=\"{qtt}\" data-qtw=\"{qtw}\" data-selected=\"true\" id=\"{id}\" onmouseout=\"outSentence(this);\" onmouseover=\"overSentence(this);\" style=\"font-weight: bold;background: #EDE981;\"> {content}</span>";
            newSentHtml = newSentHtml.replace("{id}", newId);
            newSentHtml = newSentHtml.replace("{pair}", sentenceId);
            newSentHtml = newSentHtml.replace("{qtt}", ps['qtt']);
            newSentHtml = newSentHtml.replace("{qtw}", ps['qtw']);
            newSentHtml = newSentHtml.replace("{content}", unitedSentence);
    
            $("#divTextTo").html($("#divTextTo").html().replace(s, newSentHtml));
            var operations = fromSentence.getAttribute('data-operations');
            fromSentence.setAttribute('data-operations', operations.replace(operation, ''));
            this.updateOperationsList(sentenceId, null);
            document.getElementById(sentenceId).setAttribute('data-pair', newId);
          }
    
        }
      });

    });

}

function undoUnion(sentences, sentenceId, operation) {
    var fromSentence = document.getElementById(sentenceId);
    var fromPair = fromSentence.getAttribute("data-pair");

    sentences.forEach(s => {
        if (s.indexOf("id=\"" + fromPair + "\"") > 0) {        
            var ps = this.parseSentence(s);

            this.editSentenceDialog(this, "Desfazer União de Sentença", ps['content'], function (context, ret, text) {
                if (ret) {
                    var parsedText = splitText(text);
                    var parsedSentences = parsedText['paragraphs'][0]['sentences']; 

                    var toPairs = ps['pair'].split(',');

                    var fromSentence1 = document.getElementById(toPairs[0]);
                    var fromSentence2 = document.getElementById(toPairs[1]);

                    var newHtml = "";

                    for (var i = 0; i < 2; i++) {
                      var newSentHtml = "<span data-pair=\"{pair}\" data-qtt=\"{qtt}\" data-qtw=\"{qtw}\" data-selected=\"true\" id=\"{id}\" onmouseout=\"outSentence(this);\" onmouseover=\"overSentence(this);\" style=\"font-weight: bold;background: #EDE981;\"> {content}</span>";
                      newSentHtml = newSentHtml.replace("{id}", toPairs[i].replace('f.s.', 't.s.'));
                      newSentHtml = newSentHtml.replace("{pair}", toPairs[i]);
                      newSentHtml = newSentHtml.replace("{qtt}", parsedSentences[i]['qtt']);
                      newSentHtml = newSentHtml.replace("{qtw}", parsedSentences[i]['qtw']);
                      newSentHtml = newSentHtml.replace("{content}", parsedSentences[i]['text']);
                      newHtml += newSentHtml;    
                    }
                   
                    $("#divTextTo").html($("#divTextTo").html().replace(s, newHtml));

                    var operations = fromSentence1.getAttribute('data-operations');
                    fromSentence1.setAttribute('data-operations', operations.replace(operation, ''));
                    operations = fromSentence2.getAttribute('data-operations');
                    fromSentence2.setAttribute('data-operations', operations.replace(operation, ''));
                    context.updateOperationsList(toPairs[0], null);
                    context.updateOperationsList(toPairs[1], null);

                    fromSentence1.setAttribute('data-pair', toPairs[0].replace('f.s.', 't.s.'));
                    fromSentence2.setAttribute('data-pair', toPairs[1].replace('f.s.', 't.s.'));

                  }
  
            });
                
        }
    });
}


function undoRewrite(sentenceId, operation) {
    var fromSentence = document.getElementById(sentenceId);

    var operations = fromSentence.getAttribute('data-operations');
    fromSentence.setAttribute('data-operations', operations.replace(operation, ''));
    this.updateOperationsList(sentenceId, null);

}

function undoRemotion(sentences, sentenceId, operation) {
    var fromSentence = document.getElementById(sentenceId);
    var fromPair = fromSentence.getAttribute("data-pair");

    sentences.forEach(s => {
      if (s.indexOf(fromPair) > 0) {
          $("#divTextTo").html($("#divTextTo").html().replace(s, s.replace('#rem#', fromSentence.innerText)));
          
          var operations = fromSentence.getAttribute('data-operations');
          fromSentence.setAttribute('data-operations', operations.replace(operation, ''));
          this.updateOperationsList(sentenceId, null);

      }
    });
}
