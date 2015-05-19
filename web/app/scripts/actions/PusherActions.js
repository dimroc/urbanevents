var AppDispatcher = require('../dispatcher/AppDispatcher.js');
var AppConstants = require('../constants/AppConstants');

var pusher = new Pusher('81be37a4f4ee0f471476');

var PusherActions = {
  listen: function(key) {
    if (this.key != key) {
      this.key = key;
      this.channel = pusher.subscribe(key);
      this.channel.bind('tweet', this.handlePush, this);
    }
  },

  stop: function() {
    this.channel.unbind('tweet', this.handlePush);
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
