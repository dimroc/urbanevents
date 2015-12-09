import React from 'react';
import Router from 'react-router';
import FluxComponent from 'flummox/component';
import Flux from '../flux';
import RenderToString from './RenderToString';
import loadProps from '#app/utils/loadProps';
import { Promise } from 'when';


import routes from './routes';

const flux = new Flux();

export function run() {
  // share flux instance
  window.flux = flux;
  // init promise polyfill
  window.Promise = window.Promise || Promise;
  // init fetch polyfill
  window.self = window;
  require('whatwg-fetch');

  fetch('/api/v1/conf').then((r) => {
    return r.json();
  }).then((conf) => {

    flux.getStore('app').setAppConfig(conf);
    if (process.env.NODE_ENV !== 'production'){
      flux.on('dispatch', (action) => {
        const {actionId, body} = action;
        console.log('%c[FLUX] %c%s', 'color: green', 'color: grey', actionId, body);
      });
    }
    Router.run(routes, Router.HistoryLocation, (Handler, state) => {
      const routeHandlerInfo = { flux, state };
      loadProps(state.routes, 'loadProps', routeHandlerInfo).then(()=> {
        React.render(
          <FluxComponent flux={flux}>
            <Handler />
          </FluxComponent>,
          document.getElementById('app')
        );
      });
    });

  });
}

export const renderToString = RenderToString;

require('../styles');

// Style live reload
if (module.hot) {
  let c = 0;
  module.hot.accept('../styles', () => {
    require('../styles');
    const a = document.createElement('a');
    const link = document.querySelector('link[rel="stylesheet"]');
    a.href = link.href;
    a.search = '?' + ++c;
    link.href = a.href;
  });
}
