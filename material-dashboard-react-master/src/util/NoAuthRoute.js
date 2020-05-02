import React from 'react';
import { Route, Redirect } from 'react-router-dom';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

const NoAuthRoute = ({ component: Component, authenticated, ...rest }) => (
  <Route
    {...rest}
    render={(props) => authenticated === "false" ? <Component {...props} /> : <Redirect to="/" />
    }
  />
);

// function NoAuthRoute({ component: Component, authenticated, ...rest }) {
//   return (
//     <Route
//       {...rest}
//       render={(props) => {
//         return authenticated === "false" ? <Component {...props} /> : <Redirect to="/" />
//       }
//       }
//     />
//   );
// }

const mapStateToProps = (state) => ({
  authenticated: state.user.authenticated
});

NoAuthRoute.propTypes = {
  user: PropTypes.object
};

export default connect(mapStateToProps)(NoAuthRoute);
