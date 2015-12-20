import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router';
import { resultsgrid } from './styles';
import { connect } from 'react-redux';

export class ResultsGrid extends Component {
  /*eslint-disable */
  static onEnter({store, nextState, replaceState, callback}) {
    // Load here any data.
    callback(); // this call is important, don't forget it
  }
  /*eslint-enable */

  render() {
    return <div className={resultsgrid}>
      <h1>Bunch of Results for {this.props.city.display}</h1>
      <ul>
        {this.props.geoevents.map(function(geoevent) {
          return <li key={geoevent.id}>
            <div className="full-name">{geoevent.fullName}</div>
            <div className="text">{geoevent.text}</div>
            <div className="media-url"><a href={geoevent.mediaUrl} target="_blank">{geoevent.mediaUrl}</a></div>
            <div className="neighborhoods">{geoevent.neighborhoods}</div>
            <div className="service">{geoevent.service}</div>
            <div className="created-at">{geoevent.createdAt}</div>
          </li>
        })}
      </ul>
    </div>;
  }
}

ResultsGrid.propTypes = {
  city: React.PropTypes.object.isRequired
};

function select(state) {
  return {
    geoevents: state.geoevents
  }
}

export default connect(select)(ResultsGrid)
