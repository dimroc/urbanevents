'use strict';

var React = require('react');
var PusherStore = require('../stores/PusherStore');
var CityStore = require('../stores/CityStore');
var PusherActions = require("../actions/PusherActions");

var PushedItems = React.createClass({
  render: function() {
    var createItem = function(geoevent) {
      var key = geoevent.id + "listed";
      return (
        <li key={key}>
          <span>{geoevent.payload}</span>
          <div className="location-info text text-muted">
            <span className="type">{geoevent.geojson.type}</span>
            <span className="coordinates">{geoevent.geojson.coordinates}</span>
            <span className="created-at">{geoevent.createdAt}</span>
          </div>
          <hr/>
        </li>
      );
    };
    return <ul>{this.props.items.map(createItem)}</ul>;
  }
});

var PusherEvents = React.createClass({
  contextTypes: {
    router: React.PropTypes.func
  },
  getInitialState: function() {
    var { cityId } = this.context.router.getCurrentParams();
    this.city = CityStore.get(cityId);

    PusherActions.listen(this.city.key);
    return {items: PusherStore.getAll(this.city.key)};
  },
  handlePush: function() {
    this.setState({items: PusherStore.getAll(this.city.key)});
  },
  componentDidMount: function() {
    PusherStore.addChangeListener(this.city.key, this.handlePush);
  },
  componentWillUnmount: function() {
    PusherStore.removeChangeListener(this.city.key, this.handlePush);
  },
  render: function() {
    return (
      <div className="container-fluid below">
        <h2>{this.city.display} Events</h2>
        <PushedItems items={this.state.items} />
      </div>
    );
  }
});

module.exports = PusherEvents;

