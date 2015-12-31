import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { getCitiesAsync } from '#app/actions';
import GitHubForkRibbon from 'react-github-fork-ribbon';
import GoogleAnalytics from 'react-g-analytics';

class App extends Component {
  render() {
    return <div>
      <Helmet title='New Tweet City Media Search' />
      <GoogleAnalytics id="UA-43559997-3" />
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
