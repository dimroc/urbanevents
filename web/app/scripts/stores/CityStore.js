var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('events').EventEmitter;
var AppConstants = require('../constants/AppConstants');
var assign = require('object-assign');

var CHANGE_EVENT = 'change';

var _cities = [
  {key: 'nyc', display: 'New York City', bounds: [[40.462,-74.3], [40.95,-73.65]], center: [40.7737, -73.9800]}
];

// SF -122.75,36.8,-121.75,37.8

var CityStore = assign({}, EventEmitter.prototype, {
  getAll: function() {
    return _cities;
  },

  get: function(key) {
    var matches = _cities.filter(function(value) {
      return value.key == key;
    });

    return matches[0];
  },

  emitChange: function() {
    this.emit();
  },

  addChangeListener: function(callback) {
    this.on(null, callback);
  },

  removeChangeListener: function(callback) {
    this.removeListener(callback);
  }
});

$.ajax({
  url: "http://localhost:8080/api/v1/settings",
  context: CityStore
}).done(function(data) {
  _cities = data.cities;
  this.emitChange();
}).fail(function() {
  console.log("failed to retrieve cities");
});

module.exports = CityStore;
