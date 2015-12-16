(function() {
  conn = new WebSocket("ws://127.0.0.1:8080/ws/music");
  conn.onmessage = function(evt) {
    var a = JSON.parse(evt.data); 
    if (a.op=='next') {
      document.getElementsByClassName('nxt')[0].click();
    }
    else if(a.op=='prev') {
      document.getElementsByClassName('prv')[0].click();
    }
    else if(a.op=='play') {
      document.getElementsByClassName('ply')[0].click();
    }
  }
  alert('success')
})();
