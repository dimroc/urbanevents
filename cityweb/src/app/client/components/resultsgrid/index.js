import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router';
import { resultsgrid } from './styles';

export default class ResultsGrid extends Component {
  /*eslint-disable */
  static onEnter({store, nextState, replaceState, callback}) {
    // Load here any data.
    callback(); // this call is important, don't forget it
  }
  /*eslint-enable */

  render() {
    return <div className={resultsgrid}>
      <h1>Bunch of Results for {this.props.city.display}</h1>
    </div>;
  }
}

ResultsGrid.propTypes = {
  city: React.PropTypes.object.isRequired
};
