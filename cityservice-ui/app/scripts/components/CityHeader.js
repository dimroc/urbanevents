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
    var names = this.getRoutes().map(function(item) {
      return item.name;
    });

    var isActive = function(key) { if (names.indexOf(key) >= 0) return 'active' }

    if(!this.state.city) {
      return (<div className="holder">Loading...</div>);
    }

    return (
      <div className="holder">
        <header className="row">
          <div className="col-xs-12">
            <ButtonGroup justified>
              <Button href={"/#/cities"}>Cities</Button>
              <Button href={"/#/cities/" + this.state.city.key + "/map"} className={isActive('map')}>{this.state.city.display} Map</Button>
              <Button href={"/#/cities/" + this.state.city.key + "/events"} className={isActive('events')}>Events</Button>
            </ButtonGroup>
          </div>
        </header>

        <RouteHandler/>
      </div>
    );
  }
});

module.exports = CityHeader;


