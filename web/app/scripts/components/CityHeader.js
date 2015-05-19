'use strict';

var React = require('react');
var Router = require('react-router');
var { Route, Redirect, RouteHandler, Link } = Router;
var ButtonGroup = require('react-bootstrap/lib/buttonGroup');
var Button = require('react-bootstrap/lib/button');
var CityStore = require('../stores/CityStore');

var CityHeader = React.createClass({
  mixins: [ Router.State ],
  render: function () {
    var names = this.getRoutes().map(function(item) {
      return item.name;
    });

    var isActive = function(key) { if (names.indexOf(key) >= 0) return 'active' }
    var { cityId } = this.context.router.getCurrentParams();
    this.city = CityStore.get(cityId);

    return (
      <div className="holder">
        <header className="row">
          <div className="col-xs-12">
            <ButtonGroup justified>
              <Button href={"/#/cities"}>Cities</Button>
              <Button href={"/#/cities/" + cityId + "/map"} className={isActive('map')}>{this.city.display} Map</Button>
              <Button href={"/#/cities/" + cityId + "/events"} className={isActive('events')}>Events</Button>
            </ButtonGroup>
          </div>
        </header>

        <RouteHandler/>
      </div>
    );
  }
});

module.exports = CityHeader;


