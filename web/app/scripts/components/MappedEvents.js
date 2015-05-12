'use strict';

var React = require('react');

var MappedEvents = React.createClass({
  getInitialState: function() {
    return {items: []};
  },
  handlePush: function(data) {
    var nextItems = [data].concat(this.state.items);
    this.setState({items: nextItems});
  },
  componentDidMount: function() {
    var pusher = new Pusher('81be37a4f4ee0f471476');
    this.channel = pusher.subscribe('nyc');
    this.channel.bind('tweet', this.handlePush, this, this);

    var map = L.map('map').setView([40.7127, -74.0059], 11);
    L.tileLayer('http://{s}.tile.stamen.com/toner/{z}/{x}/{y}.png', {
      attribution: 'Stamen toner: <a href="http://maps.stamen.com/toner/#12/37.7704/-122.3781">Link</a>',
      maxZoom: 18
    }).addTo(map);
  },
  componentWillUnmount: function() {
    this.channel.unbind('tweet', this.handlePush, this, this);
  },
  render: function() {
    return (
      <div>
        <div id="map"></div>
      </div>
    );
  }
});

module.exports = MappedEvents;


