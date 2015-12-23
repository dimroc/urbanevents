import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { connect } from 'react-redux';
import { pushPath } from 'redux-simple-router';
import { getCitiesAsync, clearAcross, getAcrossAsync } from '#app/actions';
import { citytile, citytileGrid, searchBar } from './styles';
import Geoevent from '#app/components/geoevent';

export class Homepage extends Component {
  /*eslint-disable */
  static onEnter({store, nextState, replaceState, callback}) {
    // Load here any data.
    callback();
  }
  /*eslint-enable */

  constructor(props) {
    super(props);
    let { query } = this.props.location
    let q = query && query.q
    this.state = { q: q }
  }

  componentWillMount() {
    console.log("will mount")
  }

  componentDidMount() {
    console.log("did mount")
    store.dispatch(getCitiesAsync())
    let { query } = this.props.location
    let q = query && query.q
    if (q) {
      store.dispatch(getAcrossAsync(q.trim()));
    }
  }

  componentWillReceiveProps(nextProps) {
    console.log("will receive props")
  }

  shouldComponentUpdate() {
    console.log("should component update", arguments)
    return true
  }

  componentWillUpdate() {
    console.log("component will update", arguments)
  }

  componentDidUpdate() {
    console.log("component did update", arguments)
  }

  handleQueryChange(e) {
    this.setState({q: e.target.value});
  }

  handleSubmit(e) {
    e.preventDefault();
    var q = this.state.q;
    if(q) {
      store.dispatch(pushPath('?q='+q));
      store.dispatch(getAcrossAsync(q.trim()));
    }

    this.refs.q.blur();
  }

  /* Change this landing page to a list of cities?
   * Show a few tiles showing the hearts of the city perhaps as
   * a jpg or a leaflet map?
   */
  render() {
    console.log("rendering");
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
        {this.props.cities.map(function(city) {
          return <div key={city.key} className={citytile}>
            <h1>{city.display}</h1>
            <div className="uk-flex uk-flex-column uk-flex-middle uk-flex-nowrap">
              {(city.geoevents || []).map(function(geoevent) {
                return <Geoevent geoevent={geoevent} key={geoevent.id}/>
              })}
            </div>
          </div>
        })}
      </div>
    </div>;
  }
}

function select(state) {
  return {
    cities: state.cities.cities
  }
}

export default connect(select)(Homepage);
