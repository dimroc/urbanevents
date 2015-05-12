'use strict';

var React = require('react');

var PushedItems = React.createClass({
  render: function() {
    var createItem = function(geoevent) {
      return (
        <li>
          <span>{geoevent.payload}</span>
          <div className="location-info text text-muted">
            <span className="type">{geoevent.geojson.type}</span>
            <span className="coordinates">{geoevent.geojson.coordinates}</span>
          </div>
          <hr/>
        </li>
      );
    };
    return <ul>{this.props.items.map(createItem)}</ul>;
  }
});

var PusherEvents = React.createClass({
  getInitialState: function() {
    var pusher = new Pusher('81be37a4f4ee0f471476');
    var channel = pusher.subscribe('nyc');
    channel.bind('tweet', function(data) {
      var nextItems = [data].concat(this.state.items);
      this.setState({items: nextItems});
    }, this);

    return {items: []};
  },
  render: function() {
    return (
      <div>
        <h2>Real Time Events</h2>
        <PushedItems items={this.state.items} />
      </div>
    );
  }
});

module.exports = PusherEvents;

