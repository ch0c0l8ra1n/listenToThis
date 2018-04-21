function temp(){
  if (this.readyState != 4 || this.status != 200) {
    return      
  }
  r = JSON.parse(req.responseText)

  var url = r[0]["Url"]

  var embeddedUrl = embed(url)

  console.log(embeddedUrl)


  holder = document.getElementById("player")

  holder.innerHTML = ifr.format(embeddedUrl)
}


// Download the playerapi library first.
var tag = document.createElement('script');
tag.src = "https://www.youtube.com/iframe_api";
var firstScriptTag = document.getElementsByTagName('script')[0];
firstScriptTag.parentNode.insertBefore(tag, firstScriptTag);

// Now we create a player instance
var player;

// Once the library has finished loading,
// this will be executed.
function onYouTubeIframeAPIReady() {
  player = new YT.Player('player', {
    height: '390',
    width: '640',
    playerVars : {'autoplay' : 1},
    events: {
      'onReady': onPlayerReady,
      'onStateChange': onPlayerStateChange
    }
  });
}

function playlistRetreived(){
  if (this.readyState == 4 && this.status == 200) {
      playlist = JSON.parse(this.responseText).Posts;
      playlist = playlist.map(x => getId(x["Url"]) )
      player.loadPlaylist(shuffle(playlist))
    }
}

function onPlayerReady(event) {
  var xhr = new XMLHttpRequest();
  xhr.open("GET","http://127.0.0.1:8080/youtubevideos?sub=listentothis",true);
  xhr.onreadystatechange= playlistRetreived;
  xhr.send();
}

var done = false;
function onPlayerStateChange(event) {
  if (event.data != YT.PlayerState.PLAYING && !done) {
    player.playVideo()
    done = true;
  }
}

function stopVideo() {
  player.stopVideo();
}















