import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { getCitiesAsync } from '#app/actions';

class App extends Component {
  render() {
    return <div>
      <Helmet title='New Tweet City Media Search' />
      {this.props.children}
    </div>;
  }
}

App.contextTypes = {
  history: React.PropTypes.object.isRequired
};

export default App
