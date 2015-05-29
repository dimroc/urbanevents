var keyMirror = require('keymirror');

module.exports = $.extend(keyMirror({
  PUSHER_CHANGE_CHANNEL: null,
  PUSHER_RESET_STORE: null,
  PUSHER_TWEET: null
}), {
  CITYSERVICE_URL: __ENV_CITYSERVICE_URL__,
  PUSHER_KEY: __ENV_PUSHER_KEY__
});
