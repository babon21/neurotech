import React from 'react';
import { Route, Redirect } from 'react-router-dom';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

const NoAuthRoute = ({ component: Component, authenticated, ...rest }) => (
  <Route
    {...rest}
    render={(props) => {
      console.log('from no auth rote, authenticated', authenticated)
      return authenticated === false ? <Component {...props} /> : <Redirect to="/admin" />
    }
    }
  />
);

const mapStateToProps = (state) => ({
  authenticated: state.user.authenticated
});

NoAuthRoute.propTypes = {
  user: PropTypes.object
};

export default connect(mapStateToProps)(NoAuthRoute);
