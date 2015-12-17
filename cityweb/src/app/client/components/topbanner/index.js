import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router';
import { topbanner } from './styles';
import { createHistory } from 'history';
import urlParameters from '#app/utils/urlParameters';
import * as actions from '#app/actions';

export default class TopBanner extends Component {
  /*eslint-disable */
  static onEnter({store, nextState, replaceState, callback}) {
    // Load here any data.
    callback(); // this call is important, don't forget it
  }
  /*eslint-enable */

  constructor(props) {
    super(props);
    this.state = {q: urlParameters('q')};
  }

  handleQueryChange(e) {
    this.setState({q: e.target.value});
  }

  handleSubmit(e) {
    e.preventDefault();
    var q = this.state.q.trim();
    if(!q) { return; }

    // Send request to server
    $.ajax({
      url: "/api/v1/cities/nyc/search",
      data: {q: q},
      dataType: 'json',
      cache: false,
      success: function(data) {
        console.log(data);
      }.bind(this),
      error: function(xhr, status, err) {
        console.warn("Query for " + q + " failed", err);
      }.bind(this),
      complete: function() {
        var history = createHistory();
        history.push({
          pathname: '/',
          search: '?q='+q,
          state: { q: q }
        })
      }.bind(this)
    });
  }

  render() {
    let name = this.props.name;
    return <div className={topbanner}>
      <h1> {name} </h1>
      <form onSubmit={this.handleSubmit.bind(this)}>
        <input type="text" name="q" placeholder="Enter your search query"
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
  name: React.PropTypes.string.isRequired
};

TopBanner.defaultProps = { name: 'nyc' };
