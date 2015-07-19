
function get(url) {
  return fetch(url,{credentials: 'same-origin'} );
}
function postJSON(url,data) {
  return fetch(url,{
    credentials: 'same-origin',
    method: 'post',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
  });
}
function post(url,data) {
  return fetch(url,{
    credentials: 'same-origin',
    method: 'post',
      headers: {
        'Accept': 'text/plain',
        'Content-Type': 'text/plain'
      },
      body: data
  });
}


var fetchPostOptions = {credentials: 'same-origin'};
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
   //addLine("You entered:"+str);
   post("/api/message",str);
   input.innerHTML = '';
   focusInput();

});
setTimeout(focusInput,100);
get('/api/user').then(function(r) {
  return r.text();
}).then(function(t) { addLine("Message from server:"+t);});



var channel;
function join() {
    get('/api/join').then(function(r) { return r.text(); }).then(function(token) {
        console.log("joining", token);
        channel = new goog.appengine.Channel(token);
        socket = channel.open();
        socket.onopen = onOpened;
        socket.onmessage = onMessage;
        socket.onerror = onError;
        socket.onclose = onClose;
    });
}
function onOpened() {
    addLine("Channel opened");
}
function sanitize(txt) {
   var d = document.createElement('div');
   d.innerHTML = txt;
   return d.innerText;
}
function onMessage(m) {
    var message = JSON.parse(m.data);
    message.When = new Date(message.When);
    console.log(message);
    addLine(message.When.getHours()+":"+message.When.getMinutes()+ " <span class='name'>"+message.Sender+"</span>: <span class='message'>"+message.Text+"</span>");
}
function onError() {
    addLine("Error in channel");
}
function onClose() {
    addLine("Channel closed");
}
join();
