
var React = window.React = require('react'),
    Timer = require("./ui/Timer"),
    PusherEvents = require("./components/PusherEvents"),
    mountNode = document.getElementById("app");


React.render(<PusherEvents />, mountNode);

