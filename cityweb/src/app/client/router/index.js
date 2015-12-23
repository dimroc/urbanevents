import React from 'react'
import { render } from 'react-dom'
import Router from 'react-router'
import { Provider } from 'react-redux'
import { DevTools, DebugPanel, LogMonitor } from 'redux-devtools/lib/react'
import { createHistory } from 'history'
import toString from './toString'
import { Promise } from 'when'
import createRoutes from './routes'
import { createStore } from '../store'
import { syncReduxAndRouter } from 'redux-simple-router'

export function run() {
  // init promise polyfill
  window.Promise = window.Promise || Promise;
  // init fetch polyfill
  window.self = window;
  require('whatwg-fetch');
  const store = createStore(window['--app-initial']);
  window.store = store;

  if (process.env.NODE_ENV !== 'production'){
    store.subscribe(() => {
      console.log('%c[STORE]', 'color: green', store.getState());
    });
  }

  const history = createHistory()
  // Listen for changes to the current location. The
  // listener is called once immediately.
  let unlisten = history.listen(location => {
    console.log("debug listen: ", location);
  })

  syncReduxAndRouter(history, store)

  render(
    <Provider store={store} >
      <Router history={history}>{createRoutes({store, first: { time: true }})}</Router>
    </Provider>,
    document.getElementById('app')
  );

  if (process.env.NODE_ENV !== 'production'){
    const node = document.createElement('div');
    document.body.appendChild(node);
    render(
      <DebugPanel top right bottom>
        <DevTools store={store} monitor={LogMonitor} visibleOnLoad={false}/>
      </DebugPanel>,
      node
    );
  }
}

// Export it to render on the Golang sever, keep the name sync with -
// https://github.com/olebedev/go-starter-kit/blob/master/src/app/server/react.go#L65
export const renderToString = toString;

require('../css');

// Style live reloading
if (module.hot) {
  let c = 0;
  module.hot.accept('../css', () => {
    require('../css');
    const a = document.createElement('a');
    const link = document.querySelector('link[rel="stylesheet"]');
    a.href = link.href;
    a.search = '?' + c++;
    link.href = a.href;
  });
}
