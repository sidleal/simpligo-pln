var paragraph;
var sentence = "";
var lastToken;
var paragraphIdx = 0;
var sentenceIdx = 0;
var tokenIdx = 0;
var wordIdx = 1;
var totWords = 0;
var totParagraphs = 0;

var clozeData;
var stage;

$("#clozeWord").keypress(function(e) {
    if(e.which == 13) {
        nextWord();
    }
});


function nextWord() {
    console.log(wordIdx, sentenceIdx, tokenIdx);
    console.log($('#clozeWord').val());

    totSentences = clozeData.prgphs[paragraphIdx].sentences.length;
    totTokens = clozeData.prgphs[paragraphIdx].sentences[sentenceIdx].tokens.length;

    lastToken = clozeData.prgphs[paragraphIdx].sentences[sentenceIdx].tokens[tokenIdx];
    if (lastToken.w == 1) {
        wordIdx++;
        if ('.,)]}!?:'.indexOf(lastToken.token) < 0 && '([{'.indexOf(sentence.substr(-1)) < 0) {
            sentence += ' ' + lastToken.token;
        } else {
            sentence += lastToken.token;
        }
    }

    tokenIdx++;

    lastToken = clozeData.prgphs[paragraphIdx].sentences[sentenceIdx].tokens[tokenIdx];
    while (lastToken.w == 0) {
        if ('.,)]}!?:'.indexOf(lastToken.token) < 0 && '([{'.indexOf(sentence.substr(-1)) < 0) {
            sentence += ' ' + lastToken.token;
        } else {
            sentence += lastToken.token;
        }
        tokenIdx++;
        if (tokenIdx >= totTokens) {
            break;
        }
        lastToken = clozeData.prgphs[paragraphIdx].sentences[sentenceIdx].tokens[tokenIdx];
    }

    if (tokenIdx >= totTokens) {
        tokenIdx = 0;
        sentenceIdx++;
        //primeira palavra da nova sentença
        lastToken = clozeData.prgphs[paragraphIdx].sentences[sentenceIdx].tokens[tokenIdx];
        sentence += ' ' + lastToken.token;
        tokenIdx++;
        wordIdx++;
    }

    if (sentenceIdx >= totSentences) {
        console.log('Fim paragrafo');
    }

    $('#clozeSentence').html(sentence);

    $('#statusTest').html("Parágrafo 1 de " + totParagraphs + " - Palavra " + wordIdx + " de " + totWords + ".");

    $('#clozeWord').val('');
    $('#clozeWord').focus();


}

function refresh() {
    $("#newTest").hide();
    $("#apply").hide();
    $("#" + this.stage).show();
}


function continueTest(clozeCode) {

    $.ajax({
        type: 'POST',
        url: '/cloze/apply/new',
        data: JSON.stringify({
            code: clozeCode,
            name: $('#participantName').val(),
            org: $('#participantOrg').val(),
            ra: $('#participantRA').val(),
            sem: $('#participantSem').val(),
        }),
        contentType: "application/json"
    }).done(function(data) {
        clozeData = JSON.parse(data);
        console.log(clozeData);

        stage = "apply";
        refresh();

        totParagraphs = clozeData.prgphs.length;
        clozeData.prgphs[paragraphIdx].sentences.forEach(s => {
            totWords += s.qtw;
        })

        nextWord();

    }).fail(function(error) {
        alert( "Erro" );
    });


}