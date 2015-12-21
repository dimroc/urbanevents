import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router';
import TopBanner from '#app/components/topbanner';
import NeighborhoodMap from '#app/components/neighborhoodmap';
import ResultsGrid from '#app/components/resultsgrid';
import { setCurrentCity } from '#app/actions';
import { connect } from 'react-redux';
import { pushPath } from 'redux-simple-router';
import { Button } from 'react-bootstrap';
import { getCitiesAsync, clearGeoevents } from '#app/actions';
import { cities, citytile, citytileGrid } from './styles';

export class Homepage extends Component {
  /*eslint-disable */
  static onEnter({store, nextState, replaceState, callback}) {
    // Load here any data.
    callback();
  }
  /*eslint-enable */

  componentDidMount() {
    store.dispatch(getCitiesAsync())
    store.dispatch(clearGeoevents())
  }

  /* Change this landing page to a list of cities?
   * Show a few tiles showing the hearts of the city perhaps as
   * a jpg or a leaflet map?
   */
  render() {
    return <div className={cities}>
      <Helmet
        title='New Tweet City'
        meta={[
          {
            property: 'og:title',
            content: 'New Tweet City Media Search'
          }
        ]}
      />
      <div className={citytileGrid}>
        {this.props.cities.map(function(city) {
          return <div key={city.key} className={citytile + " uk-panel-box"}>
            <div>
              <Link to={city.key} >{city.display}</Link>
            </div>
          </div>
        })}
      </div>
    </div>;
  }
}

function select(state) {
  return {
    cities: state.cities.cities
  }
}

export default connect(select)(Homepage);
