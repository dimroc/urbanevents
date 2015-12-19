import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { getCitiesAsync } from '#app/actions';

export default class App extends Component {
  render() {
    return <div>
      <Helmet title='New Tweet City Media Search' />
      {this.props.children}
    </div>;
  }
}
