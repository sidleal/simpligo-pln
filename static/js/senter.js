$('#newAbbrev').keypress(function (e) {
    if (e.which == 13) {
       saveAbbrev();
       $('#newAbbrev').val('');
    }
  });

loadAbbreviations();

function showAbbreviations() {
    $('#btnShowAbbrev').hide();
    $('#btnHideAbbrev').show();
    $('#abbreviations').show(); 
}

function hideAbbreviations() {
    $('#btnShowAbbrev').show();
    $('#btnHideAbbrev').hide();
    $('#abbreviations').hide(); 
}

var abbreviations = []

function saveAbbrev() {

    $.ajax({
        type: 'POST',
        url: '/senter/abbrev/new',
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        },
        data: JSON.stringify({name: $('#newAbbrev').val()}),
        contentType: "application/json"
    }).done(function(data) {
        loadAbbreviations()
    }).fail(function(error) {
        alert( "Erro" );
    });

}

function loadAbbreviations() {
    $.ajax({
        type: 'GET',
        url: '/senter/abbrev/list',
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        }
    }).done(function(data) {
        var result = JSON.parse(data);
        var lista = "";
        abbreviations = [];
        result.list.forEach(item => {
            abbreviations.push(item.name);
            lista += item.name + "<span><a onclick='removeAbbrev(\"" + item.id + "\")'><i class='fa fa-trash-o inline-inner-button' data-toggle='tooltip' title='Excluir'></i> </a></span>"
        })
        $('#listaAbbrev').html(lista);
        
    }).fail(function(error) {
        alert( "Erro" );
    });
}

function removeAbbrev(abbrevKey) {
    $.ajax({
        type: 'DELETE',
        url: '/senter/abbrev/' + abbrevKey,
        headers: {
            "Authorization": sessionStorage.getItem('simpligo.pln.jwt.key')
        },
    }).done(function(data) {
        loadAbbreviations()
    }).fail(function(error) {
        alert( "Erro" );
    });

}

/* --------------------------------- */

var totalParagraphs;
var totalSentences;
var totalTokens;
var rawOutput = '';
var xmlOutput = '';
var jsonOutput = '';
var jsonOutputNoFormat = '';

function backToMenu() {
  loadMenu("/");
}

function split() {
    $('#results').show();

    var parsedText = splitText($('#rawContent').val());
  
    this.totalParagraphs = parsedText['totP'];
    this.totalSentences =  parsedText['totS'];
    this.totalTokens =  parsedText['totT'];
  
    this.rawOutput = this.outputToRaw(parsedText);
    this.xmlOutput = this.outputToXML(parsedText);
    this.jsonOutput = this.outputToJSON(parsedText);
    this.jsonOutputNoFormat = this.outputToJSONNoFormat(parsedText);
    
    $('#rawOutput').html(rawOutput);
    $('#xmlOutput').val(xmlOutput);
    $('#jsonOutput').val(jsonOutput);
    $('#jsonOutputNoFormat').val(jsonOutputNoFormat);

    $('#totalParagraphs').text(totalParagraphs);
    $('#totalSentences').text(totalSentences);
    $('#totalTokens').text(totalTokens);

}


function outputToRaw(parsedText) {
  var out = '';
  parsedText['paragraphs'].forEach(p => {
    p['sentences'].forEach(s => {
      out += 's' + s['idx'] + ' - ' + s['text'] + '<br/>';
    });
  });
  return out;
}

function outputToXML(parsedText) {
  var out = '';
  out = '<text>\n';
  parsedText['paragraphs'].forEach(p => {
    out += '  <p i=\"' + p['idx'] + '\">\n';
    p['sentences'].forEach(s => {
      out += '    <s i=\"' + s['idx'] + '\">' + s['text'] + "</s>\n";          
    });
    out += '  </p>\n';
  });
  out += '<totP>' + parsedText['totP'] + '</totP>\n'
  out += '<totS>' + parsedText['totS'] + '</totS>\n'
  out += '</text>';
  return out;
}

