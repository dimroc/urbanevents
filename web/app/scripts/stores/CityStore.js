var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('events').EventEmitter;
var AppConstants = require('../constants/AppConstants');
var assign = require('object-assign');

var CHANGE_EVENT = 'change';

var _cities = [
  {key: 'nyc', display: 'New York City', center: [40.7737, -73.9800]}
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
  }
});

module.exports = CityStore;
