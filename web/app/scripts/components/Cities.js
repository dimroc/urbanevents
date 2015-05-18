'use strict';

var React = require('react');
var Router = require('react-router');
var { Route, Redirect, RouteHandler, Link } = Router;

var Cities = React.createClass({
  getInitialState: function() {
    return {items: ['nyc']};
  },
  render: function() {
    return (
      <div className="container-fluid">
        <h2>Cities</h2>
        <ul>
          <li><Link to="map" params={{cityId: 'nyc'}}>NYC</Link></li>
        </ul>

        <RouteHandler/>
      </div>
    );
  }
});

module.exports = Cities;


