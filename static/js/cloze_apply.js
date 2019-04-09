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

var isTraining = true;

$("#clozeWord").keypress(function(e) {
    if(e.which == 13) {
        nextWordAction();
    }
});

$("#clozeWordTrain").keypress(function(e) {
    if(e.which == 13) {
        nextWordTrainAction();
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

function nextWordTrainAction() {
    word = $('#clozeWordTrain').val();
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

    if (!isTraining) {
        saveWord();
    }

    if (wordIdx > totWords) {
        console.log('Fim paragrafo');

        $('#clozeSentenceTrain').html(sentence);
        $('#clozeSentence').html(sentence);

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

        if (isTraining) {            
            showTrainEnd();

        } else {

            console.log(paragraphIdx, totParagraphs);
            if (paragraphIdx+1 > totParagraphs) {
                console.log("Fim teste.");
                // alert("Agradecemos imensamente sua participação. Pode fechar essa página. Tudo foi gravado corretamente.")
                $("#btnNextWord").hide();
                $("#clozeWord").hide();                
                $("#msgEndTest").show();
            } else {
                // nextWord();
                $("#btnNextWord").hide();
                $("#clozeWord").hide();                
                $("#msgEndParagraph").show();
                $("#btnNextPar").show();
            }
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
    if (isTraining) {
        $('#clozeSentenceTrain').html(sentence);

        $('#statusTestTrain').html("Treinamento - Palavra " + wordIdx + " de " + totWords + ".");

        $('#clozeWordTrain').val('');
        $('#clozeWordTrain').focus();

    } else {
        $('#clozeSentence').html(sentence);

        $('#statusTest').html("Parágrafo " + (paragraphIdx+1) + " de " + totParagraphs + " - Palavra " + wordIdx + " de " + totWords + ".");

        $('#clozeWord').val('');
        $('#clozeWord').focus();
    }
}

function nextParagraph() {
    $("#btnNextWord").show();
    $("#clozeWord").show();   
    $("#msgEndParagraph").hide();
    $("#btnNextPar").hide();
    nextWord();
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
    $("#inicio").hide();
    $("#form").hide();
    $("#train").hide();
    $("#apply").hide();
    $("#" + this.stage).show();
}


function showForm() {
        stage = "form";
        refresh();
}

function showTrainEnd() {
    $("#btnNextWordTrain").hide();
    $("#clozeWordTrain").hide();   
    $("#msgEndTrain").show();
    $("#btnShowApply").show();
}


function showApply() {
    stage = "apply";
    refresh();
    isTraining = false;
    nextWord();
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
            birth: $('participantBirthdate').val(),
            email: $('participantEmail').val(),
            phone: $('participantPhone').val(),
            rg: $('participantRG').val(),
            cpf: $('participantCPF').val()
        }),
        contentType: "application/json"
    }).done(function(data) {
        clozeData = JSON.parse(data);
        console.log(clozeData);

        stage = "train";
        refresh();

        totParagraphs = clozeData.prgphs.length - 1;
        clozeData.prgphs[paragraphIdx].sentences.forEach(s => {
            totWords += s.qtw;
        })

        nextWord();

    }).fail(function(error) {
        alert( "Erro" );
    });


}