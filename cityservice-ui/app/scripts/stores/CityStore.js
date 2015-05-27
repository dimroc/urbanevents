var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('events').EventEmitter;
var AppConstants = require('../constants/AppConstants');
var assign = require('object-assign');

var CHANGE_EVENT = 'change';

var _cities = [
  {key: 'nyc', display: 'New York City', bounds: [[40.462,-74.3], [40.95,-73.65]], center: [40.7737, -73.9800]}
];

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
  url: "http://localhost:8080/api/v1/settings",
  context: CityStore
}).done(function(data) {
  _cities = normalizeCities(data.cities);
  this.emitChange();
}).fail(function() {
  console.log("failed to retrieve cities");
});

module.exports = CityStore;
