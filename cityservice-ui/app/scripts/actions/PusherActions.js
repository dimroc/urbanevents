var AppDispatcher = require('../dispatcher/AppDispatcher.js');
var AppConstants = require('../constants/AppConstants');

var pusher = new Pusher(AppConstants.PUSHER_KEY);
var _channels = {};

var PusherActions = {
  listen: function(key) {
    if (!_channels[key]) {
      var newChannel = pusher.subscribe(key);
      newChannel.bind('tweet', this.handlePush, this);
      _channels[key] = newChannel;
    }
  },

  handlePush: function(data) {
    if(data.locationType === "coordinate") {
      AppDispatcher.dispatch({
        actionType: AppConstants.PUSHER_TWEET,
        geoevent: data
      });
    }
  }
}

module.exports = PusherActions;
