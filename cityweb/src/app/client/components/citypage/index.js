import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router';
import TopBanner from '#app/components/topbanner';
import NeighborhoodMap from '#app/components/neighborhoodmap';
import ResultsGrid from '#app/components/resultsgrid';
import { connect } from 'react-redux';
import { setCurrentCity } from '#app/actions';

export class Citypage extends Component {
  /*eslint-disable */
  static onEnter({store, nextState, replaceState, callback}) {
    // Load here any data.
    callback(); // this call is important, don't forget it
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
      <h1>{this.props.city.display}</h1>
    </div>;
  }
}
      /*
      <TopBanner city={this.props.city}/>
      <NeighborhoodMap city={this.props.city}/>
      <ResultsGrid city={this.props.city}/>
      */

function select(state) {
  return {
    city: state.cities.current
  }
}

export default connect(select)(Citypage)
