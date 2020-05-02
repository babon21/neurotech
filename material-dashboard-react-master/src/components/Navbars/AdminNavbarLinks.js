import React from "react";
import Button from "components/CustomButtons/Button.js";
import {logoutUser} from "../../redux/actions/userActions"
import { connect } from 'react-redux';

function LogoutButton(props) {
  return (
      <Button
      onClick={props.logoutUser}
      >
        Logout
      </Button>
  );
}

const mapActionsToProps = {
  logoutUser
};

export default connect(
  null,
  mapActionsToProps
)(LogoutButton);
