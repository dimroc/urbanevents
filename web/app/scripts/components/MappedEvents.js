'use strict';

var React = require('react');
var Leaflet = require('react-leaflet');
var PusherStore = require('../stores/PusherStore');

var Map = Leaflet.Map,
  TileLayer = Leaflet.TileLayer,
  Circle = Leaflet.Circle

var mapKey = 0;

var MappedEvents = React.createClass({
  contextTypes: {
    router: React.PropTypes.func
  },
  getInitialState: function() {
    return {items: PusherStore.getAll(), key: mapKey++};
  },

  handlePush: function() {
    this.setState({items: PusherStore.getAll()});
  },

  componentDidMount: function() {
    PusherStore.addChangeListener(this.handlePush);
  },
  componentWillUnmount: function() {
    PusherStore.removeChangeListener(this.handlePush);
  },
  render: function() {
    var position = [40.7737, -73.9800];
    return (
      <Map key={this.state.key} center={position} zoom={12} className="real-time-map">
        <TileLayer url="http://{s}.tile.stamen.com/toner/{z}/{x}/{y}.png" attribution='Tiles by <a href="http://maps.stamen.com/toner/#12/37.7704/-122.3781">Stamen Toner</a>'/>
        {
          this.state.items.map(function(geoevent) {
            return (<Circle key={geoevent.id} center={geoevent.geojson.coordinates.reverse()}
                radius={500} color="red" fillColor="#f03" fillOpacity={0.5}/>);
          })
        }
      </Map>
    );
  }
});

module.exports = MappedEvents;


