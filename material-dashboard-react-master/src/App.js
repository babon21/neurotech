import React, { Component } from "react";
import { createBrowserHistory } from "history";
import { Router, Switch, Redirect } from "react-router-dom";

// core components
import Admin from "layouts/Admin.js";
import AuthRoute from "./util/AuthRoute"
import NoAuthRoute from "./util/NoAuthRoute"
import login from "./pages/signin"

// Redux
import { Provider } from 'react-redux';
import store from './redux/store';
import StudyMaterialsList from "views/StudyMaterialsList";

export const hist = createBrowserHistory();

class App extends Component {
  render() {
    return (
        <Provider store={store}>
          <Router history={hist}>
              <Switch>
                <AuthRoute path="/admin" component={Admin} />
                <NoAuthRoute exact path="/login" component={login} />
                <Redirect from="/" to="/admin/dashboard" />
                <AuthRoute exact path="/admin/study-materials/0" component={StudyMaterialsList} />
              </Switch>
          </Router>
        </Provider>
    );
  }
}

export default App;