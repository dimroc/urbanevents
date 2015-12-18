import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { getCitiesAsync } from '#app/actions';

export default class App extends Component {
  /*eslint-disable */
  static onEnter({store, nextState, replaceState, callback}) {
    // Load here any data.
    callback();
  }
  /*eslint-enable */

  render() {
    return <div>
      <Helmet title='New Tweet City Media Search' />
      {this.props.children}
    </div>;
  }
}
