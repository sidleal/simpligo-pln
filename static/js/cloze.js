var editor = document.getElementById("editor")
var range = document.createRange();
var sel = window.getSelection();
range.setStart(editor.childNodes[0], editor.innerText.length);
range.collapse(true);
sel.removeAllRanges();
sel.addRange(range);

editor.focus();

var lastWord = '';

function treatKeyUp(ev) {
//	console.log(ev.key)
//	console.log(ev.keyCode)

    if (ev.keyCode == 32) {
        console.log(lastWord);
        lastWord = '';
    } else {
        lastWord = lastWord + ev.key
    }

}
