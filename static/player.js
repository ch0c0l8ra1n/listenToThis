var tag = document.createElement('script');
tag.src = "https://www.youtube.com/iframe_api";
var firstScriptTag = document.getElementsByTagName('script')[0];
firstScriptTag.parentNode.insertBefore(tag, firstScriptTag);

var xhr = new XMLHttpRequest()
xhr.open("GET","http://127.0.0.1:8080/youtubevideos?sub=listentothis",true)
xhr.onload = onPlaylistReceived

//onYouTubeIframeAPIReady()

var player;
function onYouTubeIframeAPIReady() {
  console.log("Player loaded")
  player = new YT.Player('player', {
    height: '390',
    width: '640',
    playerVars : {
      'autoplay' : 1,
      'modestbranding' : 1,
      'enablejsapi' : 1
    },
    events: {
      'onReady': onPlayerReady,
      'onStateChange': onPlayerStateChange,
      }
    });
}


function onPlaylistReceived(){
  var playlistJsonRaw = JSON.parse(this.responseText)
  var playlistJson = playlistJsonRaw.Posts
  shuffle(playlistJson)
  playlistJson = playlistJson.map(x => getId(x["Url"])).filter(x => x)
  console.log(playlistJson)
  player.loadPlaylist(playlistJson)
}


function onPlayerReady(event) {
  xhr.send()
}


function onPlayerStateChange(event) {
  if (event.data == YT.PlayerState.UNSTARTED) {
    player.playVideo()
  }
}

function onCueChange() {
  player.playVideo();
}


function searchKeyPressed(event){
  if (event.keyCode == 13){
    switchPlaylist()
  }
}

function switchPlaylist(){
  searchField = document.getElementsByClassName("searchWrapper")[0].children[0]
  searchTerm = searchField.value
  searchTerm = searchTerm.filter(x => x.match(/^[0-9a-zA-Z]+$/)  ).join("")
  searchField.value = searchTerm

  console.log(searchTerm)

  var xhr = new XMLHttpRequest()
  xhr.open("GET","http://127.0.0.1:8080/youtubevideos?sub=" + searchTerm,true)
  xhr.onload = onPlaylistReceived
  xhr.send()

}






