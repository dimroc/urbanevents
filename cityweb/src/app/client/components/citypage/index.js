import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router';
import TopBanner from '#app/components/topbanner';
import NeighborhoodMap from '#app/components/neighborhoodmap';
import ResultsGrid from '#app/components/resultsgrid';
import Geoevent from '#app/components/geoevent';
import { connect } from 'react-redux';
import { setCurrentCity, getCitiesAsync, getGeoeventsAsync } from '#app/actions';
import urlParameters from '#app/utils/urlParameters';

export class Citypage extends Component {
  /*eslint-disable */
  static onEnter({store, nextState, replaceState, callback}) {
    // Load here any data.
    let cityKey = location.pathname.substr(1);
    let q = urlParameters('q');

    store.dispatch(getCitiesAsync()).then(() => {
      store.dispatch(setCurrentCity(cityKey));
      return store.dispatch(getGeoeventsAsync(cityKey, q));
    }).then(() => {
      callback(); // this call is important, don't forget it
    })
  }
  /*eslint-enable */

  componentDidMount() {
    let { cityKey } = this.props.params
    store.dispatch(setCurrentCity(cityKey));
  }

  render() {
    return <div>
      <Helmet
        title='New Tweet City'
        meta={[
          {
            property: 'og:title',
            content: 'New Tweet City Media Search'
          }
        ]} />
      <TopBanner city={this.props.city}/>
      <NeighborhoodMap city={this.props.city}/>
      <ResultsGrid city={this.props.city}/>
    </div>;
  }
}

function select(state) {
  return {
    city: state.cities.current
  }
}

export default connect(select)(Citypage)
