var keyMirror = require('keymirror');

var urlParam = function(name) {
  var results = new RegExp('[\?&]' + name + '=([^&#]*)').exec(window.location.href);
  if (results==null){
    return null;
  }
  else{
    return results[1] || 0;
  }
};

var getCityserviceUrl = function() {
  var param = urlParam('url');
  if(param) {
    return decodeURIComponent(param);
  } else {
    return __ENV_CITYSERVICE_URL__;
  }
};

module.exports = $.extend(keyMirror({
  PUSHER_CHANGE_CHANNEL: null,
  PUSHER_RESET_STORE: null,
  PUSHER_TWEET: null
}), {
  CITYSERVICE_URL: getCityserviceUrl(),
  PUSHER_KEY: __ENV_PUSHER_KEY__
});
