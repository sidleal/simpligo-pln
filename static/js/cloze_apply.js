var paragraph;
var sentence = "";
var lastToken;
var paragraphIdx = 0;
var sentenceIdx = 0;
var tokenIdx = 0;
var wordIdx = 1;
var totWords = 0;
var totParagraphs = 0;
var targetWord = "";

var clozeData;
var stage;

var startDateTime;

$("#clozeWord").keypress(function(e) {
    if(e.which == 13) {
        nextWordAction();
    }
});

function nextWordAction() {
    word = $('#clozeWord').val();
    if (word == "") {
        alert("Preencha com a palavra que ache mais provável que venha a seguir.")
        return
    }
    nextWord();
}

function nextWord() {

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
        targetWord = lastToken.token;
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

    saveWord();

    if (wordIdx > totWords) {
        console.log('Fim paragrafo');

        paragraphIdx++;
        tokenIdx = 0;
        wordIdx = 1;
        sentenceIdx = 0;
        totWords = 0;
        sentence = "";

        clozeData.prgphs[paragraphIdx].sentences.forEach(s => {
            totWords += s.qtw;
        })

        $('#clozeWord').val('');

        if (paragraphIdx+1 > totParagraphs) {
            console.log("Fim teste.");
            alert("Agradecemos imensamente sua participação. Pode fechar essa página. Tudo foi gravado corretamente.")
        } else {
            nextWord();
        }

        return
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


    startDateTime = new Date().getTime()
    $('#clozeSentence').html(sentence);

    $('#statusTest').html("Parágrafo " + (paragraphIdx+1) + " de " + totParagraphs + " - Palavra " + wordIdx + " de " + totWords + ".");

    $('#clozeWord').val('');
    $('#clozeWord').focus();


}

function saveWord() {
    word = $('#clozeWord').val();
    var elapsed = 0;
    if (word != "") {
        elapsed = new Date().getTime() - startDateTime;

        $.ajax({
            type: 'POST',
            url: '/cloze/apply/save',
            data: JSON.stringify({
                part: clozeData.part.id,
                para: paragraphIdx+1,
                sent: sentenceIdx+1,
                wseq: wordIdx-1,
                tword: targetWord,
                word: word,
                time: elapsed,
                par_id: clozeData.prgphs[paragraphIdx].idx,
                sen_id: clozeData.prgphs[paragraphIdx].sentences[sentenceIdx].idx,
                tok_id: clozeData.prgphs[paragraphIdx].sentences[sentenceIdx].tokens[tokenIdx-1].idx,
            }),
            contentType: "application/json"
        }).done(function(data) {
            console.log("Saved.")
        }).fail(function(error) {
            alert( "Desculpe. Ocorreu um erro ao salvar, tente recomeçar o teste, se o problema persistir informe o administrador." );
        });

    }
    console.log(clozeData.part.id, paragraphIdx, sentenceIdx, wordIdx, tokenIdx, targetWord, word, elapsed);
    // console.log(clozeData.prgphs[paragraphIdx].idx, clozeData.prgphs[paragraphIdx].sentences[sentenceIdx].idx);

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