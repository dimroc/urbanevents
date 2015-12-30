import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { getCitiesAsync } from '#app/actions';
import GitHubForkRibbon from 'react-github-fork-ribbon';

class App extends Component {
  render() {
    return <div>
      <Helmet title='New Tweet City Media Search' />
      <GitHubForkRibbon href="https://github.com/dimroc/urbanevents/"
                        target="_blank"
                        position="right">
        Fork me on GitHub
      </GitHubForkRibbon>
      {this.props.children}
    </div>;
  }
}

App.contextTypes = {
  history: React.PropTypes.object.isRequired
};

export default App
