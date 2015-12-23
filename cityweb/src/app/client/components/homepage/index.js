import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { pushPath } from 'redux-simple-router';
import { getCitiesAsync, getAcrossAsync } from '#app/actions';
import { citytile, citytileGrid, searchBar } from './styles';
import Geoevent from '#app/components/geoevent';

class Homepage extends Component {
  constructor(props) {
    super(props);

    let { query } = this.props.location
    let q = query && query.q
    this.state = { q: q }
  }

  componentDidMount() {
    const { getCitiesAsync, getAcrossAsync } = this.props;
    getCitiesAsync();

    let { query } = this.props.location
    let q = query && query.q
    if (q) {
      getAcrossAsync(q.trim());
    }
  }

  handleQueryChange(e) {
    this.setState({q: e.target.value});
  }

  handleSearch(q) {
    if(q) {
      const { history } = this.props;
      const transitionTo = history.pushState.bind(history, null);
      const { getAcrossAsync } = this.props;
      getAcrossAsync(q.trim(), transitionTo);
    }
  }

  handleSubmit(e) {
    e.preventDefault();
    this.handleSearch(this.state.q);
    this.refs.q.blur();
  }

  /* Change this landing page to a list of cities?
   * Show a few tiles showing the hearts of the city perhaps as
   * a jpg or a leaflet map?
   */
  render() {
    let { cities, across } = this.props;
    return <div>
      <Helmet
        title='New Tweet City'
        meta={[
          {
            property: 'og:title',
            content: 'New Tweet City Media Search'
          }
        ]}
      />

      <form onSubmit={this.handleSubmit.bind(this)} className={searchBar + " uk-form"}>
        <input type="search" name="q" ref="q" placeholder="Enter a word"
          tabIndex="1"
          value={this.state.q}
          onChange={this.handleQueryChange.bind(this)}
        />
        <input className="uk-button" type="submit" tabIndex="2"/>
      </form>

      <div className={citytileGrid}>
        {cities.map(function(city) {
          return <div key={city.key} className={citytile}>
            <h1>{city.display}</h1>
            <div className="uk-flex uk-flex-column uk-flex-middle uk-flex-nowrap">
              {(across[city.key] || []).map(function(geoevent) {
                return <Geoevent geoevent={geoevent} key={geoevent.id}/>
              })}
            </div>
          </div>
        })}
      </div>
    </div>;
  }
}

function mapStateToProps(state) {
  return {
    cities: state.cities.cities,
    across: state.cities.across
  }
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({
    getCitiesAsync,
    getAcrossAsync
  }, dispatch);
}

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Homepage);
