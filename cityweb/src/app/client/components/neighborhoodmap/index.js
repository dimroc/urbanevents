import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router';
import { neighborhoodmap } from './styles';

export default class NeighborhoodMap extends Component {
  /*eslint-disable */
  static onEnter({store, nextState, replaceState, callback}) {
    // Load here any data.
    callback();
  }
  /*eslint-enable */

  render() {
    return <div className={neighborhoodmap}>
      <h1>Map for {this.props.city.display}</h1>
    </div>;
  }
}

NeighborhoodMap.propTypes = {
  city: React.PropTypes.object.isRequired
};
