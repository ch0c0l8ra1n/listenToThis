String.prototype.filter = Array.prototype.filter

String.prototype.format = function() {
  /*
  // Python-esque string formatting.
  // "{0}{1}".format("foo","bar")
  // >> "foobar"
  */
  var args = arguments;
  this.unkeyed_index = 0;
  return this.replace(/\{(\w*)\}/g, function(match, key) { 
    if (key === '') {
      key = this.unkeyed_index;
      this.unkeyed_index++
    }
    if (key == +key) {
      return args[key] !== 'undefined'
      ? args[key]
      : match;
    } else {
      for (var i = 0; i < args.length; i++) {
        if (typeof args[i] === 'object' && typeof args[i][key] !== 'undefined') {
          return args[i][key];
        }
      }
      return match;
    }
  }.bind(this));
};

function parseURL(url) {
  /*
  // Parses the given url into an
  // appropriate object
  */
  var parser = document.createElement('a'),
      searchObject = {},
      queries, split, i;
  // Let the browser do the work
  parser.href = url;
  // Convert query string to object
  queries = parser.search.replace(/^\?/, '').split('&');
  for( i = 0; i < queries.length; i++ ) {
      split = queries[i].split('=');
      searchObject[split[0]] = split[1];
  }
  return {
      protocol: parser.protocol,
      host: parser.host,
      hostname: parser.hostname,
      port: parser.port,
      pathname: parser.pathname,
      search: parser.search,
      searchObject: searchObject,
      hash: parser.hash
  };
}

function embed(url){
  /*
  // Morphs the viden url into a embeddable 
  // iframe link
  */
  var u = parseURL(url)
  var pageId = u.searchObject.v
  return "https://www.youtube.com/embed/" + pageId
}

function getId(url){
  var u = parseURL(url)
  switch (u.hostname){
    case "www.youtube.com":
        return(u.searchObject.v)
        break;

    case "youtu.be":
        return(u.pathname.slice(1))
        break;
    default:
        return null
  }
}

function shuffle(array) {
  var currentIndex = array.length, temporaryValue, randomIndex;

  // While there remain elements to shuffle...
  while (0 !== currentIndex) {

    // Pick a remaining element...
    randomIndex = Math.floor(Math.random() * currentIndex);
    currentIndex -= 1;

    // And swap it with the current element.
    temporaryValue = array[currentIndex];
    array[currentIndex] = array[randomIndex];
    array[randomIndex] = temporaryValue;
  }

  return array;
}






