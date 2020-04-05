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
var clozeCode;
var stage;

var startDateTime;
var timeToFirstLetter;
var timeTotalTest = 0;
var timeTotalParagraph = 0;

var isTraining = true;

$("#clozeWord").keypress(function(e) {
    if(e.which == 13) {
        nextWordAction();
    } else if (timeToFirstLetter == 0) {
        timeToFirstLetter = new Date().getTime() - startDateTime;
    }
    
});

$("#clozeWordTrain").keypress(function(e) {
    if(e.which == 13) {
        nextWordTrainAction();
    }
});

function nextWordAction() {
    word = $('#clozeWord').val();
    if (word.trim() == "") {
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
    } else {
        sentence += lastToken.token;
    }

    tokenIdx++;

    lastToken = clozeData.prgphs[paragraphIdx].sentences[sentenceIdx].tokens[tokenIdx];
    while (lastToken != null && lastToken.w == 0) {
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

        timeTotalParagraph = 0;

        $('#clozeSentenceTrain').html(sentence);
        $('#clozeSentence').html(sentence);

        paragraphIdx++;

        if (isTraining) {            
            showTrainEnd();

        } else {

            console.log(paragraphIdx, totParagraphs);
            if (paragraphIdx > totParagraphs) {
                console.log("Fim teste.");
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

        tokenIdx = 0;
        wordIdx = 1;
        sentenceIdx = 0;
        totWords = 0;
        sentence = "";

        if (paragraphIdx <= totParagraphs) {
            clozeData.prgphs[paragraphIdx].sentences.forEach(s => {
                totWords += s.qtw;
            })
        }

        $('#clozeWord').val('');

        return
    }

    if (tokenIdx >= totTokens) {
        tokenIdx = 0;
        sentenceIdx++;

        /*
        //primeira palavra da nova sentença
        lastToken = clozeData.prgphs[paragraphIdx].sentences[sentenceIdx].tokens[tokenIdx];
        sentence += ' ' + lastToken.token;
        tokenIdx++;
        wordIdx++;

        nextToken = clozeData.prgphs[paragraphIdx].sentences[sentenceIdx].tokens[tokenIdx]; //trata sentença com 1 palavra
        if (nextToken != null && nextToken.w == 0) {
            nextWord();
        }
        */
    }

    startDateTime = new Date().getTime();
    timeToFirstLetter = 0;

    if (isTraining) {
        $('#clozeSentenceTrain').html(sentence);

        $('#statusTestTrain').html("Treinamento - Palavra " + wordIdx + " de " + totWords + ".");

        $('#clozeWordTrain').val('');
        $('#clozeWordTrain').focus();

    } else {
        $('#clozeSentence').html(sentence);

        $('#statusTest').html("Parágrafo " + (paragraphIdx) + " de " + totParagraphs + " - Palavra " + wordIdx + " de " + totWords + ".");

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
        typingTime = elapsed - timeToFirstLetter;
        timeTotalTest = timeTotalTest + elapsed;
        timeTotalParagraph = timeTotalParagraph + elapsed;

        $.ajax({
            type: 'POST',
            url: '/cloze/apply/save',
            data: JSON.stringify({
                part: clozeData.part.id,
                para: paragraphIdx,
                sent: sentenceIdx+1,
                wseq: wordIdx-1,
                tword: targetWord,
                word: word,
                time: elapsed,
                time_to_start: timeToFirstLetter,
                time_typing: typingTime,
                time_total: timeTotalTest,
                time_total_par: timeTotalParagraph,
                par_id: clozeData.prgphs[paragraphIdx].idx,
                sen_id: clozeData.prgphs[paragraphIdx].sentences[sentenceIdx].idx,
                tok_id: clozeData.prgphs[paragraphIdx].sentences[sentenceIdx].tokens[tokenIdx-1].idx,
                tot_words: totWords,
                code: clozeCode,
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


function showForm(clozeCode) {
    if ($("#name").val() == "" || $("#doc").val() == "") {
        alert( "Por favor, preencha seu nome completo e número do documento para continuar." );
        $("#name").focus();
        return
    }
    
    $("#doc").val($("#doc").val().replace(" ", ""));
    saveTCLE(clozeCode, $("#name").val(), $("#doc").val());

    stage = "form";
    refresh();
    $("#participantName").val($("#name").val());
    $("#participantRG").val($("#doc").val());
    $("#participantName").focus();
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
    
    if ($('#participantName').val() == "") {
        alert("Por favor, preencha seu nome.");
        $('#participantName').focus();
        return
    }
    if ($('#participantRG').val() == "") {
        alert("Por favor, preencha seu RG.");
        $('#participantRG').focus();
        return
    }
    if ($('#participantAge').val() == "") {
        alert("Por favor, preencha sua idade.");
        $('#participantAge').focus();
        return
    }
    if ($('#participantSem').val() == "") {
        alert("Por favor, preencha o semestre que está cursando.");
        $('#participantSem').focus();
        return
    }
    if ($('#participantEmail').val() == "") {
        alert("Por favor, preencha seu email para contato.");
        $('#participantEmail').focus();
        return
    }

    this.clozeCode = clozeCode;

    $.ajax({
        type: 'POST',
        url: '/cloze/apply/new',
        data: JSON.stringify({
            code: clozeCode,
            name: $('#participantName').val(),
            rg: $('#participantRG').val(),
            age: $('#participantAge').val(),
            gender: $('#participantGender').val(),
            course: $('#participantCourse').val(),
            sem: $('#participantSem').val(),
            lang: $('#participantLang').val(),
            org: $('#participantOrg').val(),
            email: $('#participantEmail').val(),
            phone: $('#participantPhone').val(),
            cpf: $('#participantCPF').val(),
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

function saveTCLE(code, name, doc) {
    $.ajax({
        type: 'POST',
        url: '/cloze/term/'+code+'/save',
        data: {
            name: name,
            doc: doc
        },
    }).done(function(data) {
        console.log("TCLE salvo.")
    }).fail(function(error) {
        console.log( "Erro TCLE" );
        console.log( error );
    });
}