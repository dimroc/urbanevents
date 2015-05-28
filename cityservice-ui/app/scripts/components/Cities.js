'use strict';

var React = require('react');
var Router = require('react-router');
var CityStore = require('../stores/CityStore');
var { Route, Redirect, RouteHandler, Link } = Router;
var BarChart = require("react-chartjs").Bar;
var moment = require("moment");

//var data = {
    //labels: ["January", "February", "March", "April", "May", "June", "July"],
    //datasets: [ { data: [65, 59, 80, 81, 56, 55, 40]}]
//};
var chartDataFor = function(city) {
  if (city.stats) {
    return {
      labels: city.stats.days.map(function(entry) { return moment(entry).format('dddd')}),
      datasets: [
        { data: city.stats.counts }
      ]
    };
  } else {
    return {
      labels: [],
      datasets: []
    };
  }
}

var Cities = React.createClass({
  getInitialState: function() {
    return {items: CityStore.getAll()};
  },

  updateCities: function() {
    this.setState({items: CityStore.getAll()});
  },

  componentDidMount: function() {
    CityStore.addChangeListener(this.updateCities);
  },

  componentWillUnmount: function() {
    CityStore.removeChangeListener(this.updateCities);
  },

  render: function() {
    var chartOptions = {maintainAspectRatio: false, responsive: true}
    return (
      <div className="container-fluid">
        <div className="row">
          <h2 className="col-sm-4">Cities</h2>
          <h2 className="col-sm-8 hidden-xs">Tweets</h2>
        </div>
        {
          this.state.items.map(function(city) {
            return (
              <div className="row" key={city.key}>
                <Link className="city-key col-sm-4" to="map" params={{cityId: city.key}}>{city.display}</Link>
                <div className="chart-container col-sm-8">
                  <BarChart className="barchart" data={chartDataFor(city)} options={chartOptions}/>
                </div>
              </div>
            )
          })
        }

        <RouteHandler/>
      </div>
    );
  }
});

module.exports = Cities;

