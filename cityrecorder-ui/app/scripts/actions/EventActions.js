var AppDispatcher = require('../dispatcher/AppDispatcher.js');
var AppConstants = require('../constants/AppConstants');
var EventSource = require('event-source');

var url = AppConstants.CITYSERVICE_URL + "/api/v1/events";

var EventActions = {
  start: function() {
    this._es = new EventSource(url);
    this._es.addEventListener("event", this.listen.bind(this));
    this._es.addEventListener("error", this.listen.bind(this));
  },

  listen: function(event) {
    var data = JSON.parse(event.data);
    if(data.locationType === "coordinate") {
      AppDispatcher.dispatch({
        actionType: AppConstants.PUSHER_TWEET,
        geoevent: data
      });
    }
  },

  error: function(event) {
    console.warn(event);
  }
}

module.exports = EventActions;
