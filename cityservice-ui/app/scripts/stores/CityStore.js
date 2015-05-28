var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('events').EventEmitter;
var AppConstants = require('../constants/AppConstants');
var assign = require('object-assign');

var CHANGE_EVENT = 'change';

var _cities = [];

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
    this.emit(CHANGE_EVENT);
  },

  addChangeListener: function(callback) {
    this.on(CHANGE_EVENT, callback);
  },

  removeChangeListener: function(callback) {
    this.removeListener(CHANGE_EVENT, callback);
  }
});

var normalizeCities = function(cities) {
  cities.forEach(function(city) {
    city.center = [
      (city.bounds[0][0] + city.bounds[1][0]) / 2.0,
      (city.bounds[0][1] + city.bounds[1][1]) / 2.0
    ].reverse();
  });

  return cities;
}

$.ajax({
  url: "http://localhost:58080/api/v1/cities",
  context: CityStore
}).done(function(data) {
  _cities = normalizeCities(data);
  this.emitChange();
}).fail(function() {
  console.log("failed to retrieve cities");
});

module.exports = CityStore;
