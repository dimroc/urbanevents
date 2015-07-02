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
  // Bounds come in as long, lat, since long is X and this complies with GeoJSON
  // Leaflet expects lat, long for historical reasons, so we have to reverse.
  var reversed = bounds.slice(0).reverse();
  var corners = [];
  corners.push([reversed[0], reversed[1]]);
  corners.push([reversed[2], reversed[1]]);
  corners.push([reversed[2], reversed[3]]);
  corners.push([reversed[0], reversed[3]]);
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
    var cityKey = this.city.key;

    return (
      <Map key={cityKey} center={this.city.center} zoom={10} className="real-time-map">
        <TileLayer url="http://{s}.tile.stamen.com/toner/{z}/{x}/{y}.png" attribution='Tiles by <a href="http://maps.stamen.com/toner/#12/37.7704/-122.3781">Stamen Toner</a>'/>
        <Polygon positions={latLongList(this.city.bbox)} color="blue"/>
        {
          this.city.circles.map(function(circle, index) {
            var key = cityKey + "circle" + index;
            return (<Circle key={key} center={circle.point.slice(0).reverse()}
                radius={circle.radius * 1000} color="gray" fillColor="#aaa" fillOpacity={0.5}/>);
          })
        }
        {
          this.state.items.map(function(geoevent) {
            var key = geoevent.id + "mapped";
            return (<Circle key={key} center={geoevent.point.slice(0).reverse()}
                radius={50} color="red" fillColor="#f03" fillOpacity={0.5}/>);
          })
        }
      </Map>
    );
  }
});

module.exports = MappedEvents;


