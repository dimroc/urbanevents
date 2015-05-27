'use strict';

var React = require('react');
var Router = require('react-router');
var CityStore = require('../stores/CityStore');
var { Route, Redirect, RouteHandler, Link } = Router;

var Cities = React.createClass({
  getInitialState: function() {
    return {items: CityStore.getAll()};
  },

  updateCities: function() {
    this.setState({items: CityStore.getAll()});
  },

  componentDidMount: function() {
    CityStore.addChangeListener(this.updateCities);
  },

  componentWillUnmount: function() {
    CityStore.removeChangeListener(this.updateCities);
  },

  render: function() {
    return (
      <div className="container-fluid">
        <h2>Cities</h2>
        <ul>
          {
            this.state.items.map(function(city) {
              return (<li key={city.key}><Link to="map" params={{cityId: city.key}}>{city.display}</Link></li>)
            })
          }
        </ul>

        <RouteHandler/>
      </div>
    );
  }
});

module.exports = Cities;

