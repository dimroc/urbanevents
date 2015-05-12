'use strict';

var React = require('react');

var MappedEvents = React.createClass({
  handlePush: function(data) {
    if(data.geojson.type === "Point") {
      var circle = L.circle(data.geojson.coordinates.reverse(), 500, {
        color: 'red',
        fillColor: '#f03',
        fillOpacity: 0.5
      }).addTo(this.map);
    }
  },
  componentDidMount: function() {
    var pusher = new Pusher('81be37a4f4ee0f471476');
    this.channel = pusher.subscribe('nyc');
    this.channel.bind('tweet', this.handlePush, this);

    this.map = L.map('map')
    this.map.setView([40.7737, -73.9800], 12);
    L.tileLayer('http://{s}.tile.stamen.com/toner/{z}/{x}/{y}.png', {
      attribution: 'Tiles by <a href="http://maps.stamen.com/toner/#12/37.7704/-122.3781">Stamen Toner</a>',
    }).addTo(this.map);
  },
  componentWillUnmount: function() {
    this.channel.unbind('tweet', this.handlePush, this);
  },
  render: function() {
    return (<div id="map"></div>);
  }
});

module.exports = MappedEvents;


