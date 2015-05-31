'use strict';

var React = require('react');
var Router = require('react-router');
var CityStore = require('../stores/CityStore');
var { Route, Redirect, RouteHandler, Link } = Router;
var BarChart = require("react-chartjs").Bar;
var AppConstants = require('../constants/AppConstants');
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
    return {items: CityStore.getAll(), url: AppConstants.CITYSERVICE_URL};
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

  handleChange: function() {
    this.setState({url: event.target.value});
  },

  render: function() {
    var chartOptions = {maintainAspectRatio: false, responsive: true}
    return (
      <div className="container-fluid">
        <div className="row">
          <h1 className="col-sm-3">Connection</h1>
          <div className="col-sm-9 connection">
            <form method="get">
              <div className="input-group">
                <input type="text" className="form-control string optional" name="url" value={this.state.url} onChange={this.handleChange}/>
                <span className="input-group-btn">
                  <input className="btn btn-default" type="submit" value="Go!"/>
                </span>
              </div>
            </form>
          </div>
        </div>

        <div className="row">
          <h2 className="col-sm-3">Cities</h2>
          <h2 className="col-sm-9 hidden-xs">Tweets</h2>
        </div>
        {
          this.state.items.map(function(city) {
            return (
              <div className="row" key={city.key}>
                <Link className="city-key col-sm-3" to="map" params={{cityId: city.key}}>{city.display}</Link>
                <div className="chart-container col-sm-9">
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

