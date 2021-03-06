'use strict';

var React = require('react');
var EventStore = require('../stores/EventStore');
var CityStore = require('../stores/CityStore');
var EventActions = require("../actions/EventActions");

var PushedItems = React.createClass({
  render: function() {
    var createItem = function(geoevent) {
      var key = geoevent.id + "listed";
      var icon;
      if (geoevent.service == "twitter") {
        icon = <i className="fa fa-twitter fa-stack-2x"></i>
      } else {
        icon = <i className="fa fa-instagram fa-stack-2x"></i>
      }

      var img;
      if (geoevent.thumbnailUrl) { img = <img src={geoevent.thumbnailUrl}></img> }

      return (
        <li key={key}>
          <a href={geoevent.link} target="_blank">
            <label className="screenName label label-info">{geoevent.username}</label>
          </a>
          <span>{geoevent.text || geoevent.text_fr}</span>
          {img}
          <div className="location-info text text-muted">
            <span className="location-type">{geoevent.locationType}</span>
            <span className="coordinates">{geoevent.point}</span>
            <span className="created-at">{geoevent.createdAt}</span>
            <span className="place">Place: {geoevent.place}</span>
            <span className="neighborhoods">Hoods: {geoevent.neighborhoods}</span>
            {icon}
          </div>
          <hr/>
        </li>
      );
    };
    return <ul>{this.props.items.map(createItem)}</ul>;
  }
});

var ListEvents = React.createClass({
  contextTypes: {
    router: React.PropTypes.func
  },
  getInitialState: function() {
    var { cityId } = this.context.router.getCurrentParams();
    this.city = CityStore.get(cityId);

    return {items: EventStore.getAll(this.city.key)};
  },
  handlePush: function() {
    this.setState({items: EventStore.getAll(this.city.key)});
  },
  componentDidMount: function() {
    EventStore.addChangeListener(this.city.key, this.handlePush);
  },
  componentWillUnmount: function() {
    EventStore.removeChangeListener(this.city.key, this.handlePush);
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

module.exports = ListEvents;

