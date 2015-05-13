var React = window.React = require('react'),
    Router = require('react-router'),
    PusherEvents = require("./components/PusherEvents"),
    MappedEvents = require("./components/MappedEvents"),
    mountNode = document.getElementById("app");

var DefaultRoute = Router.DefaultRoute;
var Route = Router.Route;
var RouteHandler = Router.RouteHandler;

var ButtonGroup = require('react-bootstrap/lib/buttonGroup');
var Button = require('react-bootstrap/lib/button');

//http://getbootstrap.com/components/#btn-groups-justified
var App = React.createClass({
  mixins: [ Router.State ],
  render: function () {
    var names = this.getRoutes().map(function(item) {
      return item.name;
    });

    var isActive = function(key) { if (names.indexOf(key) >= 0) return 'active' }

    return (
      <div className="holder">
        <header className="row">
          <div className="col-xs-12">
            <ButtonGroup justified>
              <Button href="/#/maps" className={isActive('maps') || isActive(undefined)}>Map</Button>
              <Button href="/#/events" className={isActive('events')}>Events</Button>
            </ButtonGroup>
          </div>
        </header>

        {/* this is the important part */}
        <RouteHandler/>
      </div>
    );
  }
});

var routes = (
  <Route name="app" handler={App} path="/">
    <DefaultRoute handler={MappedEvents} />
    <Route name="maps" handler={MappedEvents}/>
    <Route name="events" handler={PusherEvents}/>
  </Route>
);

Router.run(routes, function (Handler) {
  React.render(<Handler/>, mountNode);
});
