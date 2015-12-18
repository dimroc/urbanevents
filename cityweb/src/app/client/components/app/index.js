import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { getCitiesAsync } from '#app/actions';

export default class App extends Component {
  componentDidMount() {
    store.dispatch(getCitiesAsync());
    // todo: hydrate whole application based on URL.
    // So AFTER getCities succeeds, create state tree based on url.
  }

  render() {
    return <div>
      <Helmet title='Go + React + Redux = rocks!' />
      {this.props.children}
    </div>;
  }
}
