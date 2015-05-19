'use strict';

var React = require('react');
var Leaflet = require('react-leaflet');
var PusherStore = require('../stores/PusherStore');
var CityStore = require('../stores/CityStore');
var PusherActions = require("../actions/PusherActions");

var Map = Leaflet.Map,
  TileLayer = Leaflet.TileLayer,
  Circle = Leaflet.Circle

var MappedEvents = React.createClass({
  contextTypes: {
    router: React.PropTypes.func
  },
  getInitialState: function() {
    var { cityId } = this.context.router.getCurrentParams();
    this.city = CityStore.get(cityId);

    PusherActions.listen(this.city.key);
    return {items: PusherStore.getAll(this.city.key)};
  },

  handlePush: function() {
    this.setState({items: PusherStore.getAll(this.city.key)});
  },

  componentDidMount: function() {
    PusherStore.addChangeListener(this.city.key, this.handlePush);
  },
  componentWillUnmount: function() {
    PusherStore.removeChangeListener(this.city.key, this.handlePush);
  },
  render: function() {
    return (
      <Map key={this.city.key} center={this.city.center} zoom={12} className="real-time-map">
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


