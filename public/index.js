var content = document.getElementById('content');
function addLine(str) {
   var newLine = document.createElement('div');
   newLine.classList.add('item');
   newLine.innerHTML = str;
   content.appendChild(newLine);
}
addLine("Welcome");
var input = document.getElementById('input');
function focusInput() {
   var range = document.createRange();
   var sel = window.getSelection();
   range.setStart(input,0);
   range.collapse(true);
   sel.removeAllRanges();
   sel.addRange(range);
}
input.addEventListener('keypress', function(e){
    if (e.keyCode==13) {
        e.preventDefault();
        e.stopPropagation();
    }
});
input.addEventListener('keyup',  function(e) {
   if(e.keyCode!=13) return;
   var str = input.innerText;
   addLine("You entered:"+str);
   input.innerHTML = '';
   focusInput();

});
setTimeout(focusInput,100);

fetch('/api/').then(function(r) {
  return r.text();
}).then(function(t) { addLine("Message from server:"+t);});
