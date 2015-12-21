import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { IndexLink, Link } from 'react-router';
import { topbanner } from './styles';
import { createHistory } from 'history';
import urlParameters from '#app/utils/urlParameters';
import { getGeoeventsAsync } from '#app/actions';
import { connect } from 'react-redux';

export class TopBanner extends Component {
  /*eslint-disable */
  static onEnter({store, nextState, replaceState, callback}) {
    // Load here any data.
    callback(); // this call is important, don't forget it
  }
  /*eslint-enable */

  constructor(props) {
    super(props);
    this.state = {
      q: urlParameters('q')
    };
  }

  handleQueryChange(e) {
    this.setState({q: e.target.value});
  }

  handleSubmit(e) {
    e.preventDefault();
    var q = this.state.q.trim();
    if(q) {
      store.dispatch(getGeoeventsAsync(this.props.city.key, q));
    }

    this.refs.q.blur();
  }

  render() {
    var label = this.props.city ? this.props.city.display : "Select a City";
    return <div className={topbanner}>
      <IndexLink to='/'>All Cities</IndexLink>
      <h1> {label} </h1>
      <form onSubmit={this.handleSubmit.bind(this)}>
        <input type="text" name="q" ref="q" placeholder="Enter your search query"
          tabIndex="0"
          value={this.state.q}
          onChange={this.handleQueryChange.bind(this)}
        />
        <input type="submit" tabIndex="1"/>
      </form>
    </div>;
  }
}

TopBanner.propTypes = {
  cityKey: React.PropTypes.string.isRequired
};

TopBanner.defaultProps = { cityKey: 'nyc' };

function select(state) {
  return {
    city: state.cities.current, // But cities is an array newb
  }
}

export default connect(select)(TopBanner)
