'use strict';

var React = require('react');
var Leaflet = require('react-leaflet');
var EventStore = require('../stores/EventStore');
var CityStore = require('../stores/CityStore');

var Map = Leaflet.Map,
  TileLayer = Leaflet.TileLayer,
  Circle = Leaflet.Circle,
  Polygon = Leaflet.Polygon;

var latLongList = function(bounds) {
  // Bounds come in as long, lat, since long is X
  var twoCorners = bounds.map(function(bound) {
    return bound.slice(0).reverse();
  });

  var corners = [];
  corners.push(twoCorners[0]);
  corners.push([twoCorners[1][0], twoCorners[0][1]]);
  corners.push(twoCorners[1]);
  corners.push([twoCorners[0][0], twoCorners[1][1]]);
  return corners;
}

var MappedEvents = React.createClass({
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
      <Map key={this.city.key} center={this.city.center} zoom={10} className="real-time-map">
        <TileLayer url="http://{s}.tile.stamen.com/toner/{z}/{x}/{y}.png" attribution='Tiles by <a href="http://maps.stamen.com/toner/#12/37.7704/-122.3781">Stamen Toner</a>'/>
        <Polygon positions={latLongList(this.city.bounds)} color="blue"/>
        {
          this.state.items.map(function(geoevent) {
            var key = geoevent.id + "mapped";
            return (<Circle key={key} center={geoevent.point.slice(0).reverse()}
                radius={500} color="red" fillColor="#f03" fillOpacity={0.5}/>);
          })
        }
      </Map>
    );
  }
});

module.exports = MappedEvents;


