'use strict';

var React = require('react');
var PusherStore = require('../stores/PusherStore');

var MappedEvents = React.createClass({
  drawEvent: function(data) {
    if(data.geojson.type === "Point") {
      var circle = L.circle(data.geojson.coordinates.reverse(), 500, {
        color: 'red',
        fillColor: '#f03',
        fillOpacity: 0.5
      }).addTo(this.map);
    }
  },

  handlePush: function() {
    this.drawEvent(PusherStore.last());
  },

  componentDidMount: function() {
    this.map = L.map('map')
    this.map.setView([40.7737, -73.9800], 12);
    L.tileLayer('http://{s}.tile.stamen.com/toner/{z}/{x}/{y}.png', {
      attribution: 'Tiles by <a href="http://maps.stamen.com/toner/#12/37.7704/-122.3781">Stamen Toner</a>',
    }).addTo(this.map);

    PusherStore.getAll().forEach(function(data) {
      this.drawEvent(data);
    }, this);

    PusherStore.addChangeListener(this.handlePush);
  },
  componentWillUnmount: function() {
    PusherStore.removeChangeListener(this.handlePush);
  },
  render: function() {
    return (<div id="map"></div>);
  }
});

module.exports = MappedEvents;


