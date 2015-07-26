'use strict';

var React = require('react');
var Router = require('react-router');
var { Route, Redirect, RouteHandler, Link } = Router;
var ButtonGroup = require('react-bootstrap/lib/buttonGroup');
var Button = require('react-bootstrap/lib/button');
var CityStore = require('../stores/CityStore');

var CityHeader = React.createClass({
  mixins: [ Router.State ],

  getInitialState: function() {
    var { cityId } = this.context.router.getCurrentParams();
    return {city: CityStore.get(cityId)};
  },

  updateCity: function() {
    var { cityId } = this.context.router.getCurrentParams();
    this.setState({city: CityStore.get(cityId)});
  },

  componentDidMount: function() {
    CityStore.addChangeListener(this.updateCity);
  },

  componentWillUnmount: function() {
    CityStore.removeChangeListener(this.updateCity);
  },

  render: function () {
    if(!this.state.city) {
      return (<div className="holder">Loading...</div>);
    }

    return (
      <div className="holder">
        <header className="row">
          <div className="col-xs-12">
            <ButtonGroup justified>
              <Link className="btn btn-default" to="cities">Cities</Link>
              <Link className="btn btn-default" to="map" params={{cityId: this.state.city.key}}>{this.state.city.display} Map</Link>
              <Link className="btn btn-default" to="events" params={{cityId: this.state.city.key}}>{this.state.city.display} Events</Link>
            </ButtonGroup>
          </div>
        </header>

        <RouteHandler/>
      </div>
    );
  }
});

module.exports = CityHeader;


