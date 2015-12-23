import React from 'react';
import { Route, IndexRoute, Redirect } from 'react-router';
import App from '#app/components/app';
import Homepage from '#app/components/homepage';
import NotFound from '#app/components/not-found';

export default ({store}) => {
  return <Route path="/" component={App}>
    <IndexRoute component={Homepage}/>
    {/* Server redirect in action */}
    <Route path="*" component={NotFound}/>
  </Route>;
};
