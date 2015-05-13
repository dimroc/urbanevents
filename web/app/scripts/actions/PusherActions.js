var AppDispatcher = require('../dispatcher/AppDispatcher.js');
var AppConstants = require('../constants/AppConstants');

var PusherActions = {
  start: function() {
    var pusher = new Pusher('81be37a4f4ee0f471476');
    this.channel = pusher.subscribe('nyc');
    this.channel.bind('tweet', this.handlePush, this);
  },

  stop: function() {
    this.channel.unbind('tweet', this.handlePush);
  },

  handlePush: function(data) {
    if(data.geojson.type === "Point") {
      AppDispatcher.dispatch({
        actionType: AppConstants.PUSHER_EVENT,
        geoevent: data
      });
    }
  }
}

module.exports = PusherActions;
