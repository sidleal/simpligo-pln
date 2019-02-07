var paragraph;
var sentence = "";
var lastToken;
var paragraphIdx = 0;
var sentenceIdx = 0;
var tokenIdx = 0;
var wordIdx = 1;
var totWords = 0;

totParagraphs = clozeTest.paragraphs.length;
clozeTest.paragraphs[paragraphIdx].sentences.forEach(s => {
    totWords += s.qtw;
})

nextWord();

$("#clozeWord").keypress(function(e) {
    if(e.which == 13) {
        nextWord();
    }
});


function nextWord() {
    console.log(wordIdx, sentenceIdx, tokenIdx);
    console.log($('#clozeWord').val());

    totSentences = clozeTest.paragraphs[paragraphIdx].sentences.length;
    totTokens = clozeTest.paragraphs[paragraphIdx].sentences[sentenceIdx].tokens.length;

    lastToken = clozeTest.paragraphs[paragraphIdx].sentences[sentenceIdx].tokens[tokenIdx];
    if (lastToken.w == 1) {
        wordIdx++;
        if ('.,)]}!?:'.indexOf(lastToken.token) < 0 && '([{'.indexOf(sentence.substr(-1)) < 0) {
            sentence += ' ' + lastToken.token;
        } else {
            sentence += lastToken.token;
        }
    }

    tokenIdx++;

    lastToken = clozeTest.paragraphs[paragraphIdx].sentences[sentenceIdx].tokens[tokenIdx];
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
        lastToken = clozeTest.paragraphs[paragraphIdx].sentences[sentenceIdx].tokens[tokenIdx];
    }

    if (tokenIdx >= totTokens) {
        tokenIdx = 0;
        sentenceIdx++;
        //primeira palavra da nova sentença
        lastToken = clozeTest.paragraphs[paragraphIdx].sentences[sentenceIdx].tokens[tokenIdx];
        sentence += ' ' + lastToken.token;
        tokenIdx++;
        wordIdx++;
    }

    if (sentenceIdx >= totSentences) {
        console.log('Fim paragrafo');
    }

    $('#clozeSentence').html(sentence);

    $('#statusTest').html("Teste 1 de 5 - Palavra " + wordIdx + " de " + totWords + ".");

    $('#clozeWord').val('');
    $('#clozeWord').focus();


}

