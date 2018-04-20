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



ifr = `<iframe width="854" height="480" 
       src="{0}?autoplay=1" frameborder="0"
       allow="autoplay; encrypted-media" 
       allowfullscreen></iframe>`

/*
req = new XMLHttpRequest()
req.open("GET","http://127.0.0.1:8080/list?limit=10",true)
req.onreadystatechange = temp
req.send()
*/
















