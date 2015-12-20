import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router';
import { resultsgrid } from './styles';
import { connect } from 'react-redux';
import Geoevent from '#app/components/geoevent';

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
      <div>
        {this.props.geoevents.map((geoevent) => {
          return <Geoevent geoevent={geoevent}/>
        })}
      </div>
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
