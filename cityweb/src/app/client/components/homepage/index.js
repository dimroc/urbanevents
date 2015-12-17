import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router';
import TopBanner from '#app/components/topbanner';
import NeighborhoodMap from '#app/components/neighborhoodmap';
import ResultsGrid from '#app/components/resultsgrid';

export default class Homepage extends Component {
  /*eslint-disable */
  static onEnter({store, nextState, replaceState, callback}) {
    // Load here any data.
    callback(); // this call is important, don't forget it
  }
  /*eslint-enable */

  /* Change this landing page to a list of cities?
   * Show a few tiles showing the hearts of the city perhaps as
   * a jpg or a leaflet map?
   */
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
      <TopBanner cityKey="nyc"/>
      <NeighborhoodMap />
      <ResultsGrid />
    </div>;
  }
}
