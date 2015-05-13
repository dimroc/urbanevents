'use strict';

var React = require('react');
var PusherStore = require('../stores/PusherStore');

var PushedItems = React.createClass({
  render: function() {
    var createItem = function(geoevent) {
      return (
        <li key={geoevent.id}>
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
  handlePush: function() {
    this.setState({items: PusherStore.getAll()});
  },
  componentDidMount: function() {
    PusherStore.addChangeListener(this.handlePush);
    this.setState({items: PusherStore.getAll()});
  },
  componentWillUnmount: function() {
    PusherStore.removeChangeListener(this.handlePush);
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

