var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('events').EventEmitter;
var AppConstants = require('../constants/AppConstants');
var assign = require('object-assign');

var _geoevents = {};

function addEvent(geoevent) {
  _geoevents[geoevent.city] = _geoevents[geoevent.city] || [];
  var arr = _geoevents[geoevent.city];

  arr.unshift(geoevent);
  arr = arr.slice(0, 1000);
}

var EventStore = assign({}, EventEmitter.prototype, {
  getAll: function(key) {
    return _geoevents[key] || [];
  },

  last: function() {
    if (_geoevents[key]) {
      return _geoevents[key][0];
    } else {
      return null;
    }
  },

  emitChange: function(city) {
    this.emit(city);
  },

  addChangeListener: function(city, callback) {
    this.on(city, callback);
  },

  removeChangeListener: function(city, callback) {
    this.removeListener(city, callback);
  }
});

// Register callback to handle all updates
AppDispatcher.register(function(action) {
  switch(action.actionType) {
    case AppConstants.PUSHER_TWEET:
      addEvent(action.geoevent);
      EventStore.emitChange(action.geoevent.city);
      break;
    case AppConstants.PUSHER_RESET_STORE:
      _geoevents = {};
      break;
  }
});

module.exports = EventStore;
