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
  },
  componentWillUnmount: function() {
    this.channel.unbind('tweet', this.handlePush, this, this);
  },
  render: function() {
    return (
      <div className="container-fluid real-time">
        <h2>Real Time Events</h2>
        <PushedItems items={this.state.items} />
      </div>
    );
  }
});

module.exports = PusherEvents;

