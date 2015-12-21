import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router';
import { resultsgrid } from './styles';
import { connect } from 'react-redux';
import Geoevent from '#app/components/geoevent';
import { getGeoeventsAsync } from '#app/actions';

export class ResultsGrid extends Component {
  render() {
    return <div className={resultsgrid}>
      <h1>Bunch of Results for {this.props.city.display}</h1>
      <div className="uk-flex uk-flex-wrap uk-flex-left">
        {this.props.geoevents.map((geoevent) => {
          return <Geoevent geoevent={geoevent} key={geoevent.id}/>
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