function outputToJSON(parsedText) {
  var out = '';
  out = '{\n  \"text\": [\n';
  parsedText['paragraphs'].forEach(p => {
    out += '            { \"i\": ' + p['idx'] + ',\n';
    out += '              \"p\":\"' + p['text'].replace(/"/g,'\\"') + '\",\n';
    out += '              \"ps\":[\n'
    p['sentences'].forEach(s => {
      out += '                   { \"i\": ' + s['idx'] + ', \"qtw\": ' + s['qtw'] + ', \"qtt\": ' + s['qtt'] + ', \"s\": \"' + s['text'].replace(/"/g,'\\"');
      out += '\",\n                     \"st\": [\n';
      s['tokens'].forEach(t => {
        out += '                             { \"i\": ' + t['idx'] + ', \"t\": \"' + t['token'] + '\"},\n';
      });
      out += '                     ]\n';
      out += '                   },\n';
    });
    out += '              ]\n            },\n';
  });
  out += '  ],\n';
  out += '  \"totP\":' + parsedText['totP'] + ',\n';
  out += '  \"totS\":' + parsedText['totS'] + ',\n';
  out += '  \"totT\":' + parsedText['totT'] + ',\n';
  out += '}';
  return out;
}

function outputToJSONNoFormat(parsedText) {
  var out = this.outputToJSON(parsedText);
  out = out.replace(/\n\s+/g, '');
  return out;
}

/* --------------------------------- */

/*
Rules:
1. Delimita-se uma sentença sempre que uma marca de nova linha (carriage return e line
feed) é encontrada, independentemente de um sinal de fim de sentença ter sido encontrado
anteriormente;
2. Não se delimitam sentenças dentro de aspas, parênteses, chaves e colchetes;
3. Delimita-se uma sentença quando os símbolos de interrogação (?) e exclamação (!) são
encontrados;
4. Delimita-se uma sentença quando o símbolo de ponto (.) é encontrado e este não é um
ponto de número decimal, não pertence a um símbolo de reticências (...), não faz parte de
endereços de e-mail e páginas da Internet e não é o ponto que segue uma abreviatura;
5. Delimita-se uma sentença quando uma letra maiúscula é encontrada após o sinal de
reticências ou de fecha-aspas.
(Pardo, 2006)
*/
function preProcesText(rawText) {
    var out = rawText;
    var match;

    // rule 2 - " { [ ( ) ] } "
    out = this.applyGroupRule(out, /"([^"]+?)"[^A-z]/g)
    out = this.applyGroupRule(out, /“(.+?)”/g)
    out = this.applyGroupRule(out, /\{(.+?)\}/g)
    out = this.applyGroupRule(out, /\[(.+?)\]/g)
    out = this.applyGroupRule(out, /\((.+?)\)/g)

    // rule 4 - abbreviations
    this.abbreviations.forEach(abbrev => {
        var abbrevNew = abbrev.replace(/\./g, '|dot|');
        var abbrevRe = abbrev.replace(/\./g, '\\.');
        var re = new RegExp(abbrevRe, "g");
        out = out.replace(re, abbrevNew);        
    });

    // rule 4 - internet
    var regexAddress = /(http|ftp|www|@)(.+?)\s/g;
    while (match = regexAddress.exec(out)) {
        var addressOld = match[2];
        var addressNew = addressOld.replace(/\./g, '|dot|');
        out = out.replace(addressOld, addressNew);
    }
    var regexEmail = /([A-z0-9_\-\.]+)@/g;
    while (match = regexEmail.exec(out)) {
        var emailOld = match[1];
        var emailNew = emailOld.replace(/\./g, '|dot|');
        out = out.replace(emailOld, emailNew);
    }

    //rule 4 - decimals
    out = out.replace(/([0-9]+)\.([0-9]+)/g, '$1|dot|$2');

    //rule 5 - quotes
    //out = out.replace(/"\s+([A-Z])/g, '\"|dot| |||$1'); // in texts well written this rule may be disabled.

    //rule 5 - reticences
    out = out.replace(/\.\.\.\s*([A-Z])/g, '|dot||dot||dot| |||$1');

    //rule 4 - reticences
    out = out.replace(/\.\.\./g, '|dot||dot||dot|');

    // rule 3
    out = out.replace(/\./g, '.|||');
    out = out.replace(/\?/g, '?|||');
    out = out.replace(/\!/g, '!|||');


    return out;
}

function applyGroupRule(rawText, regexGroup) {
    var match;
    while (match = regexGroup.exec(rawText)) {
        var sentenceOld = match[1];
        var sentenceNew = sentenceOld.replace(/\./g, '|gdot|');
        sentenceNew = sentenceNew.replace(/\?/g, '|gint|');
        sentenceNew = sentenceNew.replace(/\!/g, '|gexc|');
        rawText = rawText.replace(sentenceOld, sentenceNew);
    }
    return rawText;
}

function parseText(rawText) {
    var parsedText = {};
    var paragraphs = rawText.split("\n");

    var idxParagraphs = 0;
    var idxSentences = 0;
    var idxTokens = 0;
    var idxWords = 0;

    parsedText['paragraphs'] = [];

    paragraphs.forEach(p => {
        p = p.trim();
        if (p != '') {
        idxParagraphs++;
        var parsedParagraph = {"idx": idxParagraphs, "sentences": [], "text": p };
        var sentences = p.split("|||");
        sentences.forEach(s => {
            s = s.trim();
            if (s.length > 1) {
            idxSentences++;
            var parsedSentence = {"idx": idxSentences, "tokens": [], "text": s, "qtt": 0, "qtw": 0};
            var tokens = this.tokenizeText(s);
            var qtw = 0;
            var qtt = 0;
            tokens.forEach(t => {
                if (t.length > 0) {
                idxTokens++;
                qtt++;
                parsedSentence['tokens'].push({"idx": idxTokens, "token": t});
                if (t.length > 1 || '{[()]}.,"?!;:-\'#'.indexOf(t) < 0) {
                    qtw++;
                    idxWords++;
                }
                }
            });
            parsedSentence['qtt'] = qtt;
            parsedSentence['qtw'] = qtw;

            parsedParagraph['sentences'].push(parsedSentence);
            }
        });
        parsedText['paragraphs'].push(parsedParagraph);
        }
    });
    parsedText['totP'] = idxParagraphs;
    parsedText['totS'] = idxSentences;
    parsedText['totT'] = idxTokens;
    parsedText['totW'] = idxWords;

    return parsedText;
}

function postProcess(parsedText) {

    parsedText['paragraphs'].forEach(p => {
        p['text'] = this.punctuateBack(p['text']);
        p['sentences'].forEach(s => {
        s['text'] = this.punctuateBack(s['text']);
        s['tokens'].forEach(t => {
            t['token'] = this.punctuateBack(t['token']);
        });
        });
    });

    return parsedText;
}

function punctuateBack(text) {
    text = text.replace(/\|dot\|/g, '.');
    text = text.replace(/\|gdot\|/g, '.');
    text = text.replace(/\|gint\|/g, '?');
    text = text.replace(/\|gexc\|/g, '!');    
    text = text.replace(/\|hyp\|/g, '-');
    text = text.replace(/\|\|\|/g, '')
    return text;
}

function splitText(rawText) {

    rawText = this.preProcesText(rawText);

    var parsedText = this.parseText(rawText);    

    parsedText = this.postProcess(parsedText);

    return parsedText;
}

function tokenizeText(rawText) {
    rawText = rawText.replace(/([A-z]+)-([A-z]+)/g, "$1|hyp|$2");
    rawText = rawText.replace(/\|gdot\|/g, ".");
    rawText = rawText.replace(/\|gint\|/g, "?");
    rawText = rawText.replace(/\|gexc\|/g, "!");
    rawText = rawText.replace(/([\.\,"\(\)\[\]\{\}\?\!;:-]{1})/g, " $1 ");
    rawText = rawText.replace(/\s+/g, ' ');
    return rawText.split(' ');
}


