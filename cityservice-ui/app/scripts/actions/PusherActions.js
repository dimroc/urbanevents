var AppDispatcher = require('../dispatcher/AppDispatcher.js');
var AppConstants = require('../constants/AppConstants');

var pusher = new Pusher('81be37a4f4ee0f471476');
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
    if(data.geojson.type === "Point") {
      AppDispatcher.dispatch({
        actionType: AppConstants.PUSHER_TWEET,
        geoevent: data
      });
    }
  }
}

module.exports = PusherActions;
