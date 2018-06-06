import React, { Component } from 'react';
import { ConnectedRouter, routerMiddleware } from 'react-router-redux';
import createHistory from 'history/createBrowserHistory';
import { Route } from 'react-router-dom';
import Home from '../scenes/Home';

const history = createHistory();
export const middleware = routerMiddleware(history);

class Navigation extends Component {
  render() {
    return (
      <ConnectedRouter history={history}>
        <div>
          <Route exact path='/' component={Home}/>
        </div>
      </ConnectedRouter>
    );
  }
}

export default Navigation;
