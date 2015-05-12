var React = window.React = require('react'),
    Router = require('react-router'),
    Timer = require("./ui/Timer"),
    PusherEvents = require("./components/PusherEvents"),
    ReactBootstrap = require('react-bootstrap'),
    mountNode = document.getElementById("app");

var DefaultRoute = Router.DefaultRoute;
var Link = Router.Link;
var Route = Router.Route;
var RouteHandler = Router.RouteHandler;

var ButtonGroup = ReactBootstrap.ButtonGroup;
var Button = ReactBootstrap.Button;

//http://getbootstrap.com/components/#btn-groups-justified
var App = React.createClass({
  render: function () {
    return (
      <div>
        <header className="row">
          <div className="col-xs-12">
            <ButtonGroup justified>
              <Button href="/#">Map</Button>
              <Button href="/#/events">Events</Button>
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
    <Route name="events" handler={PusherEvents}/>
  </Route>
);

Router.run(routes, function (Handler) {
  React.render(<Handler/>, mountNode);
});
