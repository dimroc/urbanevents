import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router';
import { topbanner } from './styles';

export default class TopBanner extends Component {
  /*eslint-disable */
  static onEnter({store, nextState, replaceState, callback}) {
    // Load here any data.
    callback(); // this call is important, don't forget it
  }
  /*eslint-enable */

  render() {
    return <div className={topbanner}>
      <h1> New York City </h1>
      <form>
        <input type="text" name="q" />
        <input type="submit" name="Search" />
      </form>
    </div>;
  }
}

