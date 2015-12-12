import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router';
import { example, p, link } from './styles';
import TopBanner from '#app/components/topbanner';

export default class Homepage extends Component {
  /*eslint-disable */
  static onEnter({store, nextState, replaceState, callback}) {
    // Load here any data.
    callback(); // this call is important, don't forget it
  }
  /*eslint-enable */

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
      <TopBanner />
      <h1 className={example}>
        Media Search
      </h1>
    </div>;
  }

}
